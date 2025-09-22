package data

import (
	"context"
	"database/sql"
	"time"
)

type ShareLink struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

type ShareLinkModel struct {
	DB *sql.DB
}

func (m ShareLinkModel) Insert(sharelink *ShareLink) error {

	query := `
		INSERT INTO share_links(post_id,token)
		VALUES ($1,$2)
		RETURNING id,created_at
			 `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, sharelink.PostID, sharelink.Token).Scan(
		&sharelink.ID,
		&sharelink.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m ShareLinkModel) GetByShareToken(token string) (*ShareLink, error) {

	query := `
		SELECT id, post_id, token, created_at
		FROM share_links
		WHERE token = $1
			 `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var shareLink ShareLink

	err := m.DB.QueryRowContext(ctx, query, token).Scan(
		&shareLink.ID,
		&shareLink.PostID,
		&shareLink.Token,
		&shareLink.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &shareLink, nil
}
