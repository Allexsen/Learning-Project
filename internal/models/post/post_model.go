package post

type Post struct {
	ID         int64  `db:"id" json:"id,omitempty"`                  // Unique post id
	UserID     int64  `db:"user_id" json:"userId,omitempty"`         // User id
	Content    string `db:"content" json:"content,omitempty"`        // Post content
	CreatedAt  string `db:"created_at" json:"createdAt,omitempty"`   // Post creation time
	LikesCount int    `db:"likes_count" json:"likesCount,omitempty"` // Number of likes
}
