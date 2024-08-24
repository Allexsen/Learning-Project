package comment

type Comment struct {
	ID        int64  `db:"id" json:"id,omitempty"`                // Unique comment id
	PostID    int64  `db:"post_id" json:"postId,omitempty"`       // Post id
	UserID    int64  `db:"user_id" json:"userId,omitempty"`       // User id
	Content   string `db:"content" json:"content,omitempty"`      // Comment content
	CreatedAt string `db:"created_at" json:"createdAt,omitempty"` // Comment creation time
}
