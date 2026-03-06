package students

import (
	"IFTP/db"
	"IFTP/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Student struct {
	ID              int      `db:"id" json:"id"`
	Name            string   `db:"name" json:"name"`
	Email           string   `db:"email" json:"email"`
	Active          bool     `db:"active" json:"active"`
	EnrolledClasses []string `db:"enrolled_classes" json:"enrolledClasses"`
}

func (s *Student) Sanitize() {
	s.Name = strings.TrimSpace(s.Name)
	s.Email = strings.ToLower(strings.TrimSpace(s.Email))
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

func AddStudent(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		var newStudent Student

		if err := json.NewDecoder(r.Body).Decode(&newStudent); err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: fmt.Sprintf("Error decoding student info: %v", err),
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Error decoding student info: %v", err)
			return
		}

		newStudent.Sanitize()

		if err := dbAddStudent(ctx, myDb, &newStudent); err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: fmt.Sprintf("Error adding student to db: %v", err),
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Error adding student to db: %v", err)
			return
		}

		utils.WriteJSONResponse(w, http.StatusOK, "Successfully added student to db")
		myDb.Logger.Printf("Successfully created new student: %v", newStudent.Name)
	}
}

// func UpdateStudent(myDb *db.MyDatabase) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		ctx := r.Context()

// 		id := r.PathValue("student_id")

// 		integerID, err := strconv.Atoi(id)
// 		if err != nil {
// 			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
// 				Status:  "error",
// 				Message: "Invalid student id- must be an integer:",
// 				Code:    http.StatusInternalServerError,
// 			})
// 			myDb.Logger.Printf("Error converting student id to int: %v", err)
// 			return
// 		}

// 		var updateStudent Student

// 		if err := json.NewDecoder(r.Body).Decode(&updateStudent); err != nil {
// 			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
// 				Status:  "error",
// 				Message: fmt.Sprintf("Error decoding student info: %v", err),
// 				Code:    http.StatusInternalServerError,
// 			})
// 			myDb.Logger.Printf("Error decoding new student info from JSON: %v", err)
// 			return
// 		}

// 		updateStudent.Sanitize()
// 		updateStudent.ID = integerID

// 		returnedStudent, err := UpdateStudentDB(ctx, myDb, integerID, &updateStudent)

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
