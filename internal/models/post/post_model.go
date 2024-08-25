package post

type Post struct {
	ID           int64  `db:"id" json:"id,omitempty"`                         // Unique post id
	UserID       int64  `db:"user_id" json:"user_id,omitempty"`               // User id
	Content      string `db:"content" json:"content,omitempty"`               // Post content
	ParentPostID int64  `db:"parent_post_id" json:"parent_post_id,omitempty"` // Parent Post
	CreatedAt    string `db:"created_at" json:"created_at,omitempty"`         // Post creation time
	LikesCount   int    `db:"likes_count" json:"likes_count,omitempty"`       // Number of likes
}
