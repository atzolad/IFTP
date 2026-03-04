package class

import (
	"IFTP/db"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

func dbListClasses(ctx context.Context, myDb *db.MyDatabase) ([]Class, error) {
	rows, err := myDb.Pool.Query(ctx,
		`SELECT c.id, name, teacher, day_of_week, time, description, capacity, COALESCE(cs.month, '0001-01-01'::date) AS month, 
		COALESCE(
			ARRAY_AGG(DISTINCT cs.session_date ORDER BY cs.session_date) 
			FILTER (WHERE cs.session_date IS NOT NULL), 
			'{}'
		) AS session_dates, 
		COUNT(DISTINCT r.student_id) AS enrolled_count
		FROM classes AS c
		LEFT JOIN class_schedule AS cs ON cs.class_id = c.id
		LEFT JOIN roster AS r ON r.class_id = c.id
		WHERE active = True
		GROUP BY cs.month, c.id
		ORDER  BY cs.month DESC`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	classes, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[Class])
	if err != nil {
		return nil, err
	}
	return classes, nil
}

func dbListClassesByMonth(ctx context.Context, myDb *db.MyDatabase, month string, studentId *int) ([]Class, error) {

	var query strings.Builder
	var args []any

	query.WriteString(`SELECT c.id, name, teacher, day_of_week, time, description, capacity, COALESCE(cs.month, '0001-01-01'::date) AS month, 
		COALESCE(
			ARRAY_AGG(DISTINCT cs.session_date ORDER BY cs.session_date) 
			FILTER (WHERE cs.session_date IS NOT NULL), 
			'{}'
		) AS session_dates, 
		COUNT(DISTINCT r.student_id) AS enrolled_count
		FROM classes AS c
		LEFT JOIN class_schedule AS cs ON cs.class_id = c.id`)

	if month != "" {
		fmt.Printf("Month: %v", month)
		args = append(args, month)
		fmt.Fprintf(&query, " AND month = $%d ", len(args))
	}

	query.WriteString(`
			LEFT JOIN roster AS r ON r.class_id = c.id AND r.class_date = cs.session_date
			WHERE c.active = True`)

	if studentId != nil {
		fmt.Printf("student id: %v ", *studentId)
		args = append(args, *studentId)
		fmt.Fprintf(&query, " AND r.student_id = $%d ", len(args))
	}

	query.WriteString(" GROUP BY cs.month, c.id, c.name, c.teacher, c.day_of_week, c.time, c.description, c.capacity")
	query.WriteString(" ORDER BY cs.month DESC")

	fmt.Println(query.String())
	rows, err := myDb.Pool.Query(ctx, query.String(), args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	classes, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[Class])
	if err != nil {
		return nil, fmt.Errorf("Error retrieving classes from db: %v", err)
	}

	return classes, nil

	// // Loop through rows, using Scan to assign column data to struct fields.
	// for rows.Next() {
	// 	var class Class
	// 	if err := rows.Scan(&class.ID, &class.Name, &class.Teacher, &class.DayOfWeek, &class.Time,
	// 		&class.Description, &class.Capacity, &class.Month, (*pq.StringArray)(&class.SessionDates), &class.EnrolledCount); err != nil {
	// 		return nil, err
	// 	}
	// 	fmt.Println(classes)
	// 	classes = append(classes, class)
	// }
	// if err = rows.Err(); err != nil {
	// 	return classes, err
	// }

}

func dbListStudentEnrolledClasses(ctx context.Context, myDb *db.MyDatabase, month string, studentId *int) ([]Class, error) {
	var query strings.Builder
	var args []any

	query.WriteString(`
	SELECT c.id, c.name, c.teacher, c.day_of_week, c.time, c.description, c.capacity, cs.month, ARRAY_AGG(DISTINCT cs.session_date ORDER BY cs.session_date) AS session_dates, COUNT(DISTINCT r_all.student_id) AS enrolled_count
	FROM classes c
	JOIN class_schedule cs ON cs.class_id = c.id
	JOIN roster r_student ON r_student.class_id = c.id AND r_student.class_date = cs.session_date
	LEFT JOIN roster r_all ON r_all.class_id = c.id AND r_all.class_date = cs.session_date
	WHERE c.active = true 
	`)

	if studentId != nil {
		args = append(args, *studentId)
		fmt.Fprintf(&query, " AND r_student.student_id = $%d ", len(args))
	}

	if month != "" {
		args = append(args, month)
		fmt.Fprintf(&query, " AND cs.month = $%d ", len(args))
	}

	query.WriteString(`GROUP BY cs.month, c.id, c.name, c.teacher, c.day_of_week, c.time, c.description, c.capacity;`)

	fmt.Println("QUERY:", query.String())

	rows, err := myDb.Pool.Query(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// var classes []Class

	classes, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[Class])
	if err != nil {
		return nil, fmt.Errorf("Error retrieving student enrolled classes from db: %v", err)
	}

	// for rows.Next() {

	// 	var class Class
	// 	if err := rows.Scan(&class.ID, &class.Name, &class.Teacher, &class.DayOfWeek, &class.Time,
	// 		&class.Description, &class.Capacity, &class.Month, (*pq.StringArray)(&class.SessionDates), &class.EnrolledCount); err != nil {
	// 		return nil, err
	// 	}
	// 	classes = append(classes, class)
	// }
	return classes, nil
}

func dbCreateClass(ctx context.Context, tx pgx.Tx, c *Class) error {
	err := tx.QueryRow(ctx,
		"INSERT INTO classes (name, teacher, day_of_week, time, description, capacity) VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
		c.Name, c.Teacher, c.DayOfWeek, c.Time, c.Description, c.Capacity).Scan(&c.ID)

	return err
}

func dbInsertClass_ScheduleRow(ctx context.Context, tx pgx.Tx, c *Class, sessionDate time.Time) error {
	_, err := tx.Exec(ctx,
		`INSERT INTO class_schedule (class_id, session_date, month, status) 
		VALUES ($1, $2, $3, $4)`,
		c.ID, sessionDate, c.Month, "Scheduled")

	return err
}

func dbUpdateClass(ctx context.Context, tx pgx.Tx, id int, c *Class) (*Class, error) {

	// month is intentionally excluded — it belongs to class_schedule, not classes

	updates := []string{}
	args := []any{}

	if c.Name != "" {
		args = append(args, c.Name)
		updates = append(updates, fmt.Sprintf("name=$%d", len(args)))

	}

	if c.Teacher != "" {
		args = append(args, c.Teacher)
		updates = append(updates, fmt.Sprintf("teacher=$%d", len(args)))

	}

	if c.DayOfWeek != "" {
		args = append(args, c.DayOfWeek)
		updates = append(updates, fmt.Sprintf("day_of_week=$%d", len(args)))

	}

	if c.Time != "" {
		args = append(args, c.Time)
		updates = append(updates, fmt.Sprintf("time=$%d", len(args)))

	}

	if c.Description != "" {
		args = append(args, c.Description)
		updates = append(updates, fmt.Sprintf("description=$%d", len(args)))

	}

	if c.Capacity != 0 {
		args = append(args, c.Capacity)
		updates = append(updates, fmt.Sprintf("capacity=$%d", len(args)))
	}

	if !c.EndDate.IsZero() {
		args = append(args, c.EndDate)
		updates = append(updates, fmt.Sprintf("end_date=$%d", len(args)))
	}

	if len(updates) == 0 {
		return nil, ErrNoFieldsToUpdate
	}

	args = append(args, id)

	query := fmt.Sprintf("UPDATE classes SET %s WHERE id=$%d RETURNING id, name, teacher, day_of_week, time, description, capacity",
		strings.Join(updates, ", "), len(args))

	var updated Class
	err := tx.QueryRow(ctx, query, args...).Scan(&updated.ID, &updated.Name, &updated.Teacher, &updated.DayOfWeek, &updated.Time, &updated.Description, &updated.Capacity)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("class with id %d not found", id)
		}
		return nil, err
	}

	return &updated, nil
}

func dbDeleteFromClassSchedule(ctx context.Context, tx pgx.Tx, id int) error {
	_, err := tx.Exec(ctx, "DELETE FROM class_schedule WHERE class_id = $1", id)
	return err
}

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
