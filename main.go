package main

import (
	"IFTP/student"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// User endpoints
	router.GET("/students", student.GetStudents)
	router.POST("/students", student.AddStudent)
	router.DELETE("/students/:id", student.SoftDeleteStudent)
	//router.PATCH("/users", student.UpdateUser)

	// Class endpoints
	// router.GET("/classes", getClasses)
	// router.POST("/classes", addClass)
	// router.

	router.Run()
}
