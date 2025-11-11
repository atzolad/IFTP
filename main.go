package main

import (
	"IFTP/user"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// User endpoints
	router.GET("/users", user.GetUsers)
	router.POST("/users", user.AddUser)
	router.DELETE("/users", user.DeleteUser)
	//router.PATCH("/users", user.UpdateUser)

	// Class endpoints
	// router.GET("/classes", getClasses)
	// router.POST("/classes", addClass)

	router.Run()
}
