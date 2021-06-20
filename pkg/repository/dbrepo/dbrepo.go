package dbrepo

import (
	"database/sql"

	"github.com/fangjjcs/bookings-app/pkg/config"
	"github.com/fangjjcs/bookings-app/pkg/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo{
	return &postgresDBRepo{
		App: a,
		DB: conn,
	}
}