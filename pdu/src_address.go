package pdu

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

// SrcAddress ...
type SrcAddress struct {
	TON    byte   `json:"source_addr_ton"` // see SMPP v5, section 4.7.1 (113p)
	NPI    byte   `json:"source_addr_npi"` // see SMPP v5, section 4.7.2 (113p)
	Source string `json:"source_addr"`
}

// ReadFrom ...
func (p *SrcAddress) ReadFrom(r io.Reader) (n int64, err error) {
	buf := bufio.NewReader(r)
	p.TON, err = buf.ReadByte()
	if err == nil {
		p.NPI, err = buf.ReadByte()
	}
	if err == nil {
		p.Source, err = readCString(buf)
	}
	return
}

// WriteTo ...
func (p SrcAddress) WriteTo(w io.Writer) (n int64, err error) {
	var buf bytes.Buffer
	buf.WriteByte(p.TON)
	buf.WriteByte(p.NPI)
	writeCString(&buf, p.Source)
	return buf.WriteTo(w)
}

// String ...
func (p SrcAddress) String() string {
	if p.TON == TypeOfNumberInternational && p.NPI == NumberingPlanE164 && len(p.Source) > 0 && p.Source[0] != '+' {
		return "+" + p.Source
	}
	return p.Source
}

// Function to auto-detect source_addr_ton and source_addr_npi
func (p *SrcAddress) AutoDetectTONNPI() {
	// Default values
	p.TON = TypeOfNumberNational
	p.NPI = NumberingPlanE164

	// Check if source address starts with '+'
	if strings.HasPrefix(p.Source, "+") {
		p.Source = p.Source[1:] // Remove '+'

		// Check if remaining part is numeric
		if _, err := fmt.Sscanf(p.Source, "%d", &p.Source); err == nil {
			p.TON = TypeOfNumberInternational
		}
	} else {
		// Check if alphanumeric
		if !isNumeric(p.Source) {
			p.TON = TypeOfNumberAlphanumeric
			p.NPI = NumberingPlanUnknown
		}
	}
}
