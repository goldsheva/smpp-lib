package coding

import (
	"os"
	"strconv"
)

func BestAllCoding(input string, isGSM7Supported bool) DataCoding {
	var codings []DataCoding
	if isGSM7Supported {
		codings = []DataCoding{
			GSM7BitCoding, ASCIICoding, Latin1Coding,
			CyrillicCoding, HebrewCoding, ShiftJISCoding,
			EUCKRCoding, UCS2Coding,
		}
	} else {
		codings = []DataCoding{
			ASCIICoding, Latin1Coding,
			CyrillicCoding, HebrewCoding, ShiftJISCoding,
			EUCKRCoding, UCS2Coding,
		}
	}

	for _, coding := range codings {
		if coding.Validate(input) {
			return coding
		}
	}
	return UCS2Coding
}

func BestSafeCoding(input string, isGSM7Supported bool) DataCoding {
	// Is GSM7 (7bit)
	if isGSM7Supported && GSM7BitCoding.Validate(input) {
		return GSM7BitCoding
	}

	// Is UCS2 (16bit)
	for _, r := range input {
		if r > 0xFF { // UCS2 for symbols exclude in Latin-1 (over 255)
			return UCS2Coding
		}
	}

	// Is Unicode (8bit)
	return OctetCoding
}

func BestCoding(input string, isGSM7Supported bool) DataCoding {
	isSafeCoding, _ := strconv.ParseBool(os.Getenv("IS_SAFE_CODING"))

	if isSafeCoding {
		return BestSafeCoding(input, isGSM7Supported)
	} else {
		return BestAllCoding(input, isGSM7Supported)
	}
}
