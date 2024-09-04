package pdu

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	_ID                   = "id"
	_ESMClass             = "ESMClass"
	MaxShortMessageLength = 140 // bytes
)

// CommandID see SMPP v5, section 4.7.5 (115p)
type CommandID uint32

// Header ...
type Header struct {
	CommandLength uint32        `json:"command_length"`
	CommandID     CommandID     `json:"command_id"`
	CommandStatus CommandStatus `json:"command_status"`
	Sequence      int32         `json:"sequence_number"`
}

// PDUError ...
type PDUError struct {
	CommandStatus CommandStatus
	Err           error
}

// ReadPDUHeader ...
func ReadPDUHeader(r io.Reader, header *Header) (err error) {
	err = binary.Read(r, binary.BigEndian, header)
	if err == nil && (header.CommandLength < 16 || header.CommandLength > 0x10000) {
		err = errors.New("InvalidCommandLength")
	}
	return
}
