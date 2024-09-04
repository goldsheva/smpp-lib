package pdu

import (
	"bytes"
	"encoding/hex"
	"io"
	"reflect"
)

// PDURequest ...
type PDURequest struct{}

// ReadPDU ...
func ReadPDU(r io.Reader) (interface{}, string, *Header, *PDUError) {
	var pdu interface{}
	var buf bytes.Buffer

	r = io.TeeReader(r, &buf)
	header := new(Header)

	if err := ReadPDUHeader(r, header); err != nil {
		// tcp dump package (Header only)
		return nil, hex.EncodeToString(buf.Bytes()), header, &PDUError{
			Err: err,
		}
	}

	_, err := r.Read(make([]byte, header.CommandLength-16))

	// tcp dump package
	hashPDU := hex.EncodeToString(buf.Bytes())

	if err != nil {
		return nil, hashPDU, header, &PDUError{
			Err:           err,
			CommandStatus: ESME_RINVCMDLEN, // 2 - Invalid Command Length
		}
	}

	if t, ok := Types[header.CommandID]; !ok {
		return nil, hashPDU, header, &PDUError{
			CommandStatus: ESME_RINVCMDID, // 3 - Invalid Command ID
		}

	} else {
		pdu = reflect.New(t).Interface()
		if _, err := UnmarshalPDU(&buf, pdu); err != nil {
			return nil, hashPDU, header, &PDUError{
				Err:           err,
				CommandStatus: ESME_RINVCMDLEN, // 2 - Invalid Command Length
			}
		}
	}

	return pdu, hashPDU, header, nil
}
