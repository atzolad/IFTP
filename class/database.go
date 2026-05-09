package class

import (
	"IFTP/db"
	"IFTP/utils"
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
		LEFT JOIN roster AS r ON r.class_id = c.id AND r.class_date = cs.session_date AND r.status = 'Enrolled'
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

func dbListClassesByMonth(ctx context.Context, myDb *db.MyDatabase, month string, studentId string) ([]Class, error) {

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
		JOIN class_schedule AS cs ON cs.class_id = c.id`)

	if month != "" {
		args = append(args, month)
		fmt.Fprintf(&query, " AND month = $%d ", len(args))
	}

	query.WriteString(` 
			LEFT JOIN roster AS r ON r.class_id = c.id AND r.class_date = cs.session_date AND r.status = 'Enrolled'
			WHERE c.active = True`)

	if studentId != "" {
		fmt.Printf("student id: %v ", studentId)
		args = append(args, studentId)
		fmt.Fprintf(&query, " AND r.student_id = $%d ", len(args))
	}

	query.WriteString(" GROUP BY cs.month, c.id, c.name, c.teacher, c.day_of_week, c.time, c.description, c.capacity")
	query.WriteString(" ORDER BY c.name, cs.month DESC")

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

func dbListStudentEnrolledClasses(ctx context.Context, myDb *db.MyDatabase, month string, studentId string) ([]Class, error) {
	var query strings.Builder
	var args []any

	var studentJoin string
	if studentId != "" {
		args = append(args, studentId)
		studentJoin = fmt.Sprintf("LEFT JOIN roster r_student ON r_student.class_id = c.id AND r_student.class_date = cs.session_date AND r_student.student_id = $%d", len(args))
	} else {
		studentJoin = "LEFT JOIN roster r_student ON r_student.class_id = c.id AND r_student.class_date = cs.session_date"
	}

	query.WriteString(fmt.Sprintf(`
	SELECT c.id, c.name, c.teacher, c.day_of_week, c.time, c.description, c.capacity, cs.month, ARRAY_AGG(DISTINCT cs.session_date ORDER BY cs.session_date) AS session_dates, COUNT(DISTINCT r_all.student_id) AS enrolled_count
	FROM classes c
	JOIN class_schedule cs ON cs.class_id = c.id
	%s
	LEFT JOIN roster r_all ON r_all.class_id = c.id AND r_all.class_date = cs.session_date AND r_all.status = 'Enrolled'
	WHERE c.active = true 
	`, studentJoin))

	if studentId != "" {
		query.WriteString(" AND r_student.student_id IS NOT NULL ")
	}

	if month != "" {
		args = append(args, month)
		fmt.Fprintf(&query, " AND cs.month = $%d ", len(args))
	}

	query.WriteString(`GROUP BY cs.month, c.id, c.name, c.teacher, c.day_of_week, c.time, c.description, c.capacity;`)

	rows, err := myDb.Pool.Query(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	classes, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[Class])
	if err != nil {
		return nil, fmt.Errorf("Error retrieving student enrolled classes from db: %v", err)
	}

	return classes, nil
}

func dbCreateClass(ctx context.Context, tx pgx.Tx, c *Class) error {
	err := tx.QueryRow(ctx,
		"INSERT INTO classes (name, teacher, day_of_week, time, description, capacity) VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
		c.Name, c.Teacher, c.DayOfWeek, c.Time, c.Description, c.Capacity).Scan(&c.ID)

	return err
}

func dbInsertClassScheduleRows(ctx context.Context, tx pgx.Tx, classId string, month time.Time, dates []time.Time) error {
	for _, date := range dates {
		_, err := tx.Exec(ctx,
			`INSERT INTO class_schedule (class_id, session_date, month, status) 
		VALUES ($1, $2, $3, $4)`,
			classId, date, month, "Scheduled")
		if err != nil {
			return err
		}
	}
	return nil
}

func dbUpdateClass(ctx context.Context, tx pgx.Tx, id string, c *Class) (*Class, error) {

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
		return nil, utils.ErrNoFieldsToUpdate
	}

	args = append(args, id)

	query := fmt.Sprintf("UPDATE classes SET %s WHERE id=$%d RETURNING id, name, teacher, day_of_week, time, description, capacity",
		strings.Join(updates, ", "), len(args))

	var updated Class
	err := tx.QueryRow(ctx, query, args...).Scan(&updated.ID, &updated.Name, &updated.Teacher, &updated.DayOfWeek, &updated.Time, &updated.Description, &updated.Capacity)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("class with id %v not found", id)
		}
		return nil, err
	}

	return &updated, nil
}

func dbDeleteFromClassSchedule(ctx context.Context, tx pgx.Tx, id string) error {
	_, err := tx.Exec(ctx, "DELETE FROM class_schedule WHERE class_id = $1", id)
	return err
}

func dbInsertScheduleApproval(ctx context.Context, tx pgx.Tx, classId string, month time.Time) (approvalId string, err error) {
	err = tx.QueryRow(ctx, `
	INSERT INTO schedule_approvals (class_id, month)
	VALUES ($1, $2)
	RETURNING id`, classId, month).Scan(&approvalId)

	if err != nil {
		return "", err
	}

	return approvalId, err

}

func dbInsertPendingDates(ctx context.Context, tx pgx.Tx, approvalId string, pendingDates []time.Time) (err error) {
	for _, date := range pendingDates {
		_, err = tx.Exec(ctx, `
			INSERT INTO schedule_approval_dates (schedule_approval_id, proposed_date)
			VALUES ($1, $2)`, approvalId, date)
		if err != nil {
			return err
		}
	}
	return nil
}

func dbGetPendingScheduleApprovals(ctx context.Context, myDb *db.MyDatabase) (approvals []MonthlyClassScheduleApproval, err error) {
	rows, err := myDb.Pool.Query(ctx, `
		SELECT sa.id, sa.class_id, c.name as class_name, c.teacher, c.day_of_week, c.time, c.capacity, sa.month, sa.status, ARRAY_AGG(sad.proposed_date ORDER BY sad.proposed_date) AS pending_dates
		FROM schedule_approvals sa
		JOIN classes c ON c.id = sa.class_id
		JOIN schedule_approval_dates sad ON sad.schedule_approval_id = sa.id
		WHERE sa.status = 'Pending'
		GROUP BY sa.id, sa.class_id, c.name, c.teacher, c.day_of_week, c.time, c.capacity, sa.month, sa.status
		ORDER BY sa.month ASC, c.name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	approvals, err = pgx.CollectRows(rows, pgx.RowToStructByNameLax[MonthlyClassScheduleApproval])
	if err != nil {
		return nil, err
	}

	return approvals, nil
}

func dbUpdateScheduleApprovalStatus(ctx context.Context, tx pgx.Tx, approvalId string, status string) (classId string, month time.Time, err error) {
	err = tx.QueryRow(ctx, `
	UPDATE schedule_approvals
	SET status = $1, reviewed_at = NOW()
	WHERE id = $2
	RETURNING class_id, month`, status, approvalId).Scan(&classId, &month)

	return
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
