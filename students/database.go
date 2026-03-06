package students

import (
	"IFTP/db"
	"IFTP/utils"
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

func dbRetrieveStudents(ctx context.Context, myDb *db.MyDatabase) ([]Student, error) {
	rows, err := myDb.Pool.Query(ctx,
		"SELECT id, name, email, active FROM students WHERE active = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	students, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[Student])
	if err != nil {
		return nil, err
	}

	return students, nil
}

func dbGetStudentsWithEnrollment(ctx context.Context, myDb *db.MyDatabase) ([]Student, error) {
	rows, err := myDb.Pool.Query(ctx,
		`SELECT s.id, s.name, s.email, s.active, COALESCE(ARRAY_AGG(DISTINCT c.name ORDER BY c.name) FILTER (WHERE c.name IS NOT NULL), '{}') AS enrolled_classes FROM students AS s 
		LEFT JOIN roster AS r on r.student_id = s.id
		LEFT JOIN classes AS c on r.class_id = c.id
		WHERE s.active = true
		GROUP BY s.name, s.id, s.email`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	students, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[Student])
	if err != nil {
		return nil, err
	}

	return students, nil
}

func dbAddStudent(ctx context.Context, myDb *db.MyDatabase, s *Student) error {
	_, err := myDb.Pool.Exec(ctx,
		"INSERT INTO students(name, email, notes, active) VALUES($1, $2, $3, true)",
		s.Name, s.Email, s.Notes)

	return err
}

func dbUpdateStudent(ctx context.Context, myDb *db.MyDatabase, s *Student) (*Student, error) {
	updates := []string{}
	args := []any{}

	if s.Name != "" {
		args = append(args, s.Name)
		updates = append(updates, fmt.Sprintf("name=$%d", len(args)))
	}

	if s.Email != "" {
		args = append(args, s.Email)
		updates = append(updates, fmt.Sprintf("email=$%d", len(args)))
	}

	if s.Notes != "" {
		args = append(args, s.Notes)
		updates = append(updates, fmt.Sprintf("notes=$%d", len(args)))
	}

	if len(updates) == 0 {
		return nil, utils.ErrNoFieldsToUpdate
	}

	args = append(args, s.ID)

	// RETURNING gives you the updated row within one request to the DB
	query := fmt.Sprintf("UPDATE students SET %s WHERE id=$%d AND active=true RETURNING id, name, email, notes, active",
		strings.Join(updates, ", "), len(args))

	rows, err := myDb.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	updatedStudent, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[Student])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("active student with id %d not found", s.ID)
		}
		return nil, err
	}

	return &updatedStudent, nil
}

// func SoftDeleteStudentDB(myDb *db.MyDatabase, id int) (string, error) {

// 	var name string

// 	err := myDb.Db.QueryRow(
// 		"SELECT name FROM students WHERE id=$1 AND active=true", id,
// 	).Scan(&name)

// 	if err == sql.ErrNoRows {
// 		return "", fmt.Errorf("student with id %d not found", id)
// 	}

// 	result, err := myDb.Db.Exec(
// 		"UPDATE students SET active = false WHERE id = $1",
// 		id)

// 	if err != nil {
// 		return "", err
// 	}

// 	// Check if any row was actually updated
// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return "", err
// 	}

// 	if rowsAffected == 0 {
// 		return "", fmt.Errorf("student with id %d not found", id)
// 	}

// 	return name, nil
// }
