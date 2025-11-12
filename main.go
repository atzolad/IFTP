package main

import (
	"IFTP/class"
	"IFTP/db"
	"IFTP/student"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {

	// Load env vars from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: no .env file found")
	}

	router := gin.Default()
	connStr := os.Getenv("CONN_STR")
	fmt.Println("Connecting with:", connStr)
	// Initialise the connection pool.
	sqldb, err := sql.Open("pgx", connStr)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Succesfully connected to database \n")
	defer sqldb.Close()

	// Test the connection
	if err := sqldb.Ping(); err != nil {
		log.Fatal(err)
	}

	// Create an instance of myDatabase containing the connection pool.
	myDb := &db.MyDatabase{Db: sqldb}

	// User endpoints- make sure to pass the database instance to each function.
	router.GET("/students", student.GetStudents(myDb))
	router.POST("/students", student.AddStudent(myDb))
	router.PATCH("/students/:id", student.UpdateStudent(myDb))
	// router.DELETE("/students/:id", student.SoftDeleteStudent(myDb))

	// Class endpoints
	router.GET("/classes", class.GetClasses)
	router.POST("/classes", class.AddClass)
	router.PATCH("/classes/:id", class.UpdateClass)
	router.DELETE("/classes/:id", class.SoftDeleteclass)

	router.Run()
}
