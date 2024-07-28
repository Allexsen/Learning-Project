package messages

// DirectMessage represents a 1-1 chat message
type DirectMessage struct {
	Message
	ChatID int64 `json:"chatID"`
}
