package msg

type BaseMessage struct {
	ID        int64  `json:"id"`
	SenderID  int64  `json:"sender_id"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
	Status    string `json:"status"`
}

type Message interface {
	GetSenderID() int64
	GetContent() string
}

func (msg BaseMessage) GetSenderID() int64 {
	return msg.SenderID
}

func (msg BaseMessage) GetContent() string {
	return msg.Content
}
