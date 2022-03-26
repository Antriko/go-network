package shared

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// Reliability is a flag that determins how the packet is sent
// True sends over TLS, false sends over DTLS
type Reliability uint8

// Reliability
const (
	Reliable Reliability = iota
	Unreliable
	Both
)

// Packets

// PacketType specifies what a packet is
// Packets represent the data which is sent, NOT the game struct itself
type PacketType uint16

const (
	// Any to Any
	A2APingPacketType PacketType = iota // keep alive etc
	A2APongPacketType

	// Client to Server
	C2SUserLoginPacketType         // Init details of user
	C2SUpdatePlayerPosPacketType   // Send X Y playerPos
	C2SUserFinishedSetupPacketType // Confirmation that player is finished setup

	C2SChatMessagePacketType // Send message to server

	// Server to Client
	S2CUserLoginResponsePacketType // Maybe get model selection of users already logged in?
	S2CUpdatePlayersPosPacketType  // Send map of all players pos
	S2CChatMessagePacketType       // User recieves a new chat message
	S2CUserInformationPacketType   // Send user information back to client
	S2CSendWorldPacketType         // Send same world to all clients
)

// Serializable ensures that the struct has a Bytes() function
type Serializable interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

// Errors
var (
	ErrUndefinedPacket   = errors.New("Undefined packet type")
	ErrMissingPacketData = errors.New("Packet data length incorrect")
)

//////////////////
// Misc Packets //
//////////////////

// A2APingPacket is an empty packet which can be used as a keep alive packet or
// to check if a connection is open
type A2APingPacket struct {
	Reliability
}

// Marshal converts a User into []byte
func (u *A2APingPacket) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, u.Reliability)
	binary.Write(buf, binary.LittleEndian, A2APingPacketType)
	return buf.Bytes(), nil
}

// Unmarshal converts []byte into a User
func (u *A2APingPacket) Unmarshal(b []byte) error {
	u.Reliability = Reliability(b[0])
	return nil
}

// A2APongPacket is an empty packet which is sent as a response to PingPacket
type A2APongPacket struct {
	Reliability
}

// Marshal converts a User into []byte
func (u *A2APongPacket) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, u.Reliability)
	binary.Write(buf, binary.LittleEndian, A2APongPacketType)
	return buf.Bytes(), nil
}

// Unmarshal converts []byte into a User
func (u *A2APongPacket) Unmarshal(b []byte) error {
	u.Reliability = Reliability(b[0])
	return nil
}

// BytesToStruct returns the correct struct depending on what the PacketType is
// Don't forget to type check/cast before using it
func BytesToStruct(b []byte) (interface{}, error) {
	origB := b[:]
	code := b[1:3] // ignore Reliability byte
	b = b[3:]
	var ret Serializable
	switch PacketType(binary.LittleEndian.Uint16(code)) {
	case A2APingPacketType:
		// Ping and Pong packets require the full bytes to be passed
		ret = &A2APingPacket{}
		ret.Unmarshal(origB)
		return ret, nil
	case A2APongPacketType:
		ret = &A2APongPacket{}
		ret.Unmarshal(origB)
		return ret, nil
	case C2SUserLoginPacketType:
		ret = &C2SUserLoginPacket{}
	case S2CUserLoginResponsePacketType:
		ret = &S2CUserLoginResponsePacket{}
	case C2SUserFinishedSetupPacketType:
		ret = &C2SUserFinishedSetupPacket{}
	case C2SUpdatePlayerPosPacketType:
		ret = &C2SUpdatePlayerPosPacket{}
	case C2SChatMessagePacketType:
		ret = &C2SChatMessagePacket{}
	case S2CChatMessagePacketType:
		ret = &S2CChatMessagePacket{}
	case S2CUpdatePlayersPosPacketType:
		ret = &S2CUpdatePlayersPosPacket{}
	case S2CUserInformationPacketType:
		ret = &S2CUserInformationPacket{}
	case S2CSendWorldPacketType:
		ret = &S2CSendWorldPacket{}
	default:
		return nil, ErrUndefinedPacket
	}
	return ret, ret.Unmarshal(b)
}
