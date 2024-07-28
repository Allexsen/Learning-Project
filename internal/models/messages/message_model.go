package messages

import "time"

type Message struct {
	ID        int64     `json:"id"`
	Sender    int64     `json:"sender"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}
