package msg

// GroupMessage represents a group chat message
type GroupMessage struct {
	Message
	RoomID int64            `json:"room_id"`
	ReadBy map[int64]string `json:"read_by"`
}
