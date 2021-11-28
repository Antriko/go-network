package shared

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/gookit/color"
	"github.com/pion/dtls/v2"
)

var red = color.New(color.FgBlack, color.BgRed).Render
var yellow = color.New(color.FgBlack, color.BgYellow).Render
var green = color.New(color.FgBlack, color.BgGreen).Render
var magenta = color.New(color.FgBlack, color.BgMagenta).Render

// Errors
var (
	ErrConnClosed = errors.New("TLS connection closed")
)

// NewConnectionHandler returns a new ConnectionHandler
func NewConnectionHandler() *ConnectionHandler {
	ch := &ConnectionHandler{}
	return ch
}

type ConnectionHandler struct {
	*DualConnection

	TLSListener  net.Listener
	DTLSListener net.Listener
}

// DualConnection binds both channels in one struct
type DualConnection struct {
	TLSConnection  net.Conn
	DTLSConnection net.Conn

	DataErrChan   chan error
	DataReadChan  chan interface{}
	DataWriteChan chan Serializable

	HandlerFunction
	writeChan chan []byte
	timeout   *time.Timer
}

// HandlerFunction is a callback used to handle connection data
type HandlerFunction func(conn *DualConnection)

// Listen sets up the ConnectionHandler to listen
func (c *ConnectionHandler) Listen(ip string, handlerFunc HandlerFunction) {
	// Certificate setup
	cer, err := tls.LoadX509KeyPair("shared/certs/server.pub.pem", "shared/certs/server.pem")
	if err != nil {
		log.Println(err)
		return
	}
	certPool := x509.NewCertPool()
	cert, err := x509.ParseCertificate(cer.Certificate[0])
	certPool.AddCert(cert)

	// DTLS setup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c.DTLSListener, err = dtls.Listen("udp", &net.UDPAddr{IP: net.ParseIP(ip), Port: 8080},
		&dtls.Config{
			Certificates:         []tls.Certificate{cer},
			ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
			ClientCAs:            certPool,
			// Create timeout context for accepted connection.
			ConnectContextMaker: func() (context.Context, func()) {
				return context.WithTimeout(ctx, time.Second)
			},
		})

	// TLS setup
	c.TLSListener, err = tls.Listen("tcp", ":8070",
		&tls.Config{
			Certificates: []tls.Certificate{cer},
		})
	if err != nil {
		log.Println(err)
		return
	}

	// Handle connections
	for {
		// log.Println("Listening", c.TLSListener.Addr().String())
		tlsConn, err := c.TLSListener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		// log.Println("tls connection: ", tlsConn.RemoteAddr())
		// FIXME Have to write something or else it hangs??
		pingRel := A2APingPacket{Reliability: Reliable}
		b, err := pingRel.Marshal()
		tlsConn.Write(b)

		dtlsConn, err := c.DTLSListener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		// log.Println("dtls connection: ", dtlsConn.RemoteAddr())
		pingUnrel := A2APingPacket{Reliability: Unreliable}
		b, err = pingUnrel.Marshal()
		tlsConn.Write(b)

		c.DualConnection = &DualConnection{
			TLSConnection:   tlsConn,
			DTLSConnection:  dtlsConn,
			DataErrChan:     make(chan error, 64),
			DataReadChan:    make(chan interface{}, 64),
			DataWriteChan:   make(chan Serializable, 64),
			HandlerFunction: handlerFunc,
			writeChan:       make(chan []byte, 1280),
			timeout:         time.NewTimer(time.Second * 10), // TODO env for this
		}

		go c.DualConnection.HandlerFunction(c.DualConnection)
		go c.DualConnection.PollConnection(func() {})
	}
}

// Dial sets up the ConnectionHandler to dial
func (c *ConnectionHandler) Dial(ip string, handlerFunc HandlerFunction) {
	// Certificate setup
	cert, err := ioutil.ReadFile("shared/certs/server.pub.pem")
	if err != nil {
		log.Fatal("Couldn't load file", err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(cert)

	var tlsConn *tls.Conn
	var dtlsConn *dtls.Conn

	// TLS setup
	tlsConn, err = tls.Dial("tcp", ip+":8070",
		&tls.Config{
			RootCAs: certPool,
		})
	if err != nil {
		log.Println(err)
		return
	}

	// DTLS setup
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	dtlsConn, err = dtls.DialWithContext(ctx, "udp", &net.UDPAddr{IP: net.ParseIP(ip), Port: 8080},
		&dtls.Config{
			ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
			RootCAs:              certPool,
		})
	if err != nil {
		log.Println(err)
		return
	}

	c.DualConnection = &DualConnection{
		TLSConnection:   tlsConn,
		DTLSConnection:  dtlsConn,
		HandlerFunction: handlerFunc,
		DataErrChan:     make(chan error, 64),
		DataReadChan:    make(chan interface{}, 64),
		DataWriteChan:   make(chan Serializable, 64),
		writeChan:       make(chan []byte, 1280),
		timeout:         time.NewTimer(time.Second * 10), // TODO env for this
	}

	go c.DualConnection.HandlerFunction(c.DualConnection)
	go c.PollConnection(func() {})
}

// PollConnection waits for connection reads and also reads/writes from the
// channels provided
func (c *DualConnection) PollConnection(removeConnectionCallback func()) {
	buf := make([]byte, 4096)
	go func() { // TLS Connection
		for {
			n, err := c.TLSConnection.Read(buf)
			if err != nil {
				// log.Println(err)
				c.DTLSConnection.Close()
				c.DataErrChan <- ErrConnClosed
				return
			}
			c.resetTimer()

			if s, err := BytesToStruct(buf[:n]); err == nil {
				log.Println(yellow(" TLS "), s)
				c.DataReadChan <- s
			} else {
				// log.Println(err)
			}
		}
	}()
	go func() { // DTLS Connection
		for {
			n, err := c.DTLSConnection.Read(buf)
			if err != nil {
				// log.Println(err)
				c.DataErrChan <- ErrConnClosed
				return
			}
			c.resetTimer()

			if s, err := BytesToStruct(buf[:n]); err == nil {
				log.Println(yellow(" DTLS "), s)
				c.DataReadChan <- s
			} else {
				// log.Println(err)
			}
		}
	}()

	for {
		select {
		case <-c.timeout.C:
			removeConnectionCallback()
			c.DTLSConnection.Close()
			c.TLSConnection.Close()
			return

		case w := <-c.DataWriteChan:
			// struct -> bytes -> send
			// log.Println(w)
			b, _ := w.Marshal()
			c.writeChan <- b

		case w := <-c.writeChan:
			if Reliability(w[0]) == Reliable || Reliability(w[0]) == Both {
				if _, err := c.TLSConnection.Write(w); err != nil {
					log.Println(err)
					c.DTLSConnection.Close()
					c.DataErrChan <- ErrConnClosed
					return
				}
			}
			if Reliability(w[0]) == Unreliable || Reliability(w[0]) == Both {
				if _, err := c.DTLSConnection.Write(w); err != nil {
					log.Println(err)
					c.TLSConnection.Close()
					c.DataErrChan <- ErrConnClosed
					return
				}
			}
		}
	}
}

func (c *DualConnection) resetTimer() {
	// Reset timeout
	if !c.timeout.Stop() {
		log.Println("Timed out")
		<-c.timeout.C
	}
	c.timeout.Reset(time.Second * 10)
}
