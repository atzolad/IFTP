package roster

import (
	"IFTP/db"
	"IFTP/timeutils"
	"IFTP/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

//	type Roster struct {
//		ID                int    `json:"id"`
//		Lecture_Date      string `json:"date"`
//		Student_ID        string `json:"student_id"`
//		Lecture_ID        string `json:"class_id"`
//		Registration_date string `json:"registration_date"`
//		Active            bool   `json:"active"`
//	}

type StudentRoster struct {
	ID     string       `db:"id" json:"id"`
	Name   string       `db:"name" json:"name"`
	Email  string       `db:"email" json:"email"`
	Status RosterStatus `db:"status" json:"status"`
}

type RosterStatus string

const (
	Enrolled RosterStatus = "Enrolled"
	AWAY     RosterStatus = "Away"
)

type GetRosterRequest struct {
	ClassName     string          `db:"class_name" json:"class_name"`
	Students      []StudentRoster `db:"students" json:"students"`
	EnrolledCount int             `db:"enrolled_count" json:"enrolled_count"`
	SessionDates  []time.Time     `db:"session_dates" json:"session_dates"`
}

type RosterRequest struct {
	ClassID   string    `json:"class_id"`
	Month     time.Time `json:"month"`
	ClassDate time.Time `json:"class_date"`
}

type EnrollmentRequest struct {
	StudentID  string   `json:"student_id"`
	ClassID    string   `json:"class_id"`
	ClassDates []string `json:"class_dates"`
}

type StudentEnrollment struct {
	ClassName string    `db:"class_name" json:"class_name"`
	ClassDate time.Time `db:"class_date" json:"class_date"`
	Month     time.Time `db:"month" json:"month"`
}

type EnrollmentRequestApproval struct {
	RequestID          string     `db:"id" json:"request_id"`
	StudentName        string     `db:"name" json:"name"`
	StudentEmail       string     `db:"email" json:"email"`
	CurrentlyEnrolled  []string   `db:"currently_enrolled" json:"currently_enrolled"`
	RequestedClassID   string     `db:"requested_class_id" json:"requested_class_id"`
	RequestedClassName string     `db:"class_name" json:"requested_class_name"`
	Month              *time.Time `db:"month" json:"month"`
	Teacher            string     `db:"teacher" json:"teacher"`
	Schedule           string     `db:"schedule" json:"schedule"`
	AvailableSpots     int        `db:"available_spots" json:"available"`
	Reason             string     `db:"reason" json:"reason"`
	RequestedAt        time.Time  `db:"requested_at" json:"requested_at"`
}

type EnrollmentRequestInput struct {
	RequestedClassID string     `json:"requested_class_id"`
	Reason           string     `json:"reason"`
	Month            *time.Time `json:"month"`
}

type MakeupRequestInput struct {
	ClassID            string   `json:"class_id"`
	MissedSessionDates []string `json:"missed_session_dates"`
	Reason             string   `json:"reason"`
}

type MakeupRequest struct {
	ID           string      `db:"id" json:"id"`
	StudentID    string      `db:"student_id" json:"student_id"`
	StudentName  string      `db:"name" json:"name"`
	StudentEmail string      `db:"email" json:"email"`
	ClassID      string      `db:"class_id" json:"class_id"`
	ClassName    string      `db:"class_name" json:"class_name"`
	MissedDates  []time.Time `db:"missed_session_dates" json:"missed_session_dates"`
	Reason       string      `db:"reason" json:"reason"`
	Status       string      `db:"status" json:"status"`
	RequestedAt  time.Time   `db:"requested_at" json:"requested_at"`
}

type MakeupRedemptionReq struct {
	RequestedClassID string `json:"requested_class_id"`
	RequestedDate    string `json:"requested_date"`
	Note             string `json:"note"`
}

type MakeupRedemption struct {
	ID                string    `db:"id" json:"id"`
	StudentID         string    `db:"student_id" json:"student_id"`
	StudentName       string    `db:"name" json:"name"`
	StudentEmail      string    `db:"email" json:"email"`
	RequestedClassID  string    `db:"requested_class_id" json:"requested_class_id"`
	ClassName         string    `db:"class_name" json:"class_name"`
	Teacher           string    `db:"teacher" json:"teacher"`
	RequestedDate     time.Time `db:"requested_date" json:"requested_date"`
	Note              string    `db:"note" json:"note"`
	Status            string    `db:"status" json:"status"`
	RequestedAt       time.Time `db:"requested_at" json:"requested_at"`
	CurrentlyEnrolled int       `db:"currently_enrolled" json:"currently_enrolled"`
	AvailableSpots    int       `db:"available_spots" json:"available_spots"`
}

// GetRoster responds with the overall enrolled class lists
func GetRoster(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		monthStr := r.FormValue("month")
		classDateStr := r.FormValue("class_date")
		classId := r.PathValue("class_id")

		myDb.Logger.Printf("month: %v, class_date: %v, class_id: %v", monthStr, classDateStr, classId)

		month, err := time.Parse("2006-01-02", monthStr)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Error: month required in YYYY-MM-DD format",
				Code:    http.StatusBadRequest,
			})
			return
		}

		classDate, err := time.Parse("2006-01-02", classDateStr)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Error: class date required in YYYY-MM-DD format",
				Code:    http.StatusBadRequest,
			})
			return
		}

		fullRoster, err := dbGetRoster(ctx, myDb, classId, month, classDate)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: fmt.Sprintf("Error retrieving roster from db: %v", err),
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Error retrieving roster from db: %v", err)
			return
		}

		utils.WriteJSONResponse(w, http.StatusOK, fullRoster)
		myDb.Logger.Printf("Successfully retrieved roster\n")
	}
}

func GetStudentEnrollment(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		monthStr := r.FormValue("month")
		studentId := r.PathValue("student_id")

		myDb.Logger.Printf("Get student enrollment request- month: %v, student_id: %v", monthStr, studentId)

		// studentId, err := strconv.Atoi(studentIdStr)
		// if err != nil {
		// 	utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
		// 		Status:  "error",
		// 		Message: "Error: Student id required",
		// 		Code:    http.StatusBadRequest,
		// 	})
		// 	return
		// }

		month, err := time.Parse("2006-01-02", monthStr)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Error: month required in YYYY-MM-DD format",
				Code:    http.StatusBadRequest,
			})
		}

		studentEnrollment, err := dbGetStudentEnrollment(ctx, myDb, studentId, month)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: fmt.Sprintf("Error retrieving roster from db: %v", err),
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Error retrieving roster from db: %v", err)
			return
		}

		utils.WriteJSONResponse(w, http.StatusOK, studentEnrollment)
		myDb.Logger.Printf("Successfully retrieved roster\n")
	}
}

func CreateEnrollmentRequest(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		// TODO get student id from the session- will hardcode it here for now.

		studentId, ok := r.Context().Value(utils.CtxUserID).(string)
		if !ok || studentId == "" {
			utils.WriteJSONResponse(w, http.StatusUnauthorized, utils.ResponseData{
				Status:  "error",
				Message: "Unauthorized",
				Code:    http.StatusUnauthorized,
			})
			return
		}

		var input EnrollmentRequestInput

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Error Decoding Request: %v",
				Code:    http.StatusBadRequest,
			})
			myDb.Logger.Printf("Error decoding Request: %v,", err)
			return
		}

		myDb.Logger.Printf("New enrollment request for student id: %v and class id: %v", studentId, input.RequestedClassID)

		var newEnrollmentRequest EnrollmentRequestApproval

		newEnrollmentRequest.RequestedClassID = input.RequestedClassID
		newEnrollmentRequest.Reason = strings.TrimSpace(input.Reason)
		newEnrollmentRequest.Month = input.Month

		err := dbGetStudentInfo(ctx, myDb, &newEnrollmentRequest, studentId)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Error fetching student info from db",
				Code:    http.StatusBadRequest,
			})
			myDb.Logger.Printf("Error fetching student info from db: %v", err)
			return
		}

		err = dbGetClassInfo(ctx, myDb, &newEnrollmentRequest)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Error fetching class info from db",
				Code:    http.StatusBadRequest,
			})
			myDb.Logger.Printf("Error fetching class info from db: %v", err)
			return
		}

		tx, err := myDb.Pool.Begin(ctx)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error Begining transcation",
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Error Begining transcation: %v", err)
			return
		}

		defer tx.Rollback(ctx)

		duplicate, err := dbEnrollmentReqExists(ctx, tx, &newEnrollmentRequest, studentId)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error checking db for duplicates",
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Error checking db for duplicates: %v", err)
			return
		}

		if duplicate {
			utils.WriteJSONResponse(w, http.StatusConflict, utils.ResponseData{
				Status:  "error",
				Message: fmt.Sprintf("Enrollment request for student %v and class %v already exists", newEnrollmentRequest.StudentName, newEnrollmentRequest.RequestedClassName),
				Code:    http.StatusConflict,
			})
			myDb.Logger.Printf("Enrollment request for student %v and class %v already exists", newEnrollmentRequest.StudentName, newEnrollmentRequest.RequestedClassName)
			return
		}

		alreadyEnrolled, err := dbStudentAlreadyEnrolled(ctx, tx, &newEnrollmentRequest, studentId)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error checking db for prior enrollment",
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Error checking db for prior enrollment: %v", err)
			return
		}

		if alreadyEnrolled {
			utils.WriteJSONResponse(w, http.StatusConflict, utils.ResponseData{
				Status:  "error",
				Message: fmt.Sprintf("Student %v already enrolled in class %v", newEnrollmentRequest.StudentName, newEnrollmentRequest.RequestedClassName),
				Code:    http.StatusConflict,
			})
			myDb.Logger.Printf("Student %v already enrolled in class %v", newEnrollmentRequest.StudentName, newEnrollmentRequest.RequestedClassName)
			return
		}

		err = dbInsertEnrollmentRequest(ctx, tx, &newEnrollmentRequest, studentId)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error creating new enrollment request in db",
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Error creating new enrollment request in db: %v", err)
			return

		}

		if err := tx.Commit(ctx); err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Failed to commit database changes",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		utils.WriteJSONResponse(w, http.StatusOK, "Successfully created enrollment request")
		myDb.Logger.Printf("Successfully created enrollment request for student id : %v and class id: %v", studentId, newEnrollmentRequest.RequestedClassID)

	}
}

func GetEnrollmentRequests(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requests, err := dbGetEnrollmentRequests(ctx, myDb)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Error fetching enrollment requests",
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Error fetching enrollment requests: %v", err)
			return
		}

		utils.WriteJSONResponse(w, http.StatusOK, requests)

	}
}

func UpdateEnrollmentRequest(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestId := r.PathValue("request_id")

		var input struct {
			Status string `json:"status"` //Approved or Denied
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Error decoding enrollment request for update",
				Code:    http.StatusBadRequest,
			})
			myDb.Logger.Printf("Error decoding enrollment request for update")
			return
		}

		if input.Status != "Approved" && input.Status != "Denied" {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Status must be approved or denied",
				Code:    http.StatusBadRequest,
			})
			myDb.Logger.Printf("Error updating enrollment request- status must be Approved or Denied: %v", input.Status)
			return
		}

		tx, err := myDb.Pool.Begin(ctx)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: fmt.Sprintf("Error Begining transcation: %v", err),
				Code:    http.StatusInternalServerError,
			})
			return
		}

		defer tx.Rollback(ctx)

		studentId, classId, month, err := dbUpdateRequestStatus(ctx, tx, requestId, input.Status)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error updating enrollment request status",
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Erorr updating enrollment request status: %v", err)
			return
		}
		if input.Status == "Approved" {
			if err := dbEnrollStudent(ctx, tx, studentId, classId, month); err != nil {
				utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
					Status:  "error",
					Message: "Error enrolling student",
					Code:    http.StatusInternalServerError,
				})
				myDb.Logger.Printf("Erorr enrolling student: %v", err)
				return
			}
		}
		if err := tx.Commit(ctx); err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error commiting enrollment request update and enrollment transaction",
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Error commiting enrollment request update and enrollment transaction: %v", err)
			return
		}
		utils.WriteJSONResponse(w, http.StatusOK, utils.ResponseData{
			Status:  "success",
			Message: fmt.Sprintf("Enrollment request %v successfully", input.Status),
			Code:    http.StatusOK,
		})
	}
}

func EnrollStudent(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		classId := strings.TrimSpace(r.PathValue("class_id"))

		// classId, err := strconv.Atoi(strings.TrimSpace(r.PathValue("class_id")))
		// if err != nil {
		// 	utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
		// 		Status:  "error",
		// 		Message: "Invalid class ID",
		// 		Code:    http.StatusBadRequest,
		// 	})
		// 	return
		// }

		// TODO: replace with student ID from session

		studentId, ok := r.Context().Value(utils.CtxUserID).(string)
		if !ok || studentId == "" {
			utils.WriteJSONResponse(w, http.StatusUnauthorized, utils.ResponseData{
				Status:  "error",
				Message: "Unauthorized",
				Code:    http.StatusUnauthorized,
			})
			return
		}

		// studentId := 1

		var body struct {
			Month string `json:"month"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Invalid request body",
				Code:    http.StatusBadRequest,
			})
			return
		}

		month, err := time.Parse("2006-01-02", body.Month)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Error: month required in YYYY-MM-DD format",
				Code:    http.StatusBadRequest,
			})
			return
		}

		tx, err := myDb.Pool.Begin(ctx)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error Begining transcation",
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Error Begining transcation: %v", err)
			return
		}
		defer tx.Rollback(ctx)

		if err := dbEnrollStudent(ctx, tx, studentId, classId, month); err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error enrolling student",
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Erorr enrolling student: %v", err)
			return
		}

		if err := tx.Commit(ctx); err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error commiting enrollment",
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Error commiting enrollment: %v", err)
			return
		}

		utils.WriteJSONResponse(w, http.StatusOK, utils.ResponseData{
			Status:  "success",
			Message: "Student successfully enrolled",
			Code:    http.StatusOK,
		})
	}
}

func CreateMakeupRequest(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		studentId, ok := r.Context().Value(utils.CtxUserID).(string)
		if !ok || studentId == "" {
			utils.WriteJSONResponse(w, http.StatusUnauthorized, utils.ResponseData{
				Status:  "error",
				Message: "Unauthorized",
				Code:    http.StatusUnauthorized,
			})
			return
		}

		var makeupRequest MakeupRequestInput

		if err := json.NewDecoder(r.Body).Decode(&makeupRequest); err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Error Decoding Request",
				Code:    http.StatusBadRequest,
			})
			myDb.Logger.Printf("Error decoding Request: %v,", err)
			return
		}

		if makeupRequest.ClassID == "" {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Class ID is required",
				Code:    http.StatusBadRequest,
			})
			return
		}

		if len(makeupRequest.MissedSessionDates) == 0 {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Must provide at least one session date",
				Code:    http.StatusBadRequest,
			})
			return
		}

		parsedDates := make([]time.Time, 0, len(makeupRequest.MissedSessionDates))
		for _, dateStr := range makeupRequest.MissedSessionDates {
			date, err := timeutils.ParseDate(dateStr)
			if err != nil {
				utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
					Status:  "error",
					Message: "Error parsing missed session dates",
					Code:    http.StatusBadRequest,
				})
				myDb.Logger.Printf("Error parsing missed session dates: %v,", err)
				return
			}
			parsedDates = append(parsedDates, date)
		}

		tx, err := myDb.Pool.Begin(ctx)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error Begining transcation",
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Error Begining transcation: %v", err)
			return
		}

		defer tx.Rollback(ctx)

		if err := dbInsertMakeupRequest(ctx, tx, studentId, &makeupRequest, parsedDates); err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error creating makeup request",
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Error inserting makeup request: %v", err)
			return
		}

		if err := tx.Commit(ctx); err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Failed to commit makeup database transaction",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		utils.WriteJSONResponse(w, http.StatusOK, utils.ResponseData{
			Status:  "success",
			Message: "Makeup Request submitted successfully",
			Code:    http.StatusOK,
		})
		myDb.Logger.Printf("Successfully created makeup request for student id %v in class id %v", studentId, makeupRequest.ClassID)
	}
}

func GetMakeupRequests(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requests, err := dbGetMakeupRequests(ctx, myDb)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error fetching makeup requests",
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Error fetching makeup requests: %v", err)
			return
		}
		utils.WriteJSONResponse(w, http.StatusOK, requests)
	}
}

func UpdateMakeupRequest(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestId := r.PathValue("request_id")

		var input struct {
			Status string `json:"status"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Error decoding makeup request body",
				Code:    http.StatusBadRequest,
			})
			myDb.Logger.Printf("Error decoding makeup request body: %v", err)
			return
		}

		if input.Status != "Approved" && input.Status != "Denied" {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Makeup approval status must be Approved or Denied",
				Code:    http.StatusBadRequest,
			})
			return
		}

		tx, err := myDb.Pool.Begin(ctx)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error Begining transcation",
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Error Begining transcation: %v", err)
			return
		}

		defer tx.Rollback(ctx)

		if err := dbUpdateMakeupRequestStatus(ctx, tx, requestId, input.Status); err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error updating makeup request",
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Error updating makeup request: %v", err)
			return
		}

		if err := tx.Commit(ctx); err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Failed to commit makeup database transaction",
				Code:    http.StatusInternalServerError,
			})
			return
		}
		utils.WriteJSONResponse(w, http.StatusOK, utils.ResponseData{
			Status:  "success",
			Message: "Makeup Request status updated successfully",
			Code:    http.StatusOK,
		})
		myDb.Logger.Printf("Successfully updated makeup request with id %v to %v", requestId, input.Status)
	}
}

func CreateMakeupRedemptionRequest(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		studentId, ok := r.Context().Value(utils.CtxUserID).(string)
		if !ok || studentId == "" {
			utils.WriteJSONResponse(w, http.StatusUnauthorized, utils.ResponseData{
				Status:  "error",
				Message: "Unauthorized",
				Code:    http.StatusUnauthorized,
			})
			return
		}

		var redemptionRequest MakeupRedemptionReq

		if err := json.NewDecoder(r.Body).Decode(&redemptionRequest); err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Error Decoding Request",
				Code:    http.StatusBadRequest,
			})
			myDb.Logger.Printf("Error decoding Request: %v,", err)
			return
		}

		if redemptionRequest.RequestedClassID == "" || redemptionRequest.RequestedDate == "" {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Class ID and class date are required",
				Code:    http.StatusBadRequest,
			})
			return
		}

		parsedDate, err := timeutils.ParseDate(redemptionRequest.RequestedDate)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Error parsing requested session date",
				Code:    http.StatusBadRequest,
			})
			myDb.Logger.Printf("Error parsing missed session dates: %v,", err)
			return
		}

		tx, err := myDb.Pool.Begin(ctx)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error Begining transcation",
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Error Begining transcation: %v", err)
			return
		}

		defer tx.Rollback(ctx)

		if err := dbInsertMakeupRedemptionRequest(ctx, tx, studentId, &redemptionRequest, parsedDate); err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error creating makeup redemption request",
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Error inserting makeup redemption request: %v", err)
			return
		}

		if err := tx.Commit(ctx); err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Failed to commit makeup database transaction",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		utils.WriteJSONResponse(w, http.StatusOK, utils.ResponseData{
			Status:  "success",
			Message: "Makeup Request submitted successfully",
			Code:    http.StatusOK,
		})
		myDb.Logger.Printf("Successfully created makeup redemption request for student id %v in class id %v", studentId, redemptionRequest.RequestedClassID)

	}
}

func GetMakeupRedemptionRequests(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requests, err := dbGetMakeupRedemptionRequests(ctx, myDb)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error fetching makeup redemption requests",
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Error fetching makeup redemption requests: %v", err)
			return
		}
		utils.WriteJSONResponse(w, http.StatusOK, requests)
	}
}

func UpdateMakeupRedemptionRequest(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestId := r.PathValue("request_id")

		var input struct {
			Status string `json:"status"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Error decoding makeup request body",
				Code:    http.StatusBadRequest,
			})
			myDb.Logger.Printf("Error decoding makeup request body: %v", err)
			return
		}

		if input.Status != "Approved" && input.Status != "Denied" {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Makeup approval status must be Approved or Denied",
				Code:    http.StatusBadRequest,
			})
			return
		}

		tx, err := myDb.Pool.Begin(ctx)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error Begining transcation",
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Error Begining transcation: %v", err)
			return
		}

		defer tx.Rollback(ctx)

		if err := dbUpdateMakeupRedemptionRequestStatus(ctx, tx, requestId, input.Status); err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error updating makeup redemption request",
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Error updating makeup redemption request: %v", err)
			return
		}

		if err := tx.Commit(ctx); err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Failed to commit makeup redemption database transaction",
				Code:    http.StatusInternalServerError,
			})
			return
		}
		utils.WriteJSONResponse(w, http.StatusOK, utils.ResponseData{
			Status:  "success",
			Message: "Makeup Redemption Request status updated successfully",
			Code:    http.StatusOK,
		})
		myDb.Logger.Printf("Successfully updated makeup redemption request with id %v to %v", requestId, input.Status)
	}
}

// Enroll adds the student info in the body of the request to the class from the url.
// TODO(): figure out how we want to require full month of classes for students
// func Enroll(myDb *db.MyDatabase) gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		var newEnrollmentRequest enrollmentRequest

// 		// Call BindJSON to bind the received JSON to
// 		// newEnrollment.
// 		if err := c.BindJSON(&newEnrollmentRequest); err != nil {
// 			c.JSON(http.StatusBadRequerst, gin.H{"error": err.Error()})
// 			return
// 		}

// 		// Retrieve the class id from the url and assign the integer value to the newEnrollmentRequest struct
// 		classID := c.Param("class_id")
// 		fmt.Println(classID)
// 		intClassID, err := strconv.Atoi(classID)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		newEnrollmentRequest.ClassID = intClassID

// 		convertedDates, err := convertStrDT(newEnrollmentRequest.ClassDates)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		for _, date := range convertedDates {
// 			if err := dbEnroll(myDb, newEnrollmentRequest.ClassID, date, newEnrollmentRequest.StudentID); err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 				return
// 			}

// 		}
// 		c.JSON(http.StatusCreated, newEnrollmentRequest)
// 		fmt.Printf("Successfully enrolled student: %v into class id %v for dates: %v ",
// 			newEnrollmentRequest.StudentID, newEnrollmentRequest.ClassID, newEnrollmentRequest.ClassDates)
// 	}
// }

// func convertStrDT(dates []string) ([]time.Time, error) {
// 	convertedDates := make([]time.Time, len(dates))

// 	for i, date := range dates {

// 		// Validate that the date provided is in the correct format
// 		parsedDate, err := time.Parse("2006-01-02", date)
// 		if err != nil {
// 			return nil, fmt.Errorf("error occured during datetime conversion : %v", err)
// 		}
// 		convertedDates[i] = parsedDate
// 	}
// 	return convertedDates, nil
// }
