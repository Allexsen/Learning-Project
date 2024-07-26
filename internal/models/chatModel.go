package models

// Chat represents a 1-1 chat
type Chat struct {
	ID       int64           `json:"id"`
	Messages []DirectMessage `json:"messages"`
}
