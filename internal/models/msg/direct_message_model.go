package msg

// DirectMessage represents a 1-1 chat message
type DirectMessage struct {
	BaseMessage
	ChatID int64 `json:"chat_id"`
}
