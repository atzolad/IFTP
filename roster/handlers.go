package roster

import (
	"IFTP/db"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Roster struct {
	ID          int    `json:"id"`
	Student_ID  string `json:"student_id"`
	Class_ID    string `json:"class_id"`
	Enrolled_At string `json:"enrolled_at"`
	Active      bool   `json:"active"`
}

// GetRoster responds with the overall enrolled class lists
func GetRoster(myDb *db.MyDatabase) gin.HandlerFunc {
	return func(c *gin.Context) {

		fullRoster, err := GetRoster(myDb)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("content-type", "application/json")
		c.JSON(http.StatusOK, fullRoster)
		fmt.Printf("Successfully retrieved roster \n")
	}
}

// JoinClass adds the student to the class from the request.
func JoinClass(myDb *db.MyDatabase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newEnrollment Roster

		// Call BindJSON to bind the received JSON to
		// newStudent.
		if err := c.BindJSON(&newEnrollment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := JoinClass(myDb, &newEnrollment); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, newEnrollment)
		fmt.Printf("Successfully created new student: %v", newEnrollment)
	}
}
