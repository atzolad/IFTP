package student

import (
	"IFTP/db"
)

func RetrieveStudents(myDb *db.MyDatabase) ([]Student, error) {
	rows, err := myDb.Db.Query(
		"SELECT id, name, email, paid, active FROM students")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// A students slice to hold the data from the returned rows
	var students []Student

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var student Student
		if err := rows.Scan(&student.ID, &student.Name, &student.Email,
			&student.Paid, &student.Active); err != nil {
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
		"INSERT INTO students(name, email, paid, active) VALUES($1, $2, $3, $4) RETURNING id",
		s.Name, s.Email, s.Paid, s.Active,
	).Scan(&s.ID)

	return err
}

// // func (myDb *myDatabase) SoftDeleteStudent(student student) {
// // 	"UPDATE students(id, name, email, active) VALUES( ?, ?, ?, ?)"
// // }
