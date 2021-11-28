package shared

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"time"
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

// Unmarshal converts []byte into a C2SChatMessagePacket
func (m *S2CUpdatePlayersPosPacket) Unmarshal(b []byte) error {
	json.Unmarshal(b, &m)

	return nil
}

// S2CChatMessagePacket will relay the message sent from the client
// to server back to all or specific clients
type S2CChatMessagePacket struct {
	Username string    `json:"u"`    // Username for time being
	Type     ChatType  `json:"t"`    // All/PM/Command
	Time     time.Time `json:"time"` // Time is was recieved by server
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
