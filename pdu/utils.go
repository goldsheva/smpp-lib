package pdu

import (
	"bufio"
	"bytes"
	"fmt"
)

func readCString(buf *bufio.Reader) (string, error) {
	value, err := buf.ReadString(0)
	if err == nil {
		value = value[0 : len(value)-1]
	}
	return value, err
}

func writeCString(buf *bytes.Buffer, value string) {
	buf.WriteString(value)
	buf.WriteByte(0)
}

func getBool(v bool) byte {
	if v {
		return 1
	}
	return 0
}

func isNumeric(s string) bool {
	_, err := fmt.Sscanf(s, "%d", &s)
	return err == nil
}

func Сontains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
