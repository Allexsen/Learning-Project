package messages

// GroupMessage represents a group chat message
type GroupMessage struct {
	Message
	RoomID int64            `json:"roomID"`
	ReadBy map[int64]string `json:"readBy"`
}
