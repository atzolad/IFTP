package student

import (
	"IFTP/db"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Student struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Paid   bool   `json:"paid"`
	Active bool   `json:"active"`
}

// var students = []Student{
// 	{Name: "Alex Zolad", Email: "atzolad@gmail.com", Paid: false, Active: true},
// 	{Name: "Megan Chang", Email: "meganchang10@gmail.com", Paid: false, Active: true},
// 	{Name: "Sarah Vaughan", Email: "Sarahvaughan@gmail.com", Paid: false, Active: true},
// }

// GetStudents responds with the list of all users as JSON.
func GetStudents(myDb *db.MyDatabase) gin.HandlerFunc {
	return func(c *gin.Context) {

		students, err := RetrieveStudents(myDb)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("content-type", "application/json")
		c.JSON(http.StatusOK, students)
		fmt.Printf("Successfully retrieved student list \n")
	}
}

// AddStudent adds an user from JSON received in the request body.
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

		// // Add the new student to the slice.
		// students = append(students, newStudent)
		// c.Header("content-type", "application/json")
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

		var updatedStudent Student

		// Call BindJSON to bind the received JSON to updatedStudent
		if err := c.BindJSON(&updatedStudent); err != nil {
			return
		}

		if err := UpdatedStudentDB(myDb, integerID, updatedStudent); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		updatedStudent.ID = integerID
		c.Header("content-type", "application/json")
		c.JSON(http.StatusOK, updatedStudent)
		fmt.Printf("Successfully updated student id: %v \n", integerID)
	}
}

// // SoftDeleteStudent changes the Active status of the student to false, rather than permanently deleting.
// func SoftDeleteStudent(c *gin.Context) {
// 	id := c.Param("id")

// 	// myDb.SoftDeleteStudent

// 	for i, student := range students {
// 		if student.ID == id {
// 			students[i].Active = false
// 			c.Header("content-type", "application/json")
// 			c.JSON(http.StatusOK, students[i])
// 			fmt.Printf("Successfully deleted student: %v", students[i])
// 			return
// 		}
// 	}
// }
