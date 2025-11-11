package user

import "database/sql"

type myDatabase struct {
	db    *sql.DB
	hello int
}

func (myDb *myDatabase) AddUser(user user) {
	"INSERT INTO users(id, name, email) VALUES( ?, ?, ?)";
}

func (myDb *myDatabase) SoftDeleteUser(user user) {
	"UPDATE users(id, name, email, active) VALUES( ?, ?, ?, ?)"
}


