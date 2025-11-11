package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type user struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = []user{
	{ID: "1", Name: "Alex Zolad", Email: "atzolad@gmail.com"},
	{ID: "2", Name: "Megan Chang", Email: "meganchang10@gmail.com"},
	{ID: "3", Name: "Sarah Vaughan", Email: "Sarahvaughan@gmail.com"},
}

// getUsers responds with the list of all users as JSON.
func GetUsers(c *gin.Context) {
	c.Header("content-type", "application/json")
	c.JSON(http.StatusOK, users)
}

// postUsers adds an user from JSON received in the request body.
func AddUser(c *gin.Context) {
	var newUser user

	// Call BindJSON to bind the received JSON to
	// newUser.
	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	myDb.AddUser(newUser)

	// Add the new user to the slice.
	users = append(users, newUser)
	c.Header("content-type", "application/json")
	c.JSON(http.StatusCreated, newUser)
}

// postUsers adds an user from JSON received in the request body.
func DeleteUser(c *gin.Context) {
	var delUser user

	// Call BindJSON to bind the received JSON to
	// newUser.
	if err := c.BindJSON(&delUser); err != nil {
		return
	}

	if err := db.DeleteUser(delUser) {
		err!
	}

	// Add the new user to the slice.
	users = append(users, delUser)
	c.IndentedJSON(http.StatusCreated, delUser)
}

// postUsers adds an user from JSON received in the request body.
func UpdateUser(c *gin.Context) {
	var newUser user

	// Call BindJSON to bind the received JSON to
	// newUser.
	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	// Add the new user to the slice.
	users = append(users, newUser)
	c.IndentedJSON(http.StatusCreated, newUser)
}
