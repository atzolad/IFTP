package class

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type class struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Teacher string `json:"teacher"`
	Email   string `json:"email"`
	Day     string `json:"classes"`
	Active  bool   `json:"active"`
}

var classes = []class{
	{ID: "1", Name: "Mondays with Matthew Moore", Teacher: "Matthew Moore", Email: "matt@gmail.com", Day: "Monday", Active: true},
	{ID: "2", Name: "Mondays with Ava Abdoulah", Teacher: "Ava Abdoulah", Email: "ava@gmail.com", Day: "Monday", Active: true},
	{ID: "3", Name: "Tuesdays with Matthew Moore", Teacher: "Matthew Moore", Email: "matt@gmail.com", Day: "Tuesday", Active: true},
	{ID: "4", Name: "Tuesdays with Liam Clancy", Teacher: "Liam Clancy", Email: "liam@gmail.com", Day: "Tuesday", Active: true},
}

// GetClassesresponds with the list of all classes as JSON.
func GetClasses(c *gin.Context) {
	c.Header("content-type", "application/json")
	c.JSON(http.StatusOK, classes)
	fmt.Printf("Successfully retrieved classes list with %v classes", len(classes))
}

// AddClass adds an class from JSON received in the request body.
func AddClass(c *gin.Context) {
	var newClass class

	// Call BindJSON to bind the received JSON to newClass
	if err := c.BindJSON(&newClass); err != nil {
		return
	}

	// myDb.AddClass(newClass)

	// Add the new class to the slice.
	classes = append(classes, newClass)
	c.Header("content-type", "application/json")
	c.JSON(http.StatusCreated, newClass)
	fmt.Printf("Successfully created new class: %v", newClass)
}

// Update class updates the class details based on the JSON received in the request body.
func UpdateClass(c *gin.Context) {
	id := c.Param("id")
	var updatedClass class

	// Call BindJSON to bind the received JSON to updatedclass
	if err := c.BindJSON(&updatedClass); err != nil {
		return
	}

	// myDb.UpdateClass

	for i, class := range classes {
		if class.ID == id {
			originalClass := classes[i]
			if class.Name != "" {
				classes[i].Name = updatedClass.Name
			}
			if class.Teacher != "" {
				classes[i].Teacher = updatedClass.Teacher
			}
			if class.Day != "" {
				classes[i].Day = updatedClass.Day
			}

			classes[i].Active = updatedClass.Active
			c.Header("content-type", "application/json")
			c.JSON(http.StatusOK, classes[i])
			fmt.Printf("Successfully updated class: %v with %v", originalClass, updatedClass)
			return
		}
	}
}

// SoftDeleteclass changes the Active status of the class to false, rather than permanently deleting.
func SoftDeleteclass(c *gin.Context) {
	id := c.Param("id")

	// myDb.SoftDeleteclass

	for i, class := range classes {
		if class.ID == id {
			classes[i].Active = false
			c.Header("content-type", "application/json")
			c.JSON(http.StatusOK, classes[i])
			fmt.Printf("Successfully deleted class: %v", classes[i])
			return
		}
	}
}
