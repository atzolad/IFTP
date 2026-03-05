package roster

import (
	"IFTP/db"
	"IFTP/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

// GetRoster responds with the overall enrolled class lists
func GetRoster(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		monthStr := r.FormValue("month")
		classDateStr := r.FormValue("class_date")
		classIdStr := r.PathValue("class_id")

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
			log.Printf("Error retrieving roster from db: %v", err)
			return
		}

		utils.WriteJSONResponse(w, http.StatusOK, fullRoster)
		log.Printf("Successfully retrieved roster\n")
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
