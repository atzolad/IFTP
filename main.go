package main

import (
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// func main() {

// 	// Load env vars from .env file
// 	if err := godotenv.Load(); err != nil {
// 		log.Println("Warning: no .env file found")
// 	}

// 	router := gin.Default()
// 	connStr := os.Getenv("CONN_STR")
// 	fmt.Println("Connecting with:", connStr)
// 	// Initialise the connection pool.
// 	sqldb, err := sql.Open("pgx", connStr)

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("Succesfully connected to database \n")
// 	defer sqldb.Close()

// 	// Test the connection
// 	if err := sqldb.Ping(); err != nil {
// 		log.Fatal(err)
// 	}

// 	// Create an instance of myDatabase containing the connection pool.
// 	myDb := &db.MyDatabase{Db: sqldb}

// 	// Studentendpoints- make sure to pass the database instance to each function.
// 	router.GET("/students", student.GetStudents(myDb))
// 	router.POST("/students", student.AddStudent(myDb))
// 	router.PATCH("/students/:id", student.UpdateStudent(myDb))
// 	router.DELETE("/students/:id", student.SoftDeleteStudent(myDb))

// 	// Class endpoints
// 	router.GET("/classes", class.GetClasses(myDb))
// 	router.POST("/classes", class.AddClass(myDb))
// 	router.PATCH("/classes/:id", class.UpdateClass(myDb))
// 	router.DELETE("/classes/:id", class.SoftDeleteClass(myDb))

// 	// Roster endpoints
// 	// router.GET("/roster", roster.GetRoster(myDb))
// 	// router.GET("/roster/:student_id/classes,", roster.GetStudentClasses)
// 	// router.GET("/roster/:class_id/students", roster.GetClassStudents)
// 	router.POST("/roster/:class_id/enroll", roster.Enroll(myDb))
// 	// router.DELETE("/roster/:id", roster.LeaveClass(myDb))
// 	// router.PATCH("/roster/:id", roster.UpdateRoster(myDb))

// 	router.Run()
// }

// FirstMonday returns the day of the first Monday in the given month.
func FirstDay(weekday time.Weekday, year int, month time.Month) int {
	t := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	return (8-int(t.Weekday()))%7 + int(weekday)
}

// datesMap map[time.Weekday][]ints

// 'Friday': 3
// 'saturday':4

func main() {
	// date, _ := time.Parse("2006-1-2", "2025-12-1")
	// roster.FindMonthlyDates(time.Date(year, month, 1))

	// if datesMap['friday'].exists()
	// else {
	// 	friday = FirstDay(time.Friday, 2025, 11)
	// 	datesMap['friday'] = friday
	// }

	fmt.Println(FirstDay(time.Friday, 2025, 12))

}
