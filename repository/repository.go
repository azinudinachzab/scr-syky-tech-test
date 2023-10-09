package repository

import "database/sql"

type PgRepository struct {
	dbCore *sql.DB
}

func NewPgRepository(dbCore *sql.DB) *PgRepository {
	return &PgRepository{
		dbCore: dbCore,
	}
}

type Repository interface {}
