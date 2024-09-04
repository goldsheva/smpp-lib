package pdu

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// DstAddress ...
type DstAddress struct {
	TON  byte   `json:"dest_addr_ton"` // see SMPP v5, section 4.7.1 (113p)
	NPI  byte   `json:"dest_addr_npi"` // see SMPP v5, section 4.7.2 (113p)
	Dest string `json:"destination_addr"`
}

// DestinationAddresses ...
type DestinationAddresses struct {
	Addresses        []DstAddress
	DistributionList []string
}

// UnsuccessfulRecord ...
type UnsuccessfulRecord struct {
	DestAddr        DstAddress
	ErrorStatusCode CommandStatus
}

// UnsuccessfulRecords ...
type UnsuccessfulRecords []UnsuccessfulRecord

// ReadFrom ...
func (p *DstAddress) ReadFrom(r io.Reader) (n int64, err error) {
	buf := bufio.NewReader(r)
	p.TON, err = buf.ReadByte()
	if err == nil {
		p.NPI, err = buf.ReadByte()
	}
	if err == nil {
		p.Dest, err = readCString(buf)
	}
	return
}

// WriteTo ...
func (p DstAddress) WriteTo(w io.Writer) (n int64, err error) {
	var buf bytes.Buffer
	buf.WriteByte(p.TON)
	buf.WriteByte(p.NPI)
	writeCString(&buf, p.Dest)
	return buf.WriteTo(w)
}

// String returns the string representation of the DstAddress.
func (p DstAddress) String() string {
	if p.TON == TypeOfNumberInternational && p.NPI == NumberingPlanE164 && len(p.Dest) > 0 && p.Dest[0] != '+' {
		return "+" + p.Dest
	}
	return p.Dest
}

// ReadFrom ...
func (p *DestinationAddresses) ReadFrom(r io.Reader) (n int64, err error) {
	buf := bufio.NewReader(r)
	count, err := buf.ReadByte()
	if err != nil {
		err = errors.New("InvalidCommandLength")
		return
	}
	*p = DestinationAddresses{}
	var destFlag byte
	var value string
	var address DstAddress
	for i := byte(0); i < count; i++ {
		switch destFlag, _ = buf.ReadByte(); destFlag {
		case 1:
			if _, err = address.ReadFrom(buf); err == nil {
				p.Addresses = append(p.Addresses, address)
			}
		case 2:
			if value, err = readCString(buf); err == nil {
				p.DistributionList = append(p.DistributionList, value)
			}
		default:
			err = errors.New("InvalidDestFlag")
			return
		}
		if err != nil {
			err = errors.New("InvalidCommandLength")
			return
		}
	}
	return
}

// WriteTo ...
func (p DestinationAddresses) WriteTo(w io.Writer) (n int64, err error) {
	length := len(p.Addresses) + len(p.DistributionList)
	if length > 0xFF {
		err = errors.New("InvalidDestCount")
		return
	}
	var buf bytes.Buffer
	buf.WriteByte(byte(length))
	for _, address := range p.Addresses {
		buf.WriteByte(1)
		_, _ = address.WriteTo(&buf)
	}
	for _, distribution := range p.DistributionList {
		buf.WriteByte(2)
		writeCString(&buf, distribution)
	}
	return buf.WriteTo(w)
}

// String ...
func (i UnsuccessfulRecord) String() string {
	return fmt.Sprintf("%s#%s", i.DestAddr, i.ErrorStatusCode)
}

// ReadFrom ...
func (p *UnsuccessfulRecords) ReadFrom(r io.Reader) (n int64, err error) {
	buf := bufio.NewReader(r)
	count, err := buf.ReadByte()
	if err != nil {
		err = errors.New("InvalidCommandLength")
		return
	}
	items := UnsuccessfulRecords{}
	var item UnsuccessfulRecord
	for i := byte(0); i < count; i++ {
		_, err = item.DestAddr.ReadFrom(buf)
		if err == nil {
			err = binary.Read(buf, binary.BigEndian, &item.ErrorStatusCode)
		}
		if err != nil {
			err = errors.New("InvalidCommandLength")
			return
		}
		items = append(items, item)
	}
	*p = items
	return
}

// WriteTo ...
func (p UnsuccessfulRecords) WriteTo(w io.Writer) (n int64, err error) {
	if len(p) > 0xFF {
		err = errors.New("ItemTooMany")
		return
	}
	var buf bytes.Buffer
	buf.WriteByte(byte(len(p)))
	for _, item := range p {
		_, _ = item.DestAddr.WriteTo(&buf)
		_ = binary.Write(&buf, binary.BigEndian, item.ErrorStatusCode)
	}
	return buf.WriteTo(w)
}

// Function to auto-detect dest_addr_ton and dest_addr_npi
func (da *DstAddress) AutoDetectTONNPI() {
	// Default values
	da.TON = TypeOfNumberNational
	da.NPI = NumberingPlanE164

	// Check if destination address starts with '+'
	if strings.HasPrefix(da.Dest, "+") {
		da.Dest = da.Dest[1:] // Remove '+'
		da.TON = TypeOfNumberInternational
	} else {
		// Check if remaining part is numeric
		numericOnly := true
		for _, char := range da.Dest {
			if !unicode.IsDigit(char) {
				numericOnly = false
				break
			}
		}
		if numericOnly {
			da.TON = TypeOfNumberInternational
		} else {
			da.TON = TypeOfNumberAlphanumeric
			da.NPI = NumberingPlanUnknown
		}
	}
}
