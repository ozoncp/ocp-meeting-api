package producer

type EventType string

const (
	Create EventType = "Created"
	Update EventType = "Updated"
	Delete EventType = "Deleted"
)

type EventMessage struct {
	MeetingId uint64
	Timestamp int64
}

type Message struct {
	EventType    EventType
	EventMessage EventMessage
}

func CreateMessage(eventType EventType, eventMessage EventMessage) Message {
	return Message{
		EventType:    eventType,
		EventMessage: eventMessage,
	}
}
