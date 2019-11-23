package devinotele

type State int

const (
	// StateSentToMobileNetwork Отправлено (передано в мобильную сеть)
	StateSentToMobileNetwork State = -1
	// StateQueued В очереди
	StateQueued State = -2
	// StateStopped Остановлено
	StateStopped State = -98
	// StateDelivered Доставлено абоненту
	StateDelivered State = 0
	// StateInvalidSenderAddress Неверно введен адрес отправителя
	StateInvalidSenderAddress State = 10
	// StateInvalidReceiverAddress Неверно введен адрес получателя
	StateInvalidReceiverAddress State = 11
	// StateUnacceptableSenderAddress Недопустимый адрес получателя
	StateUnacceptableSenderAddress State = 41
	// StateRejectedBySMSGateway Отклонено смс-центром
	StateRejectedBySMSGateway State = 42
	// StatePastDue Просрочено (истек срок жизни сообщения)
	StatePastDue State = 46
	// StateDeleted Удалено
	StateDeleted State = 47
	// StateRejectedByPlatform Отклонено Платформой
	StateRejectedByPlatform State = 48
	// StateRejected Отклонено
	StateRejected State = 69
	// StateUnknown Неизвестный
	StateUnknown State = 99
	// StateVeryOld *сообщение еще не успело попасть в БД,
	// *сообщение старше 48 часов.
	StateVeryOld State = 255
)

type MessageInfo struct {
	State            State  `json:"State"`
	StateDescription string `json:"StateDescription"`
}

// {
//   "State":<Код статуса сообщения>,
//   "CreationDateUtc":<Дата создания>,
//   "SubmittedDateUtc":<Дата отправки сообщения>,
//   "ReportedDateUtc":<Дата доставки сообщения>,
//   "TimeStampUtc":"<Дата и время получения отчета>",
//   "StateDescription":"<Описание статуса>",
//   "Price":<Стоимость>
// }
