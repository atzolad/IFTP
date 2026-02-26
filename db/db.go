package db

import (
	"html/template"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MyDatabase struct {
	Pool      *pgxpool.Pool
	Logger    *log.Logger
	Templates *template.Template
}
