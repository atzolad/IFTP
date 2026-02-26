package utils

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
	// "cloud.google.com/go/storage"
	// ics "github.com/arran4/golang-ical"
	// "github.com/google/uuid"
	// "github.com/markbates/goth"
	// "github.com/markbates/goth/gothic"
	// "golang.org/x/oauth2"
	// "golang.org/x/oauth2/google"
	// "cloud.google.com/go/storage"
	// ics "github.com/arran4/golang-ical"
	// "github.com/google/uuid"
	// "github.com/markbates/goth"
	// "github.com/markbates/goth/gothic"
	// "golang.org/x/oauth2"
	// "golang.org/x/oauth2/google"
)

type Templates struct {
	Index        *template.Template
	Unauthorized *template.Template
}

type ResponseData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func LoadTemplates() Templates {
	tpl := Templates{
		Index: template.Must(template.ParseFiles("templates/index.html")),
	}
	return tpl
}

func IndexHandler(tpl Templates, baseUrl string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data := map[string]interface{}{
			"BaseUrl": baseUrl,
		}

		err := tpl.Index.Execute(w, data)
		if err != nil {
			renderError := fmt.Errorf("error rendering main template: %v", err)
			log.Println(renderError)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

type LogWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lw *LogWriter) WriteHeader(code int) {
	lw.statusCode = code
	lw.ResponseWriter.WriteHeader(code)
}

func (lw *LogWriter) Unwrap() http.ResponseWriter { return lw.ResponseWriter }

// Log the connection details for the API
func LoggingWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request details
		log.Printf("Request Type : %v to %v From %v", r.Method, r.URL.Path, r.RemoteAddr)
		lw := &LogWriter{w, http.StatusOK}
		next.ServeHTTP(lw, r)
		//Log the Reponse details
		log.Printf("Response: %s %s %s - Status: %d", r.Method, r.URL.Path, r.Body, lw.statusCode)
	})
}

func WriteJSONResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	//Encode data to response writer
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		// Handle Encoding errors
		log.Printf("failed to encode JSON response: %v", err)
	}

}

func ConvertStrDT(date string) (time.Time, error) {

	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Time{}, fmt.Errorf("error occured during datetime conversion : %v", err)
	}
	return parsedDate, nil
}

// func ConfigureGoogleOauth(client_secret_file string) *oauth2.Config {

// 	// Read the JSON file content
// 	jsonKey, err := os.ReadFile(client_secret_file)
// 	if err != nil {
// 		log.Fatalf("Unable to read client secret file: %v", err)
// 	}

// 	// Use ConfigFromJSON to create an oauth2.Config
// 	// You must specify the required scopes for your application
// 	scopes := "openid email profile"
// 	// The function returns a *oauth2.Config
// 	config, err := google.ConfigFromJSON(jsonKey, scopes)
// 	if err != nil {
// 		log.Fatalf("Unable to parse client secret file to config: %v", err)
// 	}
// 	return config

// }

// func GoogleLogin() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Manually add provider to query so gothic knows to use google. Can send this from the front-end later.
// 		q := r.URL.Query()
// 		q.Add("provider", "google")
// 		q.Add("prompt", "consent")
// 		r.URL.RawQuery = q.Encode()

// 		gothic.BeginAuthHandler(w, r)
// 	}
// }

// func HandleGoogleOauth(myDb *db.MyDatabase) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		// TODO send this from front-end later.
// 		q := r.URL.Query()
// 		q.Add("provider", "google")
// 		r.URL.RawQuery = q.Encode()

// 		user, err := gothic.CompleteUserAuth(w, r)
// 		if err != nil {
// 			log.Printf("Error completing user auth: %v", err)
// 			WriteJSONResponse(w, http.StatusInternalServerError, ResponseData{
// 				Status:  "error",
// 				Message: "Error completing user auth",
// 				Code:    http.StatusInternalServerError,
// 			})
// 			log.Printf("Error completing user auth")
// 			return
// 		}

// 		if !IsAuthorizedUser(user.Email) {
// 			redirectURL := fmt.Sprintf("/unauthorized?email=%s", url.QueryEscape(user.Email))
// 			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
// 			return
// 		}

// 		// store the user session
// 		session, err := gothic.Store.New(r, "goth_session")
// 		if err != nil {
// 			log.Printf("Error storing user session: %v", err)
// 			WriteJSONResponse(w, http.StatusInternalServerError, ResponseData{
// 				Status:  "error",
// 				Message: "Error storing user session",
// 				Code:    http.StatusInternalServerError,
// 			})
// 			return
// 		}

// 		// Save user info into the database
// 		// if err := DbInsertOrUpdateUser(myDb, user); err != nil {
// 		// 	log.Printf("Error updating user in db: %v", err)
// 		// }

// 		// Save user info in cookie- The entire user object is too big.
// 		session.Values["userId"] = user.UserID
// 		session.Values["email"] = user.Email
// 		session.Values["name"] = user.Name

// 		// save the user session
// 		if err = session.Save(r, w); err != nil {
// 			log.Printf("Error saving user session: %v", err)
// 			WriteJSONResponse(w, http.StatusInternalServerError, ResponseData{
// 				Status:  "error",
// 				Message: "Error saving user session",
// 				Code:    http.StatusInternalServerError,
// 			})
// 			return
// 		}

// 		fmt.Printf("User login email: %sv, User Name: %v,  Access token: %v", user.Name, user.Email, user.AccessToken)

// 		http.Redirect(w, r, "/", http.StatusFound)

// 	}
// }

// func GetUserAuth(myDb *db.MyDatabase) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Retrieve the session
// 		session, err := gothic.Store.Get(r, "goth_session")
// 		if err != nil {
// 			log.Printf("Error retrieving user session: %v", err)
// 			WriteJSONResponse(w, http.StatusUnauthorized, ResponseData{
// 				Status:  "error",
// 				Message: "Error retrieving user session",
// 				Code:    http.StatusUnauthorized,
// 			})
// 			return
// 		}

// 		// Get user data from session
// 		user := session.Values["user"]
// 		if user == nil {
// 			log.Printf("Error: no authenticated user found: %v", err)
// 			WriteJSONResponse(w, http.StatusUnauthorized, ResponseData{
// 				Status:  "error",
// 				Message: "No authenticated user found",
// 				Code:    http.StatusUnauthorized,
// 			})
// 			return
// 		}

// 		WriteJSONResponse(w, http.StatusOK, ResponseData{
// 			Status:  "success",
// 			Message: "User fetched successfully",
// 			Code:    http.StatusOK,
// 		})
// 	}
// }

// func DbInsertOrUpdateUser(myDb *db.MyDatabase, user goth.User) error {
// 	userName := fmt.Sprintf("%v %v", user.FirstName, user.LastName)
// 	//TODO encrypt the refresh token before inserting into db.

// 	_, err := myDb.Db.Exec(`
// 			INSERT INTO users (google_id, name, email, access_token, refresh_token, token_expiry)
// 			VALUES ($1, $2, $3, $4, $5, $6)
// 			ON CONFLICT (provider_id) DO UPDATE SET
// 				email = EXCLUDED.email,
// 				name = EXCLUDED.name,
// 				access_token = EXCLUDED.access_token,
// 				token_expiry = EXCLUDED.token_expiry,
// 				refresh_token = COALESCE(NULLIF(EXCLUDED.refresh_token, ''), users.refresh_token);
// 		`, user.UserID, userName, user.Email, user.AccessToken, user.RefreshToken, user.ExpiresAt)

// 	if err != nil {
// 		return fmt.Errorf("failed to update use in db : %v", err)
// 	}
// 	return err
// }

// func RequireAuth(next http.Handler) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		apiKey := r.Header.Get("X-API-Key")
// 		if apiKey != "" {
// 			if apiKey == os.Getenv("API_KEY") {

// 				// Set a context value so handlers know this is API key auth
// 				ctx := context.WithValue(r.Context(), "user", "api_user")
// 				ctx = context.WithValue(ctx, "auth_type", "api_key")
// 				next.ServeHTTP(w, r.WithContext(ctx))
// 				return
// 			}
// 			// Invalid API key - return JSON error for API clients
// 			WriteJSONResponse(w, http.StatusUnauthorized, ResponseData{
// 				Status:  "error",
// 				Message: "Invalid API key",
// 				Code:    http.StatusUnauthorized,
// 			})
// 			return
// 		}

// 		// gets the user session from the request
// 		session, err := gothic.Store.Get(r, "goth_session")
// 		if err != nil {
// 			http.Redirect(w, r, "/auth/google", http.StatusSeeOther)
// 			return
// 		}
// 		// gets the user from the session data
// 		userId := session.Values["userId"]
// 		if userId == nil {
// 			http.Redirect(w, r, "/auth/google", http.StatusSeeOther)
// 			return
// 		}

// 		ctx := context.WithValue(r.Context(), "userId", userId)
// 		ctx = context.WithValue(ctx, "userEmail", session.Values["email"])
// 		ctx = context.WithValue(ctx, "userName", session.Values["name"])
// 		ctx = context.WithValue(ctx, "auth_type", "session")
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	}
// }

// func HandleLogout() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("Logout handler!")

// 		// Clear the session
// 		gothic.Logout(w, r)
// 		session, err := gothic.Store.Get(r, "goth_session")
// 		if err == nil {
// 			session.Options.MaxAge = -1
// 			session.Save(r, w)
// 		}
// 		http.Redirect(w, r, "/", http.StatusSeeOther)
// 	}
// }

// func IsAuthorizedUser(email string) bool {
// 	authorizedUserList := os.Getenv("AUTHORIZED_USERS")
// 	if authorizedUserList == "" {
// 		log.Println("Authorized users env variable missing")
// 		return false
// 	}

// 	sanitizedAuthorizedUserList := strings.TrimSpace(strings.ToLower(authorizedUserList))
// 	authorizedUsers := strings.Split(sanitizedAuthorizedUserList, ",")
// 	sanitizedEmail := strings.TrimSpace(strings.ToLower(email))

// 	fmt.Printf("SanitizedAuthorizedUserList: %v \n", sanitizedAuthorizedUserList)
// 	fmt.Printf("authorizedUsers: %v \n", authorizedUsers)
// 	fmt.Printf("sanitizedEmail: %v \n", sanitizedEmail)

// 	for _, user := range authorizedUsers {
// 		fmt.Println(user)
// 	}
// 	return slices.Contains(authorizedUsers, sanitizedEmail)
// }

// type GoogleUserEmail struct {
// 	Email string
// }

// func UnAuthorizedHandler(tpl Templates) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		userEmail := r.URL.Query().Get("email")

// 		requestingUser := GoogleUserEmail{
// 			Email: userEmail,
// 		}

// 		w.WriteHeader(http.StatusForbidden)
// 		err := tpl.Unauthorized.Execute(w, requestingUser)
// 		if err != nil {
// 			renderError := fmt.Errorf("error rendering unauthorized template: %v", err)
// 			log.Println(renderError)
// 			return
// 		}
// 	}
// }

// const (
// 	dateFormatUtc = "20060102"

// 	propertyDtStart ics.Property = "DTSTART;VALUE=DATE"
// 	propertyDtEnd   ics.Property = "DTEND;VALUE=DATE"

// 	componentPropertyDtStart = ics.ComponentProperty(propertyDtStart)
// 	componentPropertyDtEnd   = ics.ComponentProperty(propertyDtEnd)
// )

// func CreateCalendarInvite(start time.Time, end time.Time, technicianEmail string, jobName string, technicianName string, eventType common.HoldStatus) *ics.Calendar {

// 	cal := ics.NewCalendar()
// 	cal.SetMethod(ics.MethodRequest) // Essential for invites

// 	// Add an event with a unique identifier (UID)
// 	event := cal.AddEvent(uuid.New().String())
// 	event.SetProperty(componentPropertyDtStart, start.Format(dateFormatUtc))
// 	event.SetProperty(componentPropertyDtEnd, end.Format(dateFormatUtc))

// 	// Set event details (use UTC for global events or VTIMEZONE for localized events)
// 	event.SetCreatedTime(time.Now())
// 	event.SetDtStampTime(time.Now())
// 	event.SetModifiedAt(time.Now())

// 	event.SetSummary(fmt.Sprintf(`%v - %v`, jobName, eventType))
// 	event.SetDescription(fmt.Sprintf(`%v - %v`, jobName, eventType))

// 	//Define the Organizer TODO
// 	event.SetOrganizer("azoladswe@gmail.com", ics.WithCN("Cinemoves Scheduler"))

// 	//Add an Attendee (required participant) TODO
// 	event.AddAttendee(
// 		technicianEmail,
// 		ics.CalendarUserTypeIndividual,
// 		ics.ParticipationStatusNeedsAction,
// 		ics.ParticipationRoleReqParticipant,
// 		ics.WithCN(technicianName),
// 		ics.WithRSVP(true),
// 	)

// 	//Serialize the calendar to an ICS string
// 	// icsOutput := cal.Serialize()

// 	// fmt.Println(icsOutput)
// 	return cal
// }

// func HealthCheckHandler(myDb *db.MyDatabase) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		{

// 			if err := myDb.Db.Ping(); err != nil {
// 				w.WriteHeader(http.StatusServiceUnavailable)
// 				w.Write([]byte("Database unavailable"))
// 				return
// 			}

// 			w.WriteHeader(http.StatusOK)
// 			w.Write([]byte("OK"))
// 		}
// 	}
// }

// func SetupBucketCORS(ctx context.Context, bucket *storage.BucketHandle) error {

// 	// Get the domain from Railway env, default to localhost for dev
// 	domain := os.Getenv("RAILWAY_PUBLIC_DOMAIN")

// 	origins := []string{"http://localhost:8080"}

// 	if domain != "" {
// 		origins = append(origins, "https://"+domain)
// 	}

// 	policy := []storage.CORS{
// 		{
// 			MaxAge:          3600,
// 			Methods:         []string{"GET", "POST", "PUT"},
// 			Origins:         origins,
// 			ResponseHeaders: []string{"Content-Type"},
// 		},
// 	}

// 	// Apply the CORS configuration to the bucket
// 	if _, err := bucket.Update(ctx, storage.BucketAttrsToUpdate{CORS: policy}); err != nil {
// 		return err
// 	}
// 	return nil
// }
