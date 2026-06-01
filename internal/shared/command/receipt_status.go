package command

type ReceiptStatus string

func (s ReceiptStatus) String() string {
	return string(s)
}

const (
	StatusSucceeded ReceiptStatus = "succeeded"
	StatusFailed    ReceiptStatus = "failed"
)
