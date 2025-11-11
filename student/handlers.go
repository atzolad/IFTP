package student

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type student struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Classes string `json:"classes"`
	Paid    bool   `json:"paid"`
	Active  bool   `json:"active"`
}

var students = []student{
	{ID: "1", Name: "Alex Zolad", Email: "atzolad@gmail.com", Classes: "", Paid: false, Active: true},
	{ID: "2", Name: "Megan Chang", Email: "meganchang10@gmail.com", Classes: "", Paid: false, Active: true},
	{ID: "3", Name: "Sarah Vaughan", Email: "Sarahvaughan@gmail.com", Classes: "", Paid: false, Active: true},
}

// GetStudents responds with the list of all users as JSON.
func GetStudents(c *gin.Context) {
	c.Header("content-type", "application/json")
	c.JSON(http.StatusOK, students)
	fmt.Printf("Successfully retrieved students list with %v students", len(students))
}

// AddStudent adds an user from JSON received in the request body.
func AddStudent(c *gin.Context) {
	var newStudent student

	// Call BindJSON to bind the received JSON to
	// newStudent.
	if err := c.BindJSON(&newStudent); err != nil {
		return
	}

	// myDb.AddStudent(newStudent)

	// Add the new student to the slice.
	students = append(students, newStudent)
	c.Header("content-type", "application/json")
	c.JSON(http.StatusCreated, newStudent)
	fmt.Printf("Successfully created new student: %v", newStudent)
}

// Update Student updates the student details based on the JSON received in the request body.
func UpdateStudent(c *gin.Context) {
	id := c.Param("id")
	var updatedStudent student

	// Call BindJSON to bind the received JSON to updatedStudent
	if err := c.BindJSON(&updatedStudent); err != nil {
		return
	}

	// myDb.UpdateStudent

	for i, student := range students {
		if student.ID == id {
			originalStudent := students[i]
			if student.Name != "" {
				students[i].Name = updatedStudent.Name
			}
			if student.Email != "" {
				students[i].Email = updatedStudent.Email
			}
			if student.Classes != "" {
				students[i].Classes = updatedStudent.Classes
			}

			students[i].Paid = updatedStudent.Paid
			students[i].Active = updatedStudent.Active
			c.Header("content-type", "application/json")
			c.JSON(http.StatusOK, students[i])
			fmt.Printf("Successfully updated student: %v with %v", originalStudent, updatedStudent)
			return
		}
	}
}

// SoftDeleteStudent changes the Active status of the student to false, rather than permanently deleting.
func SoftDeleteStudent(c *gin.Context) {
	id := c.Param("id")

	// myDb.SoftDeleteStudent

	for i, student := range students {
		if student.ID == id {
			students[i].Active = false
			c.Header("content-type", "application/json")
			c.JSON(http.StatusOK, students[i])
			fmt.Printf("Successfully deleted student: %v", students[i])
			return
		}
	}
}
