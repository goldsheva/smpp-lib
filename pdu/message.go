package pdu

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"reflect"

	"github.com/goldsheva/smpp-lib/coding"
	"github.com/sirupsen/logrus"
)

// ShortMessage ...
type ShortMessage struct {
	DataCoding       coding.DataCoding
	DefaultMessageID byte // see SMPP v5, section 4.7.27 (134p)
	UDHeader         UserDataHeader
	Message          []byte
}

// MarshalJSON ...
func (p *ShortMessage) MarshalJSON() (data []byte, err error) {

	return json.Marshal(&struct {
		DataCoding       coding.DataCoding `json:"data_coding"`
		DefaultMessageID byte              `json:"sm_default_msg_id"`
		UDHeader         UserDataHeader    `json:"UDH,omitempty"`
		Message          string            `json:"short_message"`
	}{
		DataCoding:       p.DataCoding,
		DefaultMessageID: p.DefaultMessageID,
		UDHeader:         p.UDHeader,
		Message:          p.Decode(),
	})
}

// ReadFrom ...
func (p *ShortMessage) ReadFrom(r io.Reader) (n int64, err error) {
	buf := bufio.NewReader(r)
	if p.DataCoding != coding.NoCoding {
		cod, _ := buf.ReadByte()
		p.DataCoding = coding.DataCoding(cod)
	}
	p.DefaultMessageID, err = buf.ReadByte()
	if err == nil {
		var length byte
		if length, err = buf.ReadByte(); err == nil && p.UDHeader != nil {
			_, err = p.UDHeader.ReadFrom(buf)
		}
		if err == nil {
			p.Message = make([]byte, length-byte(p.UDHeader.Len()))
			_, err = buf.Read(p.Message)
		}
	}
	return
}

// WriteTo ...
func (p ShortMessage) WriteTo(w io.Writer) (n int64, err error) {
	var buf bytes.Buffer
	if p.DataCoding != coding.NoCoding {
		buf.WriteByte(byte(p.DataCoding))
	}
	buf.WriteByte(p.DefaultMessageID)
	start := buf.Len()
	buf.WriteByte(0)
	_, err = p.UDHeader.WriteTo(&buf)
	if err != nil {
		return
	}
	buf.Write(p.Message)
	data := buf.Bytes()
	data[start] = byte(len(data) - 1 - start)
	return buf.WriteTo(w)
}

// Prepare ...
func (p *ShortMessage) Prepare(pdu interface{}) {
	if _, ok := pdu.(*ReplaceSM); ok {
		p.DataCoding = coding.NoCoding
	} else if p.UDHeader == nil {
		v := reflect.ValueOf(pdu).Elem().FieldByName(_ESMClass)
		if v.IsValid() && v.Interface().(ESMClass).UDHIndicator {
			p.UDHeader = UserDataHeader{}
		}
	}
}

// Parse ...
func (p *ShortMessage) Parse() (string, error) {
	encoder := p.DataCoding.Encoding()

	if encoder == nil {
		return string(p.Message), nil
	}
	decoded, err := encoder.NewDecoder().Bytes(p.Message)

	return string(decoded), err
}

// Decode ...
func (p *ShortMessage) Decode() string {
	var message string

	// is GSM7 encoding
	if p.DataCoding == coding.GSM7BitCoding {
		message = coding.DecodeGSM7(p.Message)

	} else {
		// Get encoder
		encoder := p.DataCoding.Encoding()
		if encoder != nil {

			decoded, err := encoder.NewDecoder().Bytes(p.Message)
			if err == nil {
				message = string(decoded)
			} else {
				logrus.WithFields(logrus.Fields{"worker": "pdu.message"}).Errorf("Can't Decode message with data_coding %d: %s", int(p.DataCoding), err.Error())
				message = string(p.Message)
			}

		} else {
			// is Unknown coding (unicode or reserved)
			message = string(p.Message)
			p.DataCoding = coding.OctetCoding
		}
	}

	return message
}

// Encode text to "dataCoding" encoding
func EncodeMessage(message string, dataCoding coding.DataCoding) []byte {
	switch dataCoding {
	case coding.GSM7BitCoding:
		return coding.EncodeGSM7(message)
	default:
		encoding := dataCoding.Encoding()
		if encoding != nil {
			encodedMessage, err := encoding.NewEncoder().Bytes([]byte(message))
			if err == nil {
				return encodedMessage // encoded to gsm7, ucs2 etc...
			} else {
				logrus.WithFields(logrus.Fields{"worker": "pdu.message"}).Errorf("Can't Encode message with data_coding %d: %s", int(dataCoding), err.Error())
			}
		}
		return []byte(message) // Default utf-8 (coding: 2)
	}
}
