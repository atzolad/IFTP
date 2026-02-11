package class

import (
	"IFTP/db"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

func dbListClasses(myDb *db.MyDatabase) ([]Class, error) {
	rows, err := myDb.Db.Query(
		`SELECT c.id, name, teacher, day_of_week, time, description, capacity, cs.month, ARRAY_AGG(DISTINCT cs.session_date ORDER BY cs.session_date) AS session_dates, COUNT(DISTINCT r.student_id) AS enrolled_count
		FROM classes AS c
		JOIN class_schedule AS cs ON cs.class_id = c.id
		LEFT JOIN roster AS r ON r.class_id = c.id
		WHERE active = True
		GROUP BY cs.month, c.id
		ORDER  BY cs.month DESC`)

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
			&class.Description, &class.Capacity, &class.Month, (*pq.StringArray)(&class.SessionDates), &class.EnrolledCount); err != nil {
			return nil, err
		}
		classes = append(classes, class)
	}
	if err = rows.Err(); err != nil {
		return classes, err
	}
	return classes, nil
}

func dbListClassesByMonth(myDb *db.MyDatabase, month string, studentId *int) ([]Class, error) {

	var query strings.Builder
	var args []any

	query.WriteString(`SELECT c.id, name, teacher, day_of_week, time, description, capacity, cs.month, ARRAY_AGG(DISTINCT cs.session_date ORDER BY cs.session_date) AS session_dates, COUNT(DISTINCT r.student_id) AS enrolled_count
		FROM classes AS c
		JOIN class_schedule AS cs ON cs.class_id = c.id`)

	if month != "" {
		fmt.Printf("Month: %v", month)
		args = append(args, month)
		fmt.Fprintf(&query, " AND month = $%d ", len(args))
	}

	query.WriteString(`
			LEFT JOIN roster AS r ON r.class_id = c.id
			WHERE c.active = True`)

	if studentId != nil {
		fmt.Printf("student id: %v ", *studentId)
		args = append(args, *studentId)
		fmt.Fprintf(&query, " AND r.student_id = $%d ", len(args))
	}

	query.WriteString(" GROUP BY cs.month, c.id, c.name, c.teacher, c.day_of_week, c.time, c.description, c.capacity")
	query.WriteString(" ORDER BY cs.month DESC")
	// var rows *sql.Rows
	// var err error
	fmt.Println(query.String())
	rows, err := myDb.Db.Query(query.String(), args...)

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
			&class.Description, &class.Capacity, &class.Month, (*pq.StringArray)(&class.SessionDates), &class.EnrolledCount); err != nil {
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

func dbListStudentEnrolledClasses(myDb *db.MyDatabase, month string, studentId *int) ([]Class, error) {
	var query strings.Builder
	var args []any

	query.WriteString(`SELECT
    c.id,
    c.name,
    c.teacher,
    c.day_of_week,
    c.time,
    c.description,
    c.capacity,
    cs.month,
    ARRAY_AGG(DISTINCT cs.session_date ORDER BY cs.session_date) AS session_dates,
    COUNT(DISTINCT r_all.student_id) AS enrolled_count
FROM classes c
JOIN class_schedule cs ON cs.class_id = c.id
JOIN roster r_student ON r_student.class_id = c.id AND r_student.class_date = cs.session_date
LEFT JOIN roster r_all ON r_all.class_id = c.id AND r_all.class_date = cs.session_date
WHERE c.active = true `)

	if studentId != nil {
		args = append(args, *studentId)
		fmt.Fprintf(&query, " AND r_student.student_id = $%d ", len(args))
	}

	if month != "" {
		args = append(args, month)
		fmt.Fprintf(&query, " AND cs.month = $%d ", len(args))
	}

	query.WriteString(`
GROUP BY
    cs.month,
    c.id,
    c.name,
    c.teacher,
    c.day_of_week,
    c.time,
    c.description,
    c.capacity;`)

	fmt.Println("QUERY:", query.String())

	rows, err := myDb.Db.Query(query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []Class

	for rows.Next() {

		var class Class
		if err := rows.Scan(&class.ID, &class.Name, &class.Teacher, &class.DayOfWeek, &class.Time,
			&class.Description, &class.Capacity, &class.Month, (*pq.StringArray)(&class.SessionDates), &class.EnrolledCount); err != nil {
			return nil, err
		}
		classes = append(classes, class)
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
