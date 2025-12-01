package roster

import (
	"IFTP/db"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

//	type Roster struct {
//		ID                int    `json:"id"`
//		Lecture_Date      string `json:"date"`
//		Student_ID        string `json:"student_id"`
//		Lecture_ID        string `json:"class_id"`
//		Registration_date string `json:"registration_date"`
//		Active            bool   `json:"active"`
//	}
type enrollmentRequest struct {
	StudentID  int      `json:"student_id"`
	ClassID    int      `json:"class_id"`
	ClassDates []string `json:"class_dates"`
}

// // GetRoster responds with the overall enrolled class lists
// func GetRoster(myDb *db.MyDatabase) gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		fullRoster, err := GetRoster(myDb)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}

// 		c.Header("content-type", "application/json")
// 		c.JSON(http.StatusOK, fullRoster)
// 		fmt.Printf("Successfully retrieved roster \n")
// 	}
// }

// Enroll adds the student info in the body of the request to the class from the url.
// TODO(): figure out how we want to require full month of classes for students
func Enroll(myDb *db.MyDatabase) gin.HandlerFunc {
	return func(c *gin.Context) {

		var newEnrollmentRequest enrollmentRequest

		// Call BindJSON to bind the received JSON to
		// newEnrollment.
		if err := c.BindJSON(&newEnrollmentRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Retrieve the class id from the url and assign the integer value to the newEnrollmentRequest struct
		classID := c.Param("class_id")
		fmt.Println(classID)
		intClassID, err := strconv.Atoi(classID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newEnrollmentRequest.ClassID = intClassID

		convertedDates, err := convertStrDT(newEnrollmentRequest.ClassDates)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, date := range convertedDates {
			if err := dbEnroll(myDb, newEnrollmentRequest.ClassID, date, newEnrollmentRequest.StudentID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

		}
		c.JSON(http.StatusCreated, newEnrollmentRequest)
		fmt.Printf("Successfully enrolled student: %v into class id %v for dates: %v ",
			newEnrollmentRequest.StudentID, newEnrollmentRequest.ClassID, newEnrollmentRequest.ClassDates)
	}
}

func convertStrDT(dates []string) ([]time.Time, error) {
	convertedDates := make([]time.Time, len(dates))

	for i, date := range dates {

		// Validate that the date provided is in the correct format
		parsedDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			return nil, fmt.Errorf("error occured during datetime conversion : %v", err)
		}
		convertedDates[i] = parsedDate
	}
	return convertedDates, nil
}
