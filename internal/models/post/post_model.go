package post

import (
	"database/sql"
	"log"

	"github.com/Allexsen/Learning-Project/internal/models/common"
	"github.com/google/uuid"
)

// Post represents a post model. It is used to store post data in the database.
type Post struct {
	ID           int64  `db:"id" json:"id,omitempty"`                         // Unique post id
	UserID       int64  `db:"user_id" json:"user_id,omitempty"`               // User id
	Content      string `db:"content" json:"content,omitempty"`               // Post content
	ParentPostID int64  `db:"parent_post_id" json:"parent_post_id,omitempty"` // Parent Post
	CreatedAt    string `db:"created_at" json:"created_at,omitempty"`         // Post creation time
	LikesCount   int    `db:"likes_count" json:"likes_count,omitempty"`       // Number of likes
}

// New creates a new post. It returns a pointer to the newly created post.
func New(uid int64, content string, parentPostID int64) *Post {
	return &Post{
		ID:           int64(uuid.New().ID()),
		UserID:       uid,
		Content:      content,
		ParentPostID: parentPostID,
	}
}

// AddPost adds a new post to the database. It returns the id of the newly created post.
func (p *Post) AddPost(db *sql.DB) (int64, error) {
	log.Printf("[POST] Adding post %d to the database", p.ID)

	q := `INSERT INTO practice_db.posts (id, user_id, content, parent_post_id)
		VALUES(?, ?, ?, ?)`
	result, err := db.Exec(q, p.ID, p.UserID, p.Content, p.ParentPostID)
	if err != nil {
		return -1, common.GetQueryError(q, "Couldn't add a new post", p, err)
	}

	id, err := common.GetLastInsertId(result, q, p)
	if err != nil {
		return -1, err
	}

	return id, nil
}

// RetrievePostByID retrieves post by post id. It returns an error if the post is not found.
func (p *Post) RetrievePostByID(db *sql.DB) error {
	log.Printf("[POST] Retrieving post by id %d from the database", p.ID)

	q := `SELECT user_id, content, parent_post_id, created_at, likes_count
		FROM practice_db.posts
		WHERE id=?`
	err := db.QueryRow(q, p.ID).Scan(&p.UserID, &p.Content, &p.ParentPostID, &p.CreatedAt, &p.LikesCount)
	return common.GetQueryError(q, "Couldn't retrieve post by id", p, err)
}

// UpdatePostByID updates post by post id. It returns an error if the post is not found.
func (p *Post) DeletePostByID(db *sql.DB) error {
	log.Printf("[POST] Deleting post by id %d from the database", p.ID)

	q := `DELETE FROM practice_db.posts
		WHERE id=?`
	_, err := db.Exec(q, p.ID)
	return common.GetQueryError(q, "Couldn't delete post by id", p, err)
}

// RetrieveAllPostsByUserID retrieves all posts by user id. It returns an error if the user is not found.
func (p *Post) RetrieveAllPostsByUserID(db *sql.DB) ([]Post, error) {
	log.Printf("[POST] Retrieving all posts by user id %d from the database", p.UserID)

	q := `SELECT id, content, parent_post_id, created_at, likes_count
		FROM practice_db.posts
		WHERE user_id=?`
	rows, err := db.Query(q, p.UserID)
	if err != nil {
		return nil, common.GetQueryError(q, "Couldn't retrieve all posts by user id", p, err)
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.ID, &post.Content, &post.ParentPostID, &post.CreatedAt, &post.LikesCount)
		if err != nil {
			return nil, common.GetQueryError(q, "Couldn't retrieve all posts by user id", p, err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// RetrievePostRepliesByID retrieves all replies to post id. It returns an error if the post is not found.
func (p *Post) RetrievePostRepliesByID(db *sql.DB) ([]Post, error) {
	log.Printf("[POST] Retrieving all replies to post id %d from the database", p.ID)

	q := `SELECT id, user_id, content, parent_post_id, created_at, likes_count
		FROM practice_db.posts
		WHERE parent_post_id=?`
	rows, err := db.Query(q, p.ID)
	if err != nil {
		return nil, common.GetQueryError(q, "Couldn't retrieve all replies to post id", p, err)
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.ID, &post.UserID, &post.Content, &post.ParentPostID, &post.CreatedAt, &post.LikesCount)
		if err != nil {
			return nil, common.GetQueryError(q, "Couldn't retrieve all replies to post id", p, err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}
