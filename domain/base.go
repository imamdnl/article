package domain

import "github.com/jackc/pgx/v4/pgxpool"

type WebResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

//database
type BaseCapsule struct {
	Database *pgxpool.Pool
}

func NewBaseRepository(dbp *pgxpool.Pool) *BaseCapsule {
	return &BaseCapsule{
		Database: dbp,
	}
}
