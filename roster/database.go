package roster

import (
	"IFTP/db"
	"time"
)

func dbEnroll(myDb *db.MyDatabase, classID int, classDate time.Time, studentID int) error {

	var rosterID int

	err := myDb.Db.QueryRow(
		"INSERT INTO roster (class_date, student_id, class_id, registration_date) "+
			"VALUES ($1, $2, $3, NOW()) RETURNING id", classDate, studentID, classID).Scan(&rosterID)

	return err
}
