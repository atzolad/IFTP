package student

import (
	"IFTP/db"
	"fmt"
)

func RetrieveStudents(myDb *db.MyDatabase) ([]Student, error) {
	rows, err := myDb.Db.Query(
		"SELECT id, name, email FROM students WHERE active = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// A students slice to hold the data from the returned rows
	var students []Student

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var student Student
		if err := rows.Scan(&student.ID, &student.Name, &student.Email); err != nil {
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
	result, err := myDb.Db.Exec(
		"UPDATE students SET name = $1, email = $2, WHERE id = $3",
		s.Name, s.Email, s.ID)

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

func UpdatedStudentDB(myDb *db.MyDatabase, id int, s *Student) error {
	result, err := myDb.Db.Exec(
		"UPDATE students SET active = $1 WHERE id = $2",
		s.Name, s.Email, id,
	)
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

// // func (myDb *myDatabase) SoftDeleteStudent(student student) {
// // 	"UPDATE students(id, name, email, active) VALUES( ?, ?, ?, ?)"
// // }
