package roster

import (
	"IFTP/db"
	"IFTP/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/spiffe/go-spiffe/v2/logger"
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
	ID     int          `db:"id" json:"id"`
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
	ClassID   int       `json:"class_id"`
	Month     time.Time `json:"month"`
	ClassDate time.Time `json:"class_date"`
}

type EnrollmentRequest struct {
	StudentID  int      `json:"student_id"`
	ClassID    int      `json:"class_id"`
	ClassDates []string `json:"class_dates"`
}

type StudentEnrollment struct {
	ClassName string    `db:"class_name" json:"class_name"`
	ClassDate time.Time `db:"class_date" json:"class_date"`
	Month     time.Time `db:"month" json:"month"`
}

type EnrollmentRequestApproval struct {
	RequestID         int      `json:"request_id"`
	StudentName       string   `json:"name"`
	StudentEmail      string   `json:"email"`
	CurrentlyEnrolled []string `json:"currently_enrolled"`
	RequestedClassID  int      `json:"requested_class_id"`
	RequestedClassName string `json:"requested_class_name"`
	Month *time.Time  `json:"month"`
	Teacher           string   `json:"teacher"`
	AvailableSpots    int      `json:"available"`
	Reason            string   `json:"reason"`
	RequestedAt        time.Time   `json:"requested_at"`
}

type EnrollmentRequestInput struct {
	RequestedClassID int    `json:"requested_class_id"`
	Reason           string `json:"reason"`
	Month         *time.Time  `json:"month"`
}

// GetRoster responds with the overall enrolled class lists
func GetRoster(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		monthStr := r.FormValue("month")
		classDateStr := r.FormValue("class_date")
		classIdStr := r.PathValue("class_id")

		myDb.Logger.Printf("month: %v, class_date: %v, class_id: %v", monthStr, classDateStr, classIdStr)

		classId, err := strconv.Atoi(classIdStr)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Error: Class id required",
				Code:    http.StatusBadRequest,
			})
			return
		}

		month, err := time.Parse("2006-01-02", monthStr)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Error: month required in YYYY-MM-DD format",
				Code:    http.StatusBadRequest,
			})
		}

		classDate, err := time.Parse("2006-01-02", classDateStr)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Error: class date required in YYYY-MM-DD format",
				Code:    http.StatusBadRequest,
			})
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
		studentIdStr := r.PathValue("student_id")

		myDb.Logger.Printf("Get student enrollment request- month: %v, student_id: %v", monthStr, studentIdStr)

		studentId, err := strconv.Atoi(studentIdStr)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Error: Student id required",
				Code:    http.StatusBadRequest,
			})
			return
		}

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

		studentId := 1
		var input EnrollmentRequestInput

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Error Decoding Request: %v",
				Code:    http.StatusBadRequest,
			})
			myDb.logger.Printf("Error decoding Request: %v,", err)
			return
		}

		myDb.logger.Printf("New enrollment request for student id: %v and class id: %v", studentId, input.RequestedClassID)

		var newEnrollmentRequest EnrollmentRequestApproval

		newEnrollmentRequest.RequestedClassID = input.RequestedClassID
		newEnrollmentRequest.Reason = input.Reason
		newEnrollmentRequest.Month = input.Month

		err := dbGetStudentInfo(ctx, myDb, &newEnrollmentRequest, studentId)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Error fetching student info from db",
				Code:    http.StatusBadRequest,
			})
			myDb.logger.Printf("Error fetching student info from db: %v", err)
			return
		}

		err = dbGetClassInfo(ctx, myDb, &newEnrollmentRequest)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Error fetching class info from db",
				Code:    http.StatusBadRequest,
			})
			myDb.logger.Printf("Error fetching class info from db: %v", err)
			return
	}

		tx, err := myDb.Pool.Begin(ctx)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error Begining transcation",
				Code:    http.StatusInternalServerError,
			})
			myDb.logger.Printf("Error Begining transcation: %v", err)
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
			myDb.logger.Printf("Error checking db for duplicates: %v", err)
			return
		}

		if duplicate {
			utils.WriteJSONResponse(w, http.StatusConflict, utils.ResponseData{
				Status:  "error",
				Message: fmt.Sprintf("Enrollment request for student %v and class %v already exists", newEnrollmentRequest.StudentName, newEnrollmentRequest.RequestedClassName),
				Code:    http.StatusConflict,
			})
			myDb.logger.Printf("Enrollment request for student %v and class %v already exists", newEnrollmentRequest.StudentName, newEnrollmentRequest.RequestedClassName)
			return
		}

		alreadyEnrolled, err := dbStudentAlreadyEnrolled(ctx, tx, &newEnrollmentRequest, studentId)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error checking db for prior enrollment",
				Code:    http.StatusInternalServerError,
			})
			myDb.logger.Printf("Error checking db for prior enrollment: %v", err)
			return
		}


		if alreadyEnrolled {
			utils.WriteJSONResponse(w, http.StatusConflict, utils.ResponseData{
				Status:  "error",
				Message: fmt.Sprintf("Student %v already enrolled in class %v", newEnrollmentRequest.StudentName, newEnrollmentRequest.RequestedClassName),
				Code:    http.StatusConflict,
			})
			myDb.logger.Printf("Student %v already enrolled in class %v", newEnrollmentRequest.StudentName, newEnrollmentRequest.RequestedClassName)
			return
		}

		err:= dbInsertEnrollmentRequest(ctx, tx, &newEnrollmentRequest, studentId) 
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error creating new enrollment request in db",
				Code:    http.StatusInternalServerError,
			})
			myDb.logger.Printf("Error creating new enrollment request in db: %v", err)
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
		myDb.Logger.Printf("Successfully created enrollment request for student id : %v and class id: %v",studentId, newEnrollmentRequest.RequestedClassID)
		
}

// Enroll adds the student info in the body of the request to the class from the url.
// TODO(): figure out how we want to require full month of classes for students
// func Enroll(myDb *db.MyDatabase) gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		var newEnrollmentRequest enrollmentRequest

// 		// Call BindJSON to bind the received JSON to
// 		// newEnrollment.
// 		if err := c.BindJSON(&newEnrollmentRequest); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
