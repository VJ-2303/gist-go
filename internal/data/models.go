package data

import "database/sql"

type Models struct {
	User      UserModel
	Post      PostModel
	ShareLink ShareLinkModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		User:      UserModel{DB: db},
		Post:      PostModel{DB: db},
		ShareLink: ShareLinkModel{DB: db},
	}
}
