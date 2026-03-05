package roster

import (
	"IFTP/db"
	"context"
	"time"
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

// func dbEnroll(ctx context.Context, myDb *db.MyDatabase, classID int, classDate time.Time, studentID int) error {

// 	var rosterID int

// 	err := myDb.Pool.QueryRow(
// 		"INSERT INTO roster (class_date, student_id, class_id, registration_date) "+
// 			"VALUES ($1, $2, $3, NOW()) RETURNING id", classDate, studentID, classID).Scan(&rosterID)

// 	return err
// }

// Should have a function to view the roster for a specific class date- as input needs the class id and the date. We want to see the student names and emails.

//Should have a function to view the roster for a specific class for the month- needs the class id.
