package pdu

import (
	"encoding/json"
	"fmt"
)

// ESMClass see SMPP v5, section 4.7.12 (125p)// ESMClass represents the ESM class structure
type ESMClass struct {
	MessageMode  byte // __ ____ **
	MessageType  byte // __ **** __
	UDHIndicator bool // _* ____ __
	ReplyPath    bool // *_ ____ __
}

// ReadByte converts ESMClass to a byte
func (e ESMClass) ReadByte() (c byte, err error) {
	c |= e.MessageMode & 0b11
	c |= (e.MessageType & 0b1111) << 2
	c |= getBool(e.UDHIndicator) << 6
	c |= getBool(e.ReplyPath) << 7
	return
}

// WriteByte converts a byte to ESMClass
func (e *ESMClass) WriteByte(c byte) error {
	e.MessageMode = c & 0b11
	e.MessageType = c >> 2 & 0b1111
	e.UDHIndicator = c>>6&0b1 == 1
	e.ReplyPath = c>>7&0b1 == 1
	return nil
}

// String returns the string representation of ESMClass
func (e ESMClass) String() string {
	c, _ := e.ReadByte()
	return fmt.Sprintf("%d", c)
}

// MarshalJSON marshals ESMClass to JSON
func (e ESMClass) MarshalJSON() (data []byte, err error) {
	c, _ := e.ReadByte()
	return json.Marshal(c)
}
