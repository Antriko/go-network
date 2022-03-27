package server

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Antriko/go-network/shared"
	"github.com/Antriko/go-network/world"
	"github.com/gookit/color"
)

var red = color.New(color.FgBlack, color.BgRed).Render
var yellow = color.New(color.FgBlack, color.BgYellow).Render
var green = color.New(color.FgBlack, color.BgGreen).Render
var magenta = color.New(color.FgBlack, color.BgMagenta).Render

var worldMap *world.WorldStruct
var worldSize int

func server() {
	log.SetFlags(log.Lshortfile)

	AllConnections = make(Connections)
	log.Println(magenta(" - Starting Server - "))

	worldSize = 10
	worldMap = world.CreateWorld(worldSize)

	go func() {
		// TODO terminal commands
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			return
		}

		// remove the delimeter from the string
		input = strings.TrimSuffix(input, "\n")
		fmt.Println(input)
	}()

	go serverUpdates()
	Handler := shared.NewConnectionHandler()
	Handler.Listen("127.0.0.1", serverPacketHandler)
}

// serverPacketHandler validates user and appends user information
func serverPacketHandler(conn *shared.DualConnection) {
	// Block until connection is setup
	// Have some verifiction before proceeding	// i.e. check if username is in use
	connected := false
	for !connected {
		select {
		case r := <-conn.DataReadChan:
			log.Printf("%T", r)

			switch typed := r.(type) {
			case *shared.A2APingPacket:
				log.Println("server pinged", typed.Reliability)
				conn.DataWriteChan <- &shared.A2APongPacket{Reliability: typed.Reliability}
			case *shared.A2APongPacket:
				log.Println("server ponged", typed.Reliability)

			case *shared.C2SUserLoginPacket: // User logging in - verify some information if needed
				log.Println(yellow(" User logging in "), conn.DataWriteChan)

				// Send model selection of users that are logged in
				allUserModelSelection := make(map[string]shared.OtherPlayer)
				for key, value := range AllConnections {
					log.Println(key, value.UserModelSelection)
					allUserModelSelection[value.Username] = shared.OtherPlayer{
						Username:           value.Username,
						UserModelSelection: value.UserModelSelection,
					}
				}

				conn.DataWriteChan <- &shared.S2CUserLoginResponsePacket{
					Success:             true,
					UsersModelSelection: allUserModelSelection,
				}

				conn.DataWriteChan <- &shared.S2CSendWorldPacket{
					WorldTiles: worldMap.Tiles,
					Size:       worldMap.Size,
				}

				AllConnections[conn] = &UserConnection{
					Connection:         conn,
					Username:           typed.Username,
					ID:                 newID(),
					UserModelSelection: typed.UserModelSelection,
				}

				for userConn := range AllConnections {
					// Inform all users of user connect
					userConn.DataWriteChan <- &shared.S2CChatMessagePacket{
						Username: typed.Username,
						Type:     shared.UserConnect,
						Time:     time.Now(),
						Message:  "",
					}

					// Send information to all users
					userConn.DataWriteChan <- &shared.S2CUserInformationPacket{
						User: shared.PublicUserInformation{
							Username:           typed.Username,
							UserModelSelection: typed.UserModelSelection,
						},
					}
				}

			case *shared.C2SUserFinishedSetupPacket:
				log.Println(yellow(" User finished setup "))
				connected = true
				break
			default:
				log.Println("Wrong packet sent, awaiting initial UserLoginPacket, got: ", typed)
			}
		}
	}

	// user packet handler
	for {
		player, ok := AllConnections[conn]
		if !ok {
			return
		}

		isConn := true
		for isConn {
			select {
			case r := <-conn.DataReadChan:
				log.Printf("%T", r)

				switch typed := r.(type) {
				case *shared.C2SUpdatePlayerPosPacket:
					userCoordsMap[player.Username] = shared.Coords{
						Username: player.Username,
						X:        typed.X,
						Y:        typed.Y,
						Z:        typed.Z,
						Facing:   typed.Facing,
					}

				case *shared.C2SChatMessagePacket:
					for userConn := range AllConnections {
						userConn.DataWriteChan <- &shared.S2CChatMessagePacket{
							Username: typed.Username,
							Type:     typed.Type,
							Time:     time.Now(),
							Message:  typed.Message,
						}
					}
				default: // keep alive if doing nothing
					log.Println()
					conn.DataWriteChan <- &shared.A2APingPacket{Reliability: shared.Reliable}

				}

			case <-conn.DataErrChan:
				isConn = false
			}
		}
		removeConn(conn)
		break
	}
}

func serverUpdates() {
	go func() {
		timer := time.NewTimer(time.Second / 50)
		for {
			select {
			case <-timer.C: // Send updates
				for userConn := range AllConnections {
					// log.Println(yellow(" Sending coords update "))
					userConn.DataWriteChan <- &shared.S2CUpdatePlayersPosPacket{
						Coords: userCoordsMap,
					}
				}
				timer.Reset(time.Second)
			}
		}
	}()
}

func removeConn(conn *shared.DualConnection) {
	username := AllConnections[conn].Username
	delete(AllConnections, conn)
	delete(userCoordsMap, username)
	log.Println(red(" User disonnected "))
	if err := conn.Close(); err != nil {
		log.Println(err)
	}

	for userConn := range AllConnections {
		// Inform all users of user disconnect
		userConn.DataWriteChan <- &shared.S2CChatMessagePacket{
			Username: username,
			Type:     shared.UserDisconnect,
			Time:     time.Now(),
			Message:  "",
		}
	}
}

var totalIDs uint32

func newID() uint32 {
	totalIDs++
	return totalIDs
}
