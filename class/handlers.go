package class

import (
	"IFTP/db"
	"IFTP/utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Class struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	Teacher       string   `json:"teacher"`
	DayOfWeek     string   `json:"day_of_week"`
	Time          string   `json:"time"`
	Description   string   `json:"description"`
	Month         string   `json:"month"`
	Capacity      string   `json:"capacity"`
	SessionDates  []string `json:"session_dates"`
	EnrolledCount int      `json:"enrolledCount"`
}

type CalendarEventsResponse struct {
	ScheduledClasses []Class `json:"scheduledClasses"`
}

// GetClasses responds with the list of all classes as JSON.
// Nit: ListClasses
func ListClasses(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		classes, err := dbListClasses(myDb)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: fmt.Sprintf("Error fetching classes from db: %v", err),
				Code:    http.StatusInternalServerError,
			})
			return
		}

		utils.WriteJSONResponse(w, http.StatusOK, classes)
		fmt.Printf("Successfully retrieved class list \n")
	}
}

func ListClassesByMonth(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// month := c.Param("month")
		month := r.FormValue("month")
		studentId := strings.TrimSpace(r.PathValue("student_id"))
		var studentIntegerId *int
		fmt.Println("MONTH:")
		fmt.Println(month)
		fmt.Println("STUDENT ID:")
		fmt.Println(studentId)

		if studentId != "" {
			val, err := strconv.Atoi(studentId)
			if err != nil {
				utils.WriteJSONResponse(w, http.StatusBadRequest, err)
				return
			}
			studentIntegerId = &val
		}

		classes, err := dbListClassesByMonth(myDb, month, studentIntegerId)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: fmt.Sprintf("Error fetching classes from db for month %v: %v", err),
				Code:    http.StatusInternalServerError,
			})
			return
		}
		fmt.Println(classes)
		utils.WriteJSONResponse(w, http.StatusOK, classes)
		fmt.Printf("Successfully retrieved class list \n")
	}
}

// func ListStudentEnrolledClasses(myDb *db.MyDatabase) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// month := c.Param("month")
// 		month := r.FormValue("month")
// 		studentId := strings.TrimSpace(r.PathValue("student_id"))
// 		var studentIntegerId *int
// 		fmt.Println("MONTH:")
// 		fmt.Println(month)
// 		fmt.Println("STUDENT ID:")
// 		fmt.Println(studentId)

// 		if studentId != "" {
// 			val, err := strconv.Atoi(studentId)
// 			if err != nil {
// 				utils.WriteJSONResponse(w, http.StatusBadRequest, err)
// 				return
// 			}
// 			studentIntegerId = &val
// 		}

// 		classes, err := dbListStudentEnrolledClasses(myDb, month, studentIntegerId)
// 		if err != nil {
// 			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
// 				Status:  "error",
// 				Message: fmt.Sprintf("Error fetching classes from db for month: %v", err),
// 				Code:    http.StatusInternalServerError,
// 			})
// 			return
// 		}
// 		fmt.Println(classes)
// 		utils.WriteJSONResponse(w, http.StatusOK, classes)
// 		fmt.Printf("Successfully retrieved class list \n")
// 	}
// }

func GetCalendarEvents(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Getting calendar events")
		//TODO: If this endpoint gets really slow, add month
		// month := r.PathValue("month")

		classes, err := dbListClasses(myDb)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: fmt.Sprintf("Error fetching classes from db: %v", err),
				Code:    http.StatusInternalServerError,
			})
			return
		}

		utils.WriteJSONResponse(w, http.StatusOK, CalendarEventsResponse{
			ScheduledClasses: classes,
		})

	}
}

func GetCalendarEventsByStudent(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Getting calendar events by student")

		studentId := strings.TrimSpace(r.PathValue("student_id"))
		fmt.Printf("studentId: %v", studentId)
		var studentIntegerId *int

		if studentId != "" {
			val, err := strconv.Atoi(studentId)
			if err != nil {
				utils.WriteJSONResponse(w, http.StatusBadRequest, err)
				return
			}
			studentIntegerId = &val
		}

		//TODO: If this endpoint gets really slow, add month
		// month := r.PathValue("month")

		classes, err := dbListStudentEnrolledClasses(myDb, "", studentIntegerId)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: fmt.Sprintf("Error fetching classes from db %v", err),
				Code:    http.StatusInternalServerError,
			})
			return
		}

		utils.WriteJSONResponse(w, http.StatusOK, CalendarEventsResponse{
			ScheduledClasses: classes,
		})

	}
}

// func ApproveClassDates(myDb *db.MyDatabase) gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		classes, err := GetClassesDB(myDb)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}

// 		c.Header("content-type", "application/json")
// 		c.JSON(http.StatusOK, classes)
// 		fmt.Printf("Successfully retrieved class list \n")
// 	}
// }

// // GetClasses responds with the list of all classes as JSON.
// // Nit: ListClasses
// func GetClasses(myDb *db.MyDatabase) gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		classes, err := GetClassesDB(myDb)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}

// 		c.Header("content-type", "application/json")
// 		c.JSON(http.StatusOK, classes)
// 		fmt.Printf("Successfully retrieved class list \n")
// 	}
// }

// // AddClass adds a class from JSON received in the request body.
// // TODO: be consistent with your names CreateClass and dbCreateClass
// func CreateClass(myDb *db.MyDatabase) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var newClass Class

// 		// Call BindJSON to bind the received JSON to newClass.
// 		if err := c.BindJSON(&newClass); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		// Validate that the time provided is in the correct format
// 		parsedTime, err := time.Parse("15:04:05", newClass.Time)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format, expected HH:MM:SS"})
// 		}

// 		newClass.Time = parsedTime.Format("15:04:05")

// 		if err := DbCreateClass(myDb, &newClass); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}

// 		c.JSON(http.StatusCreated, newClass)
// 		fmt.Printf("Successfully created new class: %v", newClass)
// 	}
// }

// // Update Class updates the class details based on the JSON received in the request body.
// func UpdateClass(myDb *db.MyDatabase) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		id := c.Param("id")
// 		integerID, err := strconv.Atoi(id)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		}

// 		var updateClass Class

// 		// Call BindJSON to bind the received JSON to updatedStudent
// 		if err := c.BindJSON(&updateClass); err != nil {
// 			return
// 		}

// 		// Validate that the time provided is in the correct format
// 		if updateClass.Time != "" {
// 			parsedTime, err := time.Parse("15:04:05", updateClass.Time)
// 			if err != nil {
// 				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format, expected HH:MM:SS"})
// 			}
// 			updateClass.Time = parsedTime.Format("15:04:05")
// 		}

// 		returnedClass, err := UpdateClassDB(myDb, integerID, &updateClass)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}

// 		c.Header("content-type", "application/json")
// 		c.JSON(http.StatusOK, returnedClass)
// 		fmt.Printf("Successfully updated class: %v \n", returnedClass.Name)
// 	}
// }

// // SoftDeleteClass changes the Active status of the class to false, rather than permanently deleting.
// func SoftDeleteClass(myDb *db.MyDatabase) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		id := c.Param("id")

// 		integerID, err := strconv.Atoi(id)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class id"})
// 		}

// 		deletedClass, err := SoftDeleteClassDB(myDb, integerID)
// 		if err != nil {
// 			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("class %v deleted successfully", deletedClass)})
// 		fmt.Printf("Successfully soft deleted class %v with id: %v \n", deletedClass, integerID)
// 	}
// }
