package devinotele

import (
	"errors"
)

type Code int

const (
	// CodeArgumentCanNotBeNull аргументы не могут быть нулевыми
	CodeArgumentCanNotBeNull Code = 1
	// CodeInvalidArgument невалидные аргументы запроса
	CodeInvalidArgument Code = 2
	// CodeInvalidSessionId невалидный id сессии
	CodeInvalidSessionId Code = 3
	// CodeUnauthorizedAccess ошибка авторизации
	CodeUnauthorizedAccess Code = 4
	// CodeNotEnoughCredits не достаточно средств
	CodeNotEnoughCredits Code = 5
	// CodeInvalidOperation неверная операция
	CodeInvalidOperation Code = 6
	// CodeForbidden Отказано в доступе
	CodeForbidden Code = 7
	// CodeInvalidSenderAddress Неверно введен адрес отправителя
	CodeInvalidSenderAddress Code = 10
	// CodeInvalidReceiverAddress Неверно введен адрес получателя
	CodeInvalidReceiverAddress Code = 11
	// CodeUnacceptableSenderAddress Недопустимый адрес получателя
	CodeUnacceptableSenderAddress Code = 41
	// CodeRejectedBySMSGateway Отклонено смс-центром
	CodeRejectedBySMSGateway Code = 42
	// CodeMessagePastDue Просрочено (истек срок жизни сообщения)
	CodeMessagePastDue Code = 46
	// CodeMessageDeleted Удалено
	CodeMessageDeleted Code = 47
	// CodeRejectedByPlatform Отклонено Платформой
	CodeRejectedByPlatform Code = 48
	// CodeMessageRejected Отклонено
	CodeMessageRejected Code = 69
	// CodeUnknown Неизвестный
	CodeUnknown Code = 99
	// CodeMessageVeryOld *сообщение еще не успело попасть в БД,
	// *сообщение старше 48 часов.
	CodeMessageVeryOld Code = 255
)

var (
	BadRequest              = errors.New("bad request to devino")
	ErrArgumentCanNotBeNull = errors.New("argument cannot be null or empty")
	ErrInvalidArgument      = errors.New("invalid argument")
	ErrInvalidSessionId     = errors.New("invalid session id")
	ErrUnauthorizedAccess   = errors.New("unauthorized access")
	ErrNotEnoughCredits     = errors.New("not enough credits")
	ErrInvalidOperation     = errors.New("invalid operation")
	ErrForbidden            = errors.New("forbidden")

	ErrInvalidSourceAddress        = errors.New("invalid source address")
	ErrInvalidDestinationAddress   = errors.New("invalid destination address")
	ErrUnacceptableReceiverAddress = errors.New("unacceptable destination address")
	ErrRejectedBySMSGateway        = errors.New("rejected by sms getaway")
	ErrMessagePastDue              = errors.New("pass due")
	ErrMessageDeleted              = errors.New("deleted")
	ErrRejectedByPlatform          = errors.New("rejected by platform")
	ErrUnknown                     = errors.New("unknown")
	ErrMessageVeryOld              = errors.New("very old")
	ErrUnknownResponseCode         = errors.New("unknown response code")
)

type ErrorResponse struct {
	Code        Code   `json:"Code"`
	Description string `json:"Desc"`
}

// Err вернет внутреннюю ошибку системы
func (e *ErrorResponse) Err() error {
	switch e.Code {
	case CodeArgumentCanNotBeNull:
		return ErrArgumentCanNotBeNull
	case CodeInvalidArgument:
		return ErrInvalidArgument
	case CodeInvalidSessionId:
		return ErrInvalidSessionId
	case CodeUnauthorizedAccess:
		return ErrUnauthorizedAccess
	case CodeNotEnoughCredits:
		return ErrNotEnoughCredits
	case CodeInvalidOperation:
		return ErrInvalidOperation
	case CodeForbidden:
		return ErrForbidden
	case CodeInvalidSenderAddress:
		return ErrInvalidSourceAddress
	case CodeInvalidReceiverAddress:
		return ErrInvalidDestinationAddress
	case CodeUnacceptableSenderAddress:
		return ErrUnacceptableReceiverAddress
	case CodeRejectedBySMSGateway:
		return ErrRejectedBySMSGateway
	case CodeMessagePastDue:
		return ErrMessagePastDue
	case CodeMessageDeleted:
		return ErrMessageDeleted
	case CodeMessageRejected:
		return ErrMessageDeleted
	case CodeRejectedByPlatform:
		return ErrRejectedByPlatform
	case CodeUnknown:
		return ErrUnknown
	case CodeMessageVeryOld:
		return ErrMessageVeryOld
	default:
		return ErrUnknownResponseCode
	}
}
