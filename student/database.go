package student

import (
	"IFTP/db"
	"database/sql"
	"fmt"
)

func RetrieveStudents(myDb *db.MyDatabase) ([]Student, error) {
	rows, err := myDb.Db.Query(
		"SELECT id, name, email, active FROM students WHERE active = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// A students slice to hold the data from the returned rows
	var students []Student

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var student Student
		if err := rows.Scan(&student.ID, &student.Name, &student.Email, &student.Active); err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	if err = rows.Err(); err != nil {
		return students, err
	}
	return students, nil
}

func InsertStudent(myDb *db.MyDatabase, s *Student) error {
	err := myDb.Db.QueryRow(
		"INSERT INTO students(name, email) VALUES($1, $2) RETURNING id",
		s.Name, s.Email).Scan(&s.ID)

	return err
}

func UpdateStudentDB(myDb *db.MyDatabase, id int, s *Student) error {
	var current Student

	err := myDb.Db.QueryRow(
		"SELECT id, name, email, active FROM students WHERE id=$1", id,
	).Scan(&current.ID, &current.Name, &current.Email, &current.Active)

	if err == sql.ErrNoRows {
		return fmt.Errorf("student with id %d not found", id)
	}

	if err != nil {
		return err
	}

	// Check if these parameters are set to update.
	if s.Name != "" {
		current.Name = s.Name
	}

	if s.Email != "" {
		current.Email = s.Email
	}

	current.Active = s.Active

	result, err := myDb.Db.Exec(
		"UPDATE students SET name = $1, email = $2, WHERE id = $3",
		current.Name, current.Email, id)

	if err != nil {
		return err
	}

	// Check if any row was actually updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("student with id %d not found", id)
	}

	return nil
}

func SoftDeleteStudentDB(myDb *db.MyDatabase, id int) (string, error) {

	var name string

	err := myDb.Db.QueryRow(
		"SELECT name FROM students WHERE id=$1 AND active=true", id,
	).Scan(&name)

	if err == sql.ErrNoRows {
		return "", fmt.Errorf("student with id %d not found", id)
	}

	result, err := myDb.Db.Exec(
		"UPDATE students SET active = false WHERE id = $1",
		id)

	if err != nil {
		return "", err
	}

	// Check if any row was actually updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", err
	}

	if rowsAffected == 0 {
		return "", fmt.Errorf("student with id %d not found", id)
	}

	return name, nil
}
