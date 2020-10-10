package devinotele

type ReasonKey string

const (
	ReasonKeyBillingError    ReasonKey = "billing.error"
	ReasonKeyForbidden       ReasonKey = "forbidden"
	ReasonKeyUnknown         ReasonKey = "unknown"
	ReasonKeyInvalid         ReasonKey = "invalid"
	ReasonKeyLengthTooLong   ReasonKey = "length.too.long"
	ReasonKeyMustBeNotNull   ReasonKey = "must.be.not.null"
	ReasonKeyTooManyMessages ReasonKey = "too.many.messages"
)

type StatusCode string

const (
	StatusCodeOK       StatusCode = "OK"
	StatusCodeRejected StatusCode = "REJECTED"
)

type Reason struct {
	Key ReasonKey `json:"key"`
	Ref string    `json:"ref"`
}

type ErrorResponse struct {
	Description string
	Reasons     []Reason
}

func (e ErrorResponse) Error() string {
	return "err: " + e.Description
}
