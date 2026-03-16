package roster

import (
	"IFTP/db"
	"IFTP/roster"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

func dbGetRoster(ctx context.Context, myDb *db.MyDatabase, classId int, month time.Time, class_date time.Time) (GetRosterRequest, error) {

	var roster GetRosterRequest

	err := myDb.Pool.QueryRow(ctx, `
		SELECT 
			c.name AS class_name,
			-- Session dates in the specific month
			(
        SELECT COALESCE(array_agg(DISTINCT session_date ORDER BY session_date), '{}')
			FROM class_schedule 
			WHERE class_id = $1 AND month = $2
		) AS session_dates,
    -- Total enrollment count for that specific date
    (
        SELECT COUNT(DISTINCT student_id) 
			FROM ROSTER 
			WHERE class_id = $1 AND class_date = $3
    	) AS enrolled_count,
    -- Aggregate student list into a JSON array
    COALESCE(
        jsonb_agg(
            jsonb_build_object(
                'id', s.id,
                'name', s.name,
                'email', s.email,
                'status', r.status
            ) ORDER BY s.name
        ) FILTER (WHERE s.id IS NOT NULL), 
        '[]'
    ) AS students
	FROM classes c
	LEFT JOIN ROSTER r ON r.class_id = c.id AND r.class_date = $3
	LEFT JOIN Students s ON s.id = r.student_id
	WHERE c.id = $1
	GROUP BY c.name;`, classId, month, class_date).Scan(&roster.ClassName, &roster.SessionDates, &roster.EnrolledCount, &roster.Students)

	return roster, err
}

func dbGetStudentEnrollment(ctx context.Context, myDb *db.MyDatabase, studentId int, month time.Time) ([]StudentEnrollment, error) {

	rows, err := myDb.Pool.Query(ctx, `
	SELECT c.name as class_name, r.class_date, cs.month
	FROM roster r
	JOIN classes c on r.class_id = c.id
	JOIN class_schedule cs on cs.class_id = r.class_id
	AND cs.session_date = r.class_date
	WHERE r.student_id = $1
	AND cs.month = $2;
	`, studentId, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requestedStudentEnrollment, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[StudentEnrollment])
	if err != nil {
		return nil, err
	}

	return requestedStudentEnrollment, nil
}

func dbGetStudentInfo(ctx context.Context, myDb *db.MyDatabase, request *roster.EnrollmentRequestApproval, studentId int) error {
	err := myDb.Pool.QueryRow(ctx,
		`SELECT s.name, s.email, COALESCE(ARRAY_AGG(DISTINCT c.name ORDER BY c.name) FILTER (WHERE c.name IS NOT NULL), '{}') AS currently_enrolled FROM students AS s
	LEFT JOIN roster AS r on r.student_id = s.id
	LEFT JOIN classes AS c on r.class_id = c.id
	WHERE s.id = $1 AND s.active = true
	GROUP BY s.name, s.email`, studentId).Scan(&request.StudentName, &request.StudentEmail, &request.CurrentlyEnrolled)

	return err
}

func dbGetClassInfo(ctx context.Context, myDb *db.MyDatabase, request *roster.EnrollmentRequestApproval) error {
	err := myDb.Pool.QueryRow(ctx,
		`SELECT c.id, c.name, c.teacher,
		c.capacity - COUNT(DISTINCT r.student_id) as available_spots
		FROM classes AS c
		LEFT JOIN class_schedule AS cs ON cs.class_id = c.id
		LEFT JOIN roster AS r ON r.class_id = c.id AND r.class_date = cs.session_date
		WHERE c.id = $1 AND cs.month = $2 AND c.active = True
		GROUP BY cs.month, c.id
		ORDER  BY cs.month DESC`, request.ClassID, request.Month).Scan(&request.RequestedClassID, &request.RequestedClassName, &request.Teacher, &request.AvailableSpots)

	return err
}

func dbEnrollmentReqExists(ctx context.Context, tx pgx.Tx, request *roster.EnrollmentRequestApproval, studentId int) (bool, error) {
	var exists bool

	err := tx.QueryRow(ctx, `
	SELECT EXISTS FROM enrollment_requests
	WHERE student_id = $1
	AND requested_class_id = $2
	AND status = 'Pending'
	`, studentId, request.RequestedClassID).Scan(&exists)

	return exists, err
}

func dbStudentAlreadyEnrolled(ctx context.Context, tx pgx.Tx, request *roster.EnrollmentRequestApproval, studentId int) (bool, error) {
	var enrolled bool

	err := tx.QueryRow(ctx, `
	SELECT EXISTS FROM roster
	JOIN class_schedule cs on cs.class_id = r.class_id
	AND cs.session_date = r.class_date
	WHERE student_id = $1
	AND class_id = $2
	AND status = 'Enrolled'
	AND cs.month = $3
	`, studentId, request.RequestedClassID, request.Month).Scan(&enrolled)

	return enrolled, err

}

// func dbEnroll(ctx context.Context, myDb *db.MyDatabase, classID int, classDate time.Time, studentID int) error {

// 	var rosterID int

// 	err := myDb.Pool.QueryRow(
// 		"INSERT INTO roster (class_date, student_id, class_id, registration_date) "+
// 			"VALUES ($1, $2, $3, NOW()) RETURNING id", classDate, studentID, classID).Scan(&rosterID)

// 	return err
// }

// Should have a function to view the roster for a specific class date- as input needs the class id and the date. We want to see the student names and emails.

//Should have a function to view the roster for a specific class for the month- needs the class id.
