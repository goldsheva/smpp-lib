package pdu

import (
	"fmt"
	"regexp"
	"time"
)

type DeliveryReceipt struct {
	ID         string `json:"id"`
	Sub        string `json:"sub"`
	Dlvrd      string `json:"dlvrd"`
	SubmitDate string `json:"submit_date"`
	DoneDate   string `json:"done_date"`
	Status     string `json:"stat"`
	Error      string `json:"err"`
	Text       string `json:"text"`
}

// Convert Milliseconds to DLR format seconds
func ConvertMillisecondsToDLRDate(sentTime int64) string {
	sec := sentTime / 1000

	t := time.Unix(sec, 0).UTC()

	return t.Format("0201061504")
}

func (dlr *DeliveryReceipt) GenerateDLRString() string {
	return fmt.Sprintf("id:%s sub:%s dlvrd:%s submit date:%s done date:%s stat:%s err:%s text:%s",
		dlr.ID, dlr.Sub, dlr.Dlvrd, dlr.SubmitDate, dlr.DoneDate, dlr.Status, dlr.Error, dlr.Text)
}

// ParseDLR parse short_message to struct DeliveryReceipt.
func ParseDLR(dlrMessage string) (*DeliveryReceipt, error) {
	re := regexp.MustCompile(`id:(\S+) sub:(\S+) dlvrd:(\S+) submit date:(\S+) done date:(\S+) stat:(\S+) err:(\S+)(?: text:(.*))?`)

	matches := re.FindStringSubmatch(dlrMessage)
	if matches == nil {
		return nil, fmt.Errorf("invalid DLR format")
	}

	var parsedText = ""
	if len(matches) > 8 {
		parsedText = matches[8]
	}

	dlr := &DeliveryReceipt{
		ID:         matches[1],
		Sub:        matches[2],
		Dlvrd:      matches[3],
		SubmitDate: matches[4],
		DoneDate:   matches[5],
		Status:     matches[6],
		Error:      matches[7],
		Text:       parsedText,
	}

	if dlr.ID == "" {
		return nil, fmt.Errorf("invalid MessageID")
	}

	return dlr, nil
}
