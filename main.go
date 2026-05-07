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

	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	googleProvider "github.com/markbates/goth/providers/google"
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

	clientKey := os.Getenv("GOOGLE_KEY")
	clientSecret := os.Getenv("GOOGLE_SECRET")
	scopes := os.Getenv("GOOGLE_SCOPES")
	baseUrl := "http://localhost:8080"
	callBackUrl := fmt.Sprintf("%v/auth/google/callback", baseUrl)

	gp := googleProvider.New(clientKey, clientSecret, callBackUrl, scopes)
	gp.SetPrompt("select_account")
	goth.UseProviders(gp)

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

	// Configure Gothic store parameters
	key := os.Getenv("SESSION_SECRET")
	if key == "" {
		log.Fatal("Session Secret env variable missing")
	}
	maxAge := 86400 * 30 // 30 days
	isProd := false      // Set to true when serving over https
	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd
	store.Options.SameSite = http.SameSiteLaxMode
	gothic.Store = store

	// Handle file serving
	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	protectedMux := http.NewServeMux()

	// Authorization endpoints
	mux.HandleFunc("GET /auth/google", utils.GoogleLogin())
	mux.HandleFunc("GET /auth/google/callback", utils.HandleGoogleOauth(myDb))
	protectedMux.HandleFunc("GET /auth/user", utils.GetUserAuth(myDb))
	mux.HandleFunc("GET /unauthorized", utils.UnAuthorizedHandler(tpl))
	mux.HandleFunc("GET /logout", utils.HandleLogout())
	mux.HandleFunc("GET /health", utils.HealthCheckHandler(myDb))

	// Require authorization for protected endpoints and wrap all endpoints with logger.
	mux.Handle("/", utils.RequireAuth(protectedMux))
	wrappedMux := utils.LoggingWrapper(mux)

	// Wrap the serve mux with a logger for http requests and responses.
	// wrappedMux := utils.LoggingWrapper(mux)

	// Render the main index.
	protectedMux.HandleFunc("/", utils.IndexHandler(tpl, baseUrl))

	// Class Endpoints
	protectedMux.HandleFunc("GET /classes/all", class.ListClasses(myDb))
	protectedMux.HandleFunc("GET /classes", class.ListClassesByMonth(myDb))
	protectedMux.HandleFunc("GET /classes/{student_id}", class.ListClassesByMonth(myDb))
	protectedMux.HandleFunc("PATCH /classes/{class_id}", utils.RequireAdmin(class.UpdateClass(myDb)))
	protectedMux.HandleFunc("POST /classes", utils.RequireAdmin(class.CreateClass(myDb)))
	protectedMux.HandleFunc("POST /classes/schedule_approval/generate", utils.RequireAdmin(class.TriggerScheduleApprovals(myDb)))
	protectedMux.HandleFunc("GET /classes/schedule_approval", utils.RequireAdmin(class.GetPendingScheduleApprovals(myDb)))
	protectedMux.HandleFunc("PATCH /classes/schedule_approval/confirm/{approval_id}", class.ConfirmScheduleApproval(myDb))

	// Roster Endpoints
	protectedMux.HandleFunc("GET /roster/{class_id}", utils.RequireAdmin(roster.GetRoster(myDb)))
	protectedMux.HandleFunc("GET /roster/enrollment/{student_id}", utils.RequireAdmin(roster.GetStudentEnrollment(myDb)))
	protectedMux.HandleFunc("POST /roster/{class_id}/enroll", roster.EnrollStudent(myDb))

	// Enrollment Request Endpoints
	protectedMux.HandleFunc("POST /enrollment_requests", roster.CreateEnrollmentRequest(myDb))
	protectedMux.HandleFunc("GET /enrollment_requests", utils.RequireAdmin(roster.GetEnrollmentRequests(myDb)))
	protectedMux.HandleFunc("PATCH /enrollment_requests/{request_id}", utils.RequireAdmin(roster.UpdateEnrollmentRequest(myDb)))

	// Makeup Request Endpoints
	protectedMux.HandleFunc("GET /makeup_requests", utils.RequireAdmin(roster.GetMakeupRequests(myDb)))
	protectedMux.HandleFunc("GET /makeup_redemptions/available_dates/{class_id}", roster.GetAvailableRedemptionDates(myDb))
	protectedMux.HandleFunc("POST /makeup_requests", roster.CreateMakeupRequest(myDb))
	protectedMux.HandleFunc("PATCH /makeup_request/{request_id}", utils.RequireAdmin(roster.UpdateMakeupRequest(myDb)))
	protectedMux.HandleFunc("GET /makeup_redemptions", utils.RequireAdmin(roster.GetMakeupRedemptionRequests(myDb)))
	protectedMux.HandleFunc("POST /makeup_redemptions", roster.CreateMakeupRedemptionRequest(myDb))
	protectedMux.HandleFunc("PATCH /makeup_redemptions/{request_id}", utils.RequireAdmin(roster.UpdateMakeupRedemptionRequest(myDb)))

	// Calendar Endpoints
	protectedMux.HandleFunc("GET /calendarEvents", utils.RequireAdmin(class.GetCalendarEvents(myDb)))
	protectedMux.HandleFunc("GET /calendarEvents/{student_id}", class.GetCalendarEventsByStudent(myDb))

	// Student Endpoints
	protectedMux.HandleFunc("GET /students", utils.RequireAdmin(students.GetStudents(myDb)))
	protectedMux.HandleFunc("GET /students/enrollment", utils.RequireAdmin(students.GetStudentsWithEnrollment(myDb)))
	protectedMux.HandleFunc("POST /students", utils.RequireAdmin(students.AddStudent(myDb)))
	protectedMux.HandleFunc("PATCH /students/{student_id}", utils.RequireAdmin(students.UpdateStudent(myDb)))

	myDb.Logger.Printf("Server starting on :%v", port)
	if err := http.ListenAndServe(":"+port, wrappedMux); err != nil {
		log.Fatalf("Error initializing server: %v", err)
	}

}
