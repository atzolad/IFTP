package class

import (
	"IFTP/db"
	"fmt"

	"github.com/lib/pq"
)

func dbListClasses(myDb *db.MyDatabase) ([]Class, error) {
	rows, err := myDb.Db.Query(
		`SELECT c.id, name, teacher, day_of_week, time, description, capacity, ARRAY_AGG(DISTINCT cs.session_date ORDER BY cs.session_date) AS session_dates, COUNT(DISTINCT r.student_id) AS enrolled_count
		FROM classes AS c
		JOIN class_schedule AS cs ON cs.class_id = c.id
		JOIN roster AS r ON r.class_id = c.id
		WHERE active = True
		GROUP BY c.id`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// A Classes slice to hold the data from the returned rows
	var classes []Class

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var class Class
		if err := rows.Scan(&class.ID, &class.Name, &class.Teacher, &class.DayOfWeek, &class.Time,
			&class.Description, &class.Capacity, (*pq.StringArray)(&class.SessionDates), &class.EnrolledCount); err != nil {
			return nil, err
		}
		classes = append(classes, class)
	}
	if err = rows.Err(); err != nil {
		return classes, err
	}
	return classes, nil
}

func dbListClassesByMonth(myDb *db.MyDatabase, month string, studentId ...int) ([]Class, error) {
	queryStmt := `SELECT c.id, name, teacher, day_of_week, time, description, capacity, ARRAY_AGG(DISTINCT cs.session_date ORDER BY cs.session_date) AS session_dates, COUNT(DISTINCT r.student_id) AS enrolled_count
		FROM classes AS c
		JOIN class_schedule AS cs ON cs.class_id = c.id
		JOIN roster AS r ON r.class_id = c.id
		WHERE active = True`

	var args []any

	if month != "" {
		fmt.Printf("Month: %v", month)
		args = append(args, month)
		queryStmt = queryStmt + fmt.Sprintf(" AND month = $%d ", len(args))
	}
	if len(studentId) > 0 {
		fmt.Printf("student id: %v ", studentId[0])
		args = append(args, studentId[0])
		queryStmt = queryStmt + fmt.Sprintf(" AND r.student_id = $%d ", len(args))
	}

	queryStmt = queryStmt + " GROUP BY c.id "
	// var rows *sql.Rows
	// var err error

	rows, err := myDb.Db.Query(queryStmt, args...)

	// if len(args) > 0 {
	// 	rows, err = myDb.Db.Query(queryStmt, args[0])
	// } else {
	// 	rows, err = myDb.Db.Query(queryStmt)
	// }

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// A Classes slice to hold the data from the returned rows
	var classes []Class

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var class Class
		if err := rows.Scan(&class.ID, &class.Name, &class.Teacher, &class.DayOfWeek, &class.Time,
			&class.Description, &class.Capacity, (*pq.StringArray)(&class.SessionDates), &class.EnrolledCount); err != nil {
			return nil, err
		}
		fmt.Println(classes)
		classes = append(classes, class)
	}
	if err = rows.Err(); err != nil {
		return classes, err
	}
	return classes, nil
}

// func DbCreateClass(myDb *db.MyDatabase, c *Class) error {
// 	err := myDb.Db.QueryRow(
// 		"INSERT INTO classes (name, teacher, day, time) VALUES($1, $2, $3, $4) RETURNING id",
// 		c.Name, c.Teacher, c.Day, c.Time).Scan(&c.ID)

// 	return err
// }

// func UpdateClassDB(myDb *db.MyDatabase, id int, c *Class) (*Class, error) {
// 	updates := []string{}
// 	args := []any{}
// 	argCount := 1

// 	if c.Name != "" {
// 		updates = append(updates, fmt.Sprintf("name=$%d", argCount))
// 		args = append(args, c.Name)
// 		argCount++
// 	}

// 	if c.Teacher != "" {
// 		updates = append(updates, fmt.Sprintf("teacher=$%d", argCount))
// 		args = append(args, c.Teacher)
// 		argCount++
// 	}

// 	if c.Day != "" {
// 		updates = append(updates, fmt.Sprintf("day=$%d", argCount))
// 		args = append(args, c.Day)
// 		argCount++
// 	}

// 	if c.Time != "" {
// 		updates = append(updates, fmt.Sprintf("time=$%d", argCount))
// 		args = append(args, c.Time)
// 		argCount++
// 	}

// 	if len(updates) == 0 {
// 		return nil, fmt.Errorf("no fields to update")
// 	}

// 	args = append(args, id)
// 	// RETURNING gives you the updated row within one request to the DB
// 	query := fmt.Sprintf("UPDATE classes SET %s WHERE id=$%d RETURNING id, name, teacher, day, time",
// 		strings.Join(updates, ", "), argCount)

// 	var updated Class
// 	err := myDb.Db.QueryRow(query, args...).Scan(&updated.ID, &updated.Name, &updated.Teacher, &updated.Day, &updated.Time)

// 	if err == sql.ErrNoRows {
// 		return nil, fmt.Errorf("active student with id %d not found", id)
// 	}
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &updated, nil
// }

// // TODO: standard is to just call this delete not softDelete. Add comment about soft delete
// func SoftDeleteClassDB(myDb *db.MyDatabase, id int) (string, error) {

// 	var name string

// 	err := myDb.Db.QueryRow(
// 		"SELECT name FROM classes WHERE id=$1 AND active=true", id,
// 	).Scan(&name)

// 	if err == sql.ErrNoRows {
// 		return "", fmt.Errorf("class with id %d not found", id)
// 	}

// 	result, err := myDb.Db.Exec(
// 		"UPDATE classes SET active = false WHERE id = $1",
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
// 		return "", fmt.Errorf("class with id %d not found", id)
// 	}

// 	return name, nil
// }
