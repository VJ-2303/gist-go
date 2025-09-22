package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/vj-2303/gist-go/internal/validator"
)

var (
	ErrPostNotFound = errors.New("posts not found")
)

type Post struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"-"`
	Title     string    `json:"title"`
	Language  string    `json:"language"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	Version   int       `json:"version"`
}

func ValidatePosts(v *validator.Validator, post *Post) {
	v.Check(post.Title != "", "title", "must be provided")
	v.Check(len(post.Title) <= 500, "title", "must be less than 500 bytes long")
	v.Check(post.Language != "", "language", "must be provided")
	v.Check(len(post.Language) <= 15, "language", "must be less than 15 bytes Long")
	v.Check(post.Code != "", "code", "must be provided")
	v.Check(len(post.Code) >= 2, "code", "must be greater than 50 bytes Long")
}

type PostModel struct {
	DB *sql.DB
}

func (m PostModel) GetByID(postID, userID int64) (*Post, error) {

	query := `
		SELECT id,user_id,title,language,code,created_at,version
		FROM posts
		WHERE id = $1 AND user_id = $2
			 `

	var post Post

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, postID, userID).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Language,
		&post.Code,
		&post.CreatedAt,
		&post.Version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPostNotFound
		} else {
			return nil, err
		}
	}
	return &post, nil
}

func (m PostModel) Insert(post *Post) error {

	query := `
		INSERT INTO posts(user_id, title,language,code)
		VALUES ($1,$2,$3,$4)
		RETURNING id, created_at, version
			 `
	args := []any{post.UserID, post.Title, post.Language, post.Code}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.Version,
	)
	if err != nil {
		return err
	}
	return nil
}
