package coding

import (
	"fmt"
	. "unicode"

	. "golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/unicode/rangetable"
)

// DataCoding see SMPP v5, section 4.7.7 (123p)

type DataCoding byte

func (c DataCoding) GoString() string {
	return c.String()
}

func (c DataCoding) String() string {
	return fmt.Sprintf("%08b", byte(c))
}

// MessageWaitingInfo ...
func (c DataCoding) MessageWaitingInfo() (coding DataCoding, active bool, kind int) {
	kind = -1
	coding = NoCoding
	switch c >> 4 & 0b1111 {
	case 0b1100:
	case 0b1101:
		coding = GSM7BitCoding
	case 0b1110:
		coding = UCS2Coding
	default:
		return
	}
	active = c>>3 == 1
	kind = int(c & 0b11)
	return
}

// MessageClass ...
func (c DataCoding) MessageClass() (coding DataCoding, class int) {
	class = int(c & 0b11)
	coding = GSM7BitCoding
	if c>>4&0b1111 != 0b1111 {
		coding = NoCoding
		class = -1
	} else if c>>2&0b1 == 1 {
		coding = UCS2Coding
	}
	return
}

// Encoding ...
func (c DataCoding) Encoding() Encoding {
	if coding, _, kind := c.MessageWaitingInfo(); kind != -1 {
		return encodingMap[coding]
	} else if coding, class := c.MessageClass(); class != -1 {
		return encodingMap[coding]
	}

	if enc, ok := encodingMap[c]; ok {
		return enc
	}

	return nil
}

// Splitter ...
func (c DataCoding) Splitter() Splitter {
	if coding, _, kind := c.MessageWaitingInfo(); kind != -1 {
		return splitterMap[coding]
	} else if coding, class := c.MessageClass(); class != -1 {
		return splitterMap[coding]
	}

	return splitterMap[c]
}

// Validate ...
func (c DataCoding) Validate(input string) bool {

	// GSM7
	if c == GSM7BitCoding {
		return IsValidGSM7(input)
	}

	// UCS2
	if c == UCS2Coding {
		return true
	}

	// Others
	if alphabet, ok := alphabetMap[c]; ok {
		for _, r := range input {
			if !Is(alphabet, r) {
				return false
			}
		}

		return true
	}
	return false
}

func GetDataCoding(code int) DataCoding {
	switch code {
	case 0:
		return GSM7BitCoding
	case 1:
		return ASCIICoding
	case 2:
		return OctetCoding
	case 3:
		return Latin1Coding
	case 4:
		return OctetCoding4
	case 5:
		return ShiftJISCoding
	case 6:
		return CyrillicCoding
	case 7:
		return HebrewCoding
	case 8:
		return UCS2Coding
	case 10:
		return ISO2022JPCoding
	case 13:
		return EUCJPCoding
	case 14:
		return EUCKRCoding
	case 191:
		return NoCoding
	default:
		return UCS2Coding
	}
}

const (
	GSM7BitCoding   DataCoding = 0b00000000 // 0 - GSM 7Bit 					(7 bit)
	ASCIICoding     DataCoding = 0b00000001 // 1 - ASCII (ISO-8859-9) 			(8 bit)
	OctetCoding     DataCoding = 0b00000010 // 2 - Octet unspecified 			(8-bit)
	Latin1Coding    DataCoding = 0b00000011 // 3 - Latin-1 (ISO-8859-1) 		(8 bit)
	OctetCoding4    DataCoding = 0b00000100 // 4 - Octet unspecified 			(8-bit)
	ShiftJISCoding  DataCoding = 0b00000101 // 5 - JIS (X 0208-1990)			(16 bit)
	CyrillicCoding  DataCoding = 0b00000110 // 6 - Cyrllic (ISO-8859-5)			(8 bit)
	HebrewCoding    DataCoding = 0b00000111 // 7 - Latin/Hebrew (ISO-8859-8)	(8 bit)
	UCS2Coding      DataCoding = 0b00001000 // 8 - UCS2/UTF-16 (ISO/IEC-10646)	(16 bit)
	ISO2022JPCoding DataCoding = 0b00001010 // 10 - Music Codes	(ISO-2022-JP)
	EUCJPCoding     DataCoding = 0b00001101 // 13 - Extended Kanji JIS (X 0212-1990)
	EUCKRCoding     DataCoding = 0b00001110 // 14 - Korean Graphic Character Set (KS C 5601/KS X 1001)
	NoCoding        DataCoding = 0b10111111 // 15-255 Reserved (Non-specification definition)
)

var splitterMap = map[DataCoding]Splitter{
	GSM7BitCoding:   _7BitSplitter,
	OctetCoding:     _1ByteSplitter,
	ASCIICoding:     _1ByteSplitter,
	Latin1Coding:    _1ByteSplitter,
	OctetCoding4:    _1ByteSplitter,
	ShiftJISCoding:  _MultibyteSplitter,
	CyrillicCoding:  _1ByteSplitter,
	HebrewCoding:    _1ByteSplitter,
	UCS2Coding:      _UTF16Splitter,
	ISO2022JPCoding: _MultibyteSplitter,
	EUCJPCoding:     _MultibyteSplitter,
	EUCKRCoding:     _MultibyteSplitter,
}

var encodingMap = map[DataCoding]Encoding{
	ASCIICoding:     charmap.ISO8859_9,
	Latin1Coding:    charmap.ISO8859_1,
	ShiftJISCoding:  japanese.ShiftJIS,
	CyrillicCoding:  charmap.ISO8859_5,
	HebrewCoding:    charmap.ISO8859_8,
	UCS2Coding:      unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM),
	ISO2022JPCoding: japanese.ISO2022JP,
	EUCJPCoding:     japanese.EUCJP,
	EUCKRCoding:     korean.EUCKR,
}

var alphabetMap = map[DataCoding]*RangeTable{
	ASCIICoding:    _ASCII,
	Latin1Coding:   rangetable.Merge(_ASCII, Latin),
	ShiftJISCoding: rangetable.Merge(_ASCII, _Shift_JIS_Definition),
	CyrillicCoding: rangetable.Merge(_ASCII, Cyrillic),
	HebrewCoding:   rangetable.Merge(_ASCII, Hebrew),
	EUCKRCoding:    rangetable.Merge(_ASCII, _EUC_KR_Definition),
}
