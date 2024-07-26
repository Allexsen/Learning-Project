package models

// Room  represents a group chat
type Room struct {
	ID       int64           `json:"id"`
	Messages []DirectMessage `json:"messages"`
}
