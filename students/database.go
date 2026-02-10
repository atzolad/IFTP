package students

import (
	"IFTP/db"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

func dbRetrieveStudents(myDb *db.MyDatabase) ([]Student, error) {
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

func dbGetStudentsWithEnrollment(myDb *db.MyDatabase) ([]Student, error) {
	rows, err := myDb.Db.Query(
		`SELECT s.id, s.name, s.email, s.active, ARRAY_AGG(DISTINCT c.name ORDER BY c.name) AS enrolledClasses FROM students AS s 
		JOIN roster AS r on r.student_id = s.id
		JOIN classes AS c on r.class_id = c.id
		WHERE s.active = true
		GROUP BY s.name, s.id, s.email`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// A students slice to hold the data from the returned rows
	var students []Student

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var student Student
		if err := rows.Scan(&student.ID, &student.Name, &student.Email, &student.Active, (*pq.StringArray)(&student.EnrolledClasses)); err != nil {
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

func UpdateStudentDB(myDb *db.MyDatabase, id int, s *Student) (*Student, error) {
	updates := []string{}
	args := []any{}
	argCount := 1

	if s.Name != "" {
		updates = append(updates, fmt.Sprintf("name=$%d", argCount))
		args = append(args, s.Name)
		argCount++
	}

	if s.Email != "" {
		updates = append(updates, fmt.Sprintf("email=$%d", argCount))
		args = append(args, s.Email)
		argCount++
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	args = append(args, id)
	// RETURNING gives you the updated row within one request to the DB
	query := fmt.Sprintf("UPDATE students SET %s WHERE id=$%d AND active=true RETURNING id, name, email, active",
		strings.Join(updates, ", "), argCount)

	var updated Student
	err := myDb.Db.QueryRow(query, args...).Scan(&updated.ID, &updated.Name, &updated.Email, &updated.Active)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("active student with id %d not found", id)
	}
	if err != nil {
		return nil, err
	}

	return &updated, nil
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
