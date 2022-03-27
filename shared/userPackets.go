package shared

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
)

//////////////////
// User Packets //
//////////////////

// C2SUserLoginPacket is the most basic user data for login
type C2SUserLoginPacket struct {
	Username           string             `json:"u"`
	UserModelSelection UserModelSelection `json:"um"`
}

// Marshal converts a User into []byte
func (u *C2SUserLoginPacket) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, Reliable)
	binary.Write(buf, binary.LittleEndian, C2SUserLoginPacketType)

	userBytes, _ := json.Marshal(u)
	buf.Write(userBytes)

	return buf.Bytes(), nil
}

// Unmarshal converts []byte into a User
func (u *C2SUserLoginPacket) Unmarshal(b []byte) error {
	json.Unmarshal(b, &u)
	return nil
}

// S2CUserLoginResponsePacket is returned after user tries login
type S2CUserLoginResponsePacket struct {
	// ID uint32 // Server sends a new unique ID to assosiate with the player
	Success             bool                   `json:"s"`
	UsersModelSelection map[string]OtherPlayer `json:"ums"`
}

type OtherPlayer struct {
	Username           string
	UserModelSelection UserModelSelection
}

// Marshal converts a UserLoginResponsePacket into []byte
func (u *S2CUserLoginResponsePacket) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, Reliable)
	binary.Write(buf, binary.LittleEndian, S2CUserLoginResponsePacketType)

	userBytes, _ := json.Marshal(u)
	buf.Write(userBytes)

	return buf.Bytes(), nil
}

// Unmarshal converts []byte into a UserLoginResponsePacket
func (u *S2CUserLoginResponsePacket) Unmarshal(b []byte) error {
	json.Unmarshal(b, &u)
	return nil
}

// C2SUserFinishedSetupPacket is returned after client finishes setup
type C2SUserFinishedSetupPacket struct {
}

// Marshal converts a UserFinishedSetupPacket into []byte
func (u *C2SUserFinishedSetupPacket) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, Reliable)
	binary.Write(buf, binary.LittleEndian, C2SUserFinishedSetupPacketType)

	userBytes, _ := json.Marshal(u)
	buf.Write(userBytes)

	return buf.Bytes(), nil
}

// Unmarshal converts []byte into a UserFinishedSetupPacket
func (u *C2SUserFinishedSetupPacket) Unmarshal(b []byte) error {
	json.Unmarshal(b, &u)
	return nil
}
