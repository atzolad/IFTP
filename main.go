package main

import (
	"IFTP/class"
	"IFTP/db"
	"IFTP/roster"
	"IFTP/students"
	"IFTP/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
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

	connStr := os.Getenv("CONN_STR")
	fmt.Println("Connecting with:", connStr)

	// Initialise the connection pool.
	ctx := context.Background()
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v \n", err)
	}

	defer pool.Close()

	// Test the connection
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		log.Fatalf("Error: Could not ping database: %v \n", err)
	}

	fmt.Printf("Succesfully connected to database \n")

	// Load all the html templates from the templates directory.
	tpl := utils.LoadTemplates()

	// Create an instance of myDatabase containing the connection pool.
	myDb := &db.MyDatabase{
		Pool:      pool,
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
	mux.HandleFunc("GET /classes/all", class.ListClasses(myDb))
	mux.HandleFunc("GET /classes", class.ListClassesByMonth(myDb))
	mux.HandleFunc("GET /classes/{student_id}", class.ListClassesByMonth(myDb))
	mux.HandleFunc("PATCH /classes/{class_id}", class.UpdateClass(myDb))
	mux.HandleFunc("POST /classes", class.CreateClass(myDb))
	mux.HandleFunc("POST /classes/schedule_approval/generate", class.TriggerScheduleApprovals(myDb))
	mux.HandleFunc("GET /classes/schedule_approval", class.GetPendingScheduleApprovals(myDb))
	mux.HandleFunc("PATCH /classes/schedule_approval/confirm", class.ConfirmScheduleApproval(myDb))

	// Roster Endpoints
	mux.HandleFunc("GET /roster/{class_id}", roster.GetRoster(myDb))
	mux.HandleFunc("GET /roster/enrollment/{student_id}", roster.GetStudentEnrollment(myDb))

	//Enrollment Request Endpoints
	mux.HandleFunc("POST /enrollment_requests", roster.CreateEnrollmentRequest(myDb))
	mux.HandleFunc("GET /enrollment_requests", roster.GetEnrollmentRequests(myDb))
	mux.HandleFunc("PATCH /enrollment_requests/{request_id}", roster.UpdateEnrollmentRequest(myDb))

	// Calendar Endpoints
	mux.HandleFunc("GET /calendarEvents", class.GetCalendarEvents(myDb))
	mux.HandleFunc("GET /calendarEvents/{student_id}", class.GetCalendarEventsByStudent(myDb))

	// Student Endpoints
	mux.HandleFunc("GET /students", students.GetStudents(myDb))
	mux.HandleFunc("GET /students/enrollment", students.GetStudentsWithEnrollment(myDb))
	mux.HandleFunc("POST /students", students.AddStudent(myDb))
	mux.HandleFunc("PATCH /students/{student_id}", students.UpdateStudent(myDb))

	myDb.Logger.Printf("Server starting on :%v", port)
	if err := http.ListenAndServe(":"+port, wrappedMux); err != nil {
		log.Fatalf("Error initializing server: %v", err)
	}

}
