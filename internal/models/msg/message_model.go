package msg

type Message struct {
	ID        int64  `json:"id"`
	SenderID  int64  `json:"sender_id"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
	Status    string `json:"status"`
}
