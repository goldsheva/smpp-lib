package pdu

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"reflect"
	"strconv"
)

// UnmarshalPDU ...
func UnmarshalPDU(r io.Reader, packet interface{}) (n int64, err error) {
	buf := bufio.NewReader(r)
	v := reflect.ValueOf(packet)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		switch field := v.Field(i); field.Kind() {
		case reflect.String:
			var value string
			if value, err = readCString(buf); err == nil {
				field.SetString(value)
			}
		case reflect.Uint8:
			var value byte
			if value, err = buf.ReadByte(); err == nil {
				field.SetUint(uint64(value))
			}
		case reflect.Bool:
			var value byte
			if value, err = buf.ReadByte(); err == nil {
				field.SetBool(value == 1)
			}
		case reflect.Array, reflect.Map, reflect.Slice, reflect.Struct:
			switch v := (field.Addr().Interface()).(type) {
			case *Header:
				err = ReadPDUHeader(buf, v)
				if v.CommandStatus != 0 {
					return
				}
			case io.ByteWriter:
				var value byte
				if value, err = buf.ReadByte(); err == nil {
					err = v.WriteByte(value)
				}
			case io.ReaderFrom:
				if m, ok := v.(*ShortMessage); ok {
					m.Prepare(packet)
				}
				_, err = v.ReadFrom(buf)
			}
		}
		n = int64(buf.Size())
		if err != nil {
			err = errors.New("ErrUnmarshalPDUFailed")
			return
		}
	}
	return
}

// MarshalPDU ...
func MarshalPDU(w io.Writer, packet interface{}) (string, *PDUError) {
	var buf bytes.Buffer
	p := reflect.ValueOf(packet)

	if p.Kind() == reflect.Ptr {
		p = p.Elem()
	}

	for i := 0; i < p.NumField(); i++ {
		field := p.Field(i)
		switch field.Kind() {
		case reflect.String:
			writeCString(&buf, field.String())

		case reflect.Uint8:
			buf.WriteByte(byte(field.Uint()))

		case reflect.Bool:
			var value byte
			if field.Bool() {
				value = 1
			}
			buf.WriteByte(value)

		case reflect.Array, reflect.Map, reflect.Slice, reflect.Struct:
			switch v := field.Addr().Interface().(type) {
			case *Header:
				if value := p.Type().Field(i).Tag.Get(_ID); value != "" {

					parsed, err := strconv.ParseUint(value, 16, 32)
					if err != nil {
						return "", &PDUError{
							Err: err,
						}
					}

					v.CommandID = CommandID(parsed)
				}

				if v.Sequence <= 0 {
					return "", &PDUError{
						CommandStatus: ESME_RUNKNOWNERR, // 255 - Unknown Error
					}
				}

				_ = binary.Write(&buf, binary.BigEndian, v)

				if v.CommandStatus != ESME_ROK {
					goto write
				}

			case io.ByteReader:

				value, err := v.ReadByte()
				if err != nil {
					return "", &PDUError{
						Err: err,
					}
				}
				buf.WriteByte(value)

			case io.WriterTo:
				if m, ok := v.(*ShortMessage); ok {
					m.Prepare(packet)
				}

				_, err := v.WriteTo(&buf)
				if err != nil {
					return "", &PDUError{
						Err: err,
					}
				}
			}
		}
	}
write:
	if p.Field(0).Type() == reflect.TypeOf(Header{}) {
		data := buf.Bytes()
		binary.BigEndian.PutUint32(data[0:4], uint32(buf.Len()))
	}

	// tcp dump package
	hashPDU := hex.EncodeToString(buf.Bytes())

	_, err := buf.WriteTo(w)
	if err != nil {
		return "", &PDUError{
			Err: err,
		}
	}

	return hashPDU, nil
}
