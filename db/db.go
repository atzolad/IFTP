package db

import (
	"database/sql"
	"html/template"
	"log"
)

type MyDatabase struct {
	Db        *sql.DB
	Logger    *log.Logger
	Templates *template.Template
}
