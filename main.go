package main

import (
	"IFTP/class"
	"IFTP/db"
	"IFTP/students"
	"IFTP/utils"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {

	// Load env vars from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: no .env file found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Railway typically uses this as default
	}

	// router := gin.Default()
	// // Load all the html templates from the templates directory.
	// router.LoadHTMLGlob("templates/*")

	connStr := os.Getenv("CONN_STR")
	fmt.Println("Connecting with:", connStr)

	// Initialise the connection pool.
	sqldb, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Succesfully connected to database \n")
	defer sqldb.Close()

	// Test the connection
	if err := sqldb.Ping(); err != nil {
		log.Fatal(err)
	}

	// Load all the html templates from the templates directory.
	tpl := utils.LoadTemplates()

	// Create an instance of myDatabase containing the connection pool.
	myDb := &db.MyDatabase{
		Db:        sqldb,
		Logger:    log.Default(),
		Templates: tpl.Index,
	}

	mux := http.NewServeMux()

	baseUrl := "http://localhost:8080"

	// Wrap the serve mux with a logger for http requests and responses.
	wrappedMux := utils.LoggingWrapper(mux)

	// Render the main index.
	mux.HandleFunc("/", utils.IndexHandler(tpl, baseUrl))

	// Class Endpoints
	mux.HandleFunc("GET /classes", class.ListClassesByMonth(myDb))
	mux.HandleFunc("GET /classes/{student_id}", class.ListClassesByMonth(myDb))
	mux.HandleFunc("GET /calendarEvents", class.GetCalendarEvents(myDb))
	mux.HandleFunc("GET /calendarEvents/{student_id}", class.GetCalendarEventsByStudent(myDb))
	// mux.HandleFunc("POST /classes", class.CreateClass(myDb))

	// Student Endpoints
	mux.HandleFunc("GET /students", students.GetStudents(myDb))
	mux.HandleFunc("GET /students/enrollment", students.GetStudentsWithEnrollment(myDb))
	// mux.HandleFunc("Post" /students, students.AddStudent(myDb))

	log.Printf("Server starting on :%v", port)
	if err := http.ListenAndServe(":"+port, wrappedMux); err != nil {
		log.Fatalf("Error initializing server: %v", err)
	}

	// // Render the main index.
	// router.GET("/", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.html", gin.H{
	// 		"title": "IFTP Main Page",
	// 	})
	// })
	// // Studentendpoints- make sure to pass the database instance to each function.
	// router.GET("/students", student.GetStudents(myDb))
	// router.POST("/students", student.AddStudent(myDb))
	// router.PATCH("/students/:id", student.UpdateStudent(myDb))
	// router.DELETE("/students/:id", student.SoftDeleteStudent(myDb))

	// // Class endpoints
	// router.GET("/classes", class.ListClasses(myDb))
	// router.GET("/classes/:student_id", class.ListClassesByMonth(myDb))
	// // router.GET("/classes/:month", class.ListClassesByMonth(myDb))
	// // router.POST("/classes", class.CreateClass(myDb))
	// // router.PATCH("/classes/:id", class.UpdateClass(myDb))
	// // router.DELETE("/classes/:id", class.SoftDeleteClass(myDb))

	// // Roster endpoints
	// // router.GET("/roster", roster.GetRoster(myDb))
	// // router.GET("/roster/:student_id/classes,", roster.GetStudentClasses)
	// // router.GET("/roster/:class_id/students", roster.GetClassStudents)
	// router.POST("/roster/:class_id/enroll", roster.Enroll(myDb))
	// // router.DELETE("/roster/:id", roster.LeaveClass(myDb))
	// // router.PATCH("/roster/:id", roster.UpdateRoster(myDb))

	// router.Run()
}

// FirstDay returns the day of the first Monday in the given month.
// func FirstDay(weekday time.Weekday, year int, month time.Month) int {
// 	t := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
// 	return (8-int(t.Weekday()))%7 + int(weekday)
// }

// datesMap map[time.Weekday][]ints

// 'Friday': 3
// 'saturday':4

// func FindMonthlyDates(starting_date time.Time) error {
// 	// Set the starting date
// 	//Check all the days after that- increment by 7 if the month is the same, add it , if not do not

// 	// // set the starting date (in any way you wish) - replace with starting_date- if input is a string
// 	// start, err := time.Parse("2006-1-2", starting_date)
// 	// if err != nil {
// 	// 	return fmt.Errorf("error parsing starting date for monthly day / date aggregation: %v", err)
// 	// }
// 	// handle error

// 	// set d to starting date and keep adding 7 days to it as long as month doesn't change
// 	for d := starting_date; d.Month() == starting_date.Month(); d = d.AddDate(0, 0, 7) {
// 		date := d.String()
// 		fmt.Println(date)
// 	}
// 	return nil
// }

// func main() {
// 	month := time.November
// 	monday := time.Monday
// 	tuesday := time.Tuesday
// 	daysSlice := []time.Weekday{monday, tuesday}
// 	datesMap := timeutils.CreateDatesMap(daysSlice, 2025, month)
// 	for day, dates := range datesMap {
// 		fmt.Printf("Day: %v /n", day)
// 		for _, date := range dates {
// 			fmt.Printf("Date: %v", date)
// 		}
// 		fmt.Println()
// 	}

// if datesMap['friday'].exists()
// else {
// 	friday = FirstDay(time.Friday, 2025, 11)
// 	datesMap['friday'] = friday
// }

// fmt.Println(FirstDay(time.Friday, 2025, 12))

// }
