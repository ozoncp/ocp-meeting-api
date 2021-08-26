package producer

type EventType int

const (
	Create EventType = iota
	Update
	Delete
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

func (eventType EventType) String() string {
	switch eventType {
	case Create:
		return "Created"
	case Update:
		return "Updated"
	case Delete:
		return "Removed"
	default:
		return "Unknown event type"
	}
}
