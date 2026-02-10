package students

import (
	"IFTP/db"
	"IFTP/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Student struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Email           string   `json:"email"`
	Active          bool     `json:"active"`
	EnrolledClasses []string `json:"enrolledClasses"`
}

func GetStudents(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		students, err := dbRetrieveStudents(myDb)
		if err != nil {
			fmt.Printf("Error retrieving students from DB")
			utils.WriteJSONResponse(w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSONResponse(w, http.StatusOK, students)
		fmt.Printf("Successfully retrieved student list \n")
	}
}

func GetStudentsWithEnrollment(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		students, err := dbGetStudentsWithEnrollment(myDb)
		if err != nil {
			fmt.Printf("Error retrieving students from DB")
			utils.WriteJSONResponse(w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSONResponse(w, http.StatusOK, students)
		fmt.Printf("Successfully retrieved student list \n")
	}
}

// AddStudent adds a user from JSON received in the request body.
func AddStudent(myDb *db.MyDatabase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newStudent Student

		// Call BindJSON to bind the received JSON to
		// newStudent.
		if err := c.BindJSON(&newStudent); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := InsertStudent(myDb, &newStudent); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, newStudent)
		fmt.Printf("Successfully created new student: %v", newStudent)
	}
}

// Update Student updates the student details based on the JSON received in the request body.
func UpdateStudent(myDb *db.MyDatabase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		integerID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		var updateStudent Student

		// Call BindJSON to bind the received JSON to updatedStudent
		if err := c.BindJSON(&updateStudent); err != nil {
			return
		}

		returnedStudent, err := UpdateStudentDB(myDb, integerID, &updateStudent)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("content-type", "application/json")
		c.JSON(http.StatusOK, returnedStudent)
		fmt.Printf("Successfully updated student: %v \n", returnedStudent.Name)
	}
}

// SoftDeleteStudent changes the Active status of the student to false, rather than permanently deleting.
func SoftDeleteStudent(myDb *db.MyDatabase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		integerID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student id"})
		}

		deletedStudent, err := SoftDeleteStudentDB(myDb, integerID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("student %v deleted successfully", deletedStudent)})
		fmt.Printf("Successfully soft deleted student %v with id: %v \n", deletedStudent, integerID)
	}
}
