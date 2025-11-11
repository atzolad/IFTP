package student

import (
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
}

// SoftDeleteStudent changes the Active status of the student to false, rather than permanently deleting.
func SoftDeleteStudent(c *gin.Context) {
	id := c.Param("id")

	for i, student := range students {
		if student.ID == id {
			students[i].Active = false
			c.Header("content-type", "application/json")
			c.JSON(http.StatusOK, students[i])
			return
		}
	}
}
