package shared

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/Antriko/go-network/world"
)

// C2SChatMessagePacket will send the server the chat message
// so the server will be able to send it to everyone else who's connected
type S2CUpdatePlayersPosPacket struct {
	Coords map[string]Coords `json:"c"`
}

type Coords struct {
	Username string  `json:"u"`
	X        float32 `json:"x"`
	Y        float32 `json:"y"`
	Z        float32 `json:"z"`
	Facing   float32 `json:"f"`
}

// Marshal converts a C2SChatMessagePacket into []byte
func (m *S2CUpdatePlayersPosPacket) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, Unreliable) // DTLS
	binary.Write(buf, binary.LittleEndian, S2CUpdatePlayersPosPacketType)
	userBytes, _ := json.Marshal(m)

	buf.Write(userBytes)
	return buf.Bytes(), nil
}

// Unmarshal converts []byte into a S2CUpdatePlayersPosPacket
func (m *S2CUpdatePlayersPosPacket) Unmarshal(b []byte) error {
	json.Unmarshal(b, &m)

	return nil
}

// S2CChatMessagePacket will relay the message sent from the client
// to server back to all or specific clients
type S2CChatMessagePacket struct {
	Username string    `json:"u"`    // Username for time being
	Type     ChatType  `json:"t"`    // All/PM/Command
	Time     time.Time `json:"time"` // Time it was recieved by server
	Message  string    `json:"m"`
}

// Marshal converts a S2CChatMessagePacket into []byte
func (m *S2CChatMessagePacket) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, Unreliable)
	binary.Write(buf, binary.LittleEndian, S2CChatMessagePacketType)

	userBytes, _ := json.Marshal(m)
	buf.Write(userBytes)

	return buf.Bytes(), nil
}

// Unmarshal converts []byte into a S2CChatMessagePacket
func (m *S2CChatMessagePacket) Unmarshal(b []byte) error {
	json.Unmarshal(b, &m)
	return nil
}

type PublicUserInformation struct {
	Username           string             `json:"u"` // Username for time being
	UserModelSelection UserModelSelection `json:"um"`
}

// S2CUserInformation will send all players the information and model that the user has chosen
type S2CUserInformationPacket struct {
	User PublicUserInformation
}

// Marshal converts a S2CChatMessagePacket into []byte
func (m *S2CUserInformationPacket) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, Unreliable)
	binary.Write(buf, binary.LittleEndian, S2CUserInformationPacketType)

	userBytes, _ := json.Marshal(m)
	buf.Write(userBytes)

	return buf.Bytes(), nil
}

// Unmarshal converts []byte into a S2CChatMessagePacket
func (m *S2CUserInformationPacket) Unmarshal(b []byte) error {
	json.Unmarshal(b, &m)
	return nil
}

// *world.WorldStruct but want to ignore instances
type S2CSendWorldPacket struct {
	WorldTiles [][]world.MapTile `json:"t"`
	Size       int               `json:"s"`
}

// Marshal converts a UserFinishedSetupPacket into []byte
func (u *S2CSendWorldPacket) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, Unreliable)
	binary.Write(buf, binary.LittleEndian, S2CSendWorldPacketType)

	userBytes, _ := json.Marshal(u)
	buf.Write(userBytes)

	return buf.Bytes(), nil
}

// Unmarshal converts []byte into a S2CSendWorldPacket
func (u *S2CSendWorldPacket) Unmarshal(b []byte) error {
	json.Unmarshal(b, u)
	return nil
}

// If world too big, then better to split the world into chunks and send those rather than entire map at once.
