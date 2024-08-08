package msg

// GroupMessage represents a group chat message
type GroupMessage struct {
	BaseMessage
	RoomID int64            `json:"room_id,omitempty"`
	ReadBy map[int64]string `json:"read_by,omitempty"`
}
