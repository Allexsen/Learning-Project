package msg

// BaseMessage is the base message model
// that all message types will inherit from
type BaseMessage struct {
	ID        int64  `json:"id,omitempty"`
	SenderID  int64  `json:"sender_id,omitempty"`
	Content   string `json:"content,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Status    string `json:"status,omitempty"`
}

// Message is the interface that all message types will implement
type Message interface {
	GetSenderID() int64
	GetContent() string
}

// GetSenderID returns the sender ID
func (msg BaseMessage) GetSenderID() int64 {
	return msg.SenderID
}

// GetContent returns the message content
func (msg BaseMessage) GetContent() string {
	return msg.Content
}
