package pdu

const (
	TypeOfNumberUnknown         byte = 0x00
	TypeOfNumberInternational   byte = 0x01
	TypeOfNumberNational        byte = 0x02
	TypeOfNumberNetworkSpecific byte = 0x03
	TypeOfNumberSubscriber      byte = 0x04
	TypeOfNumberAlphanumeric    byte = 0x05
	TypeOfNumberAbbreviated     byte = 0x06
	TypeOfNumberExtension       byte = 0x07

	NumberingPlanUnknown   byte = 0x00
	NumberingPlanE164      byte = 0x01
	NumberingPlanX121      byte = 0x03
	NumberingPlanTelex     byte = 0x04
	NumberingPlanNational  byte = 0x08
	NumberingPlanPrivate   byte = 0x09
	NumberingPlanERMES     byte = 0x0A
	NumberingPlanExtension byte = 0x0F
)
