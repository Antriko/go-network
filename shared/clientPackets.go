package shared

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"math"
)

///////////////////////////
// Client update packets //
///////////////////////////

// C2SUpdatePlayerPosPacket updates the server with the clients desired position
type C2SUpdatePlayerPosPacket struct {
	X, Y, Z, Facing float32
}

// Marshal converts a MapPosPacket into []byte
func (m *C2SUpdatePlayerPosPacket) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, Unreliable)
	binary.Write(buf, binary.LittleEndian, C2SUpdatePlayerPosPacketType)
	binary.Write(buf, binary.LittleEndian, m.X)
	binary.Write(buf, binary.LittleEndian, m.Y)
	binary.Write(buf, binary.LittleEndian, m.Z)
	binary.Write(buf, binary.LittleEndian, m.Facing)

	return buf.Bytes(), nil
}

// Unmarshal converts []byte into a MapPosPacket
func (m *C2SUpdatePlayerPosPacket) Unmarshal(b []byte) error {
	if len(b) != 4+4+4+4 {
		return ErrMissingPacketData
	}

	m.X = math.Float32frombits(binary.LittleEndian.Uint32(b[:4]))
	m.Y = math.Float32frombits(binary.LittleEndian.Uint32(b[4:8]))
	m.Z = math.Float32frombits(binary.LittleEndian.Uint32(b[8:12]))
	m.Facing = math.Float32frombits(binary.LittleEndian.Uint32(b[12:16]))

	return nil
}

// C2SChatMessagePacket will send the server the chat message
// so the server will be able to send it to everyone else who's connected
type C2SChatMessagePacket struct {
	Username string   `json:"u"` // Username for time being
	Type     ChatType `json:"t"` // All/PM/Command
	Message  string   `json:"m"`
}

// Different types of messages
type ChatType uint32

const (
	AllChat ChatType = iota
	PrivateMessage
	UserConnect
	UserDisconnect

	CommandWorldSize
	// Different commands i.e. /list
)

// Marshal converts a C2SChatMessagePacket into []byte
func (m *C2SChatMessagePacket) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, Unreliable)
	binary.Write(buf, binary.LittleEndian, C2SChatMessagePacketType)

	userBytes, _ := json.Marshal(m)
	buf.Write(userBytes)

	return buf.Bytes(), nil
}

// Unmarshal converts []byte into a C2SChatMessagePacket
func (m *C2SChatMessagePacket) Unmarshal(b []byte) error {
	json.Unmarshal(b, &m)
	return nil
}

type UserModelSelection struct {
	Accessory int
	Hair      int
	Head      int
	Body      int
	Bottom    int
}
