package students

import (
	"IFTP/db"
	"IFTP/utils"
	"fmt"
	"net/http"
)

type Student struct {
	ID              int      `db:"id" json:"id"`
	Name            string   `db:"name" json:"name"`
	Email           string   `db:"email" json:"email"`
	Active          bool     `db:"active" json:"active"`
	EnrolledClasses []string `db:"enrolled_classes" json:"enrolledClasses"`
}

func GetStudents(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		students, err := dbRetrieveStudents(ctx, myDb)
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

		ctx := r.Context()
		students, err := dbGetStudentsWithEnrollment(ctx, myDb)
		if err != nil {
			fmt.Printf("Error retrieving students from DB: %v", err)
			utils.WriteJSONResponse(w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSONResponse(w, http.StatusOK, students)
		fmt.Printf("Successfully retrieved student list \n")
	}
}

// // AddStudent adds a user from JSON received in the request body.
// func AddStudent(myDb *db.MyDatabase) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		ctx := r.Context()

// 		var newStudent Student

// 		if err := json.NewDecoder(r.Body).Decode(&newStudent); err != nil {
// 			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
// 				Status:  "error",
// 				Message: fmt.Sprintf("Error decoding student info: %v", err),
// 				Code:    http.StatusInternalServerError,
// 			})
// 		}

// 		if err := InsertStudent(ctx, myDb, &newStudent); err != nil {
// 			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
// 				Status:  "error",
// 				Message: fmt.Sprintf("Error adding student to db: %v", err),
// 				Code:    http.StatusInternalServerError,
// 			})
// 			return
// 		}

// 		utils.WriteJSONResponse(w, http.StatusOK, "Successfully added student to db")
// 		fmt.Printf("Successfully created new student: %v", newStudent)
// 	}
// }

// // Update Student updates the student details based on the JSON received in the request body.
// func UpdateStudent(myDb *db.MyDatabase) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		id := c.Param("id")
// 		integerID, err := strconv.Atoi(id)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		}

// 		var updateStudent Student

// 		// Call BindJSON to bind the received JSON to updatedStudent
// 		if err := c.BindJSON(&updateStudent); err != nil {
// 			return
// 		}

// 		returnedStudent, err := UpdateStudentDB(myDb, integerID, &updateStudent)

// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}

// 		c.Header("content-type", "application/json")
// 		c.JSON(http.StatusOK, returnedStudent)
// 		fmt.Printf("Successfully updated student: %v \n", returnedStudent.Name)
// 	}
// }

// // SoftDeleteStudent changes the Active status of the student to false, rather than permanently deleting.
// func SoftDeleteStudent(myDb *db.MyDatabase) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		id := c.Param("id")

// 		integerID, err := strconv.Atoi(id)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student id"})
// 		}

// 		deletedStudent, err := SoftDeleteStudentDB(myDb, integerID)
// 		if err != nil {
// 			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("student %v deleted successfully", deletedStudent)})
// 		fmt.Printf("Successfully soft deleted student %v with id: %v \n", deletedStudent, integerID)
// 	}
// }
