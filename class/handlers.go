package class

import (
	"IFTP/db"
	"IFTP/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type classScheduleStatus string

// Define an enum for the class schedule status
// TODO need to define these ENUMS in Postgres db

const (
	PENDING   classScheduleStatus = "Pending"
	SCHEDULED classScheduleStatus = "Scheduled"
	CANCELLED classScheduleStatus = "Cancelled"
)

type DayOfWeek string

const (
	Monday    DayOfWeek = "Monday"
	Tuesday   DayOfWeek = "Tuesday"
	Wednesday DayOfWeek = "Wednesday"
	Thursday  DayOfWeek = "Thursday"
	Friday    DayOfWeek = "Friday"
	Saturday  DayOfWeek = "Saturday"
	Sunday    DayOfWeek = "Sunday"
)

type Class struct {
	ID            int         `db:"id" json:"id"`
	Name          string      `db:"name" json:"name"`
	Teacher       string      `db:"teacher" json:"teacher"`
	DayOfWeek     DayOfWeek   `db:"day_of_week" json:"day_of_week"`
	Time          string      `db:"time" json:"time"`
	Description   string      `db:"description" json:"description"`
	Month         *time.Time  `db:"month" json:"month"`
	Capacity      int         `db:"capacity" json:"capacity"`
	SessionDates  []time.Time `db:"session_dates" json:"session_dates"`
	EnrolledCount int         `db:"enrolled_count" json:"enrolledCount"`
	EndDate       time.Time   `db:"endDate" json:"endDate"`
}

func (c Class) MarshalJSON() ([]byte, error) {
	// Creates a new type with all of the fields of Class but none of the methods
	type Alias Class

	return json.Marshal(&struct {
		Month        string   `json:"month"`
		SessionDates []string `json:"session_dates"`
		EndDate      string   `json:"endDate"`
		Time         string   `json:"time"`
		Alias
	}{
		Month:        c.Month.Format("2006-01-02"),
		SessionDates: formatTimeSlice(c.SessionDates),
		EndDate:      c.EndDate.Format("2006-01-02"),
		Time:         formatTime(c.Time),
		Alias:        (Alias)(c),
	})
}

func formatTimeSlice(dates []time.Time) []string {
	if len(dates) == 0 {
		return []string{}
	}

	formatted := make([]string, len(dates))
	for i, d := range dates {
		formatted[i] = d.Format("2006-01-02")
	}
	return formatted
}

func formatTime(time string) string {
	if time == "" {
		return ""
	}
	return strings.Split(time, ".")[0]
}

type ClassSchedule struct {
	Id          int                 `json:"id"`
	ClassId     string              `json:"classId"`
	SessionDate time.Time           `json:"sessionDate"`
	Month       time.Time           `json:"month"`
	Status      classScheduleStatus `json:"status"`
}

type CalendarEventsResponse struct {
	ScheduledClasses []Class `json:"scheduledClasses"`
}

func ListClasses(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		classes, err := dbListClasses(ctx, myDb)
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

		ctx := r.Context()
		month := r.FormValue("month")
		studentId := strings.TrimSpace(r.PathValue("student_id"))
		var studentIntegerId *int
		// fmt.Println("MONTH:")
		// fmt.Println(month)
		// fmt.Println("STUDENT ID:")
		// fmt.Println(studentId)

		if studentId != "" {
			val, err := strconv.Atoi(studentId)
			if err != nil {
				utils.WriteJSONResponse(w, http.StatusBadRequest, err)
				return
			}
			studentIntegerId = &val
		}

		classes, err := dbListClassesByMonth(ctx, myDb, month, studentIntegerId)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: fmt.Sprintf("Error fetching classes from db for month : %v", err),
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

		ctx := r.Context()
		classes, err := dbListClasses(ctx, myDb)
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

		ctx := r.Context()
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

		classes, err := dbListStudentEnrolledClasses(ctx, myDb, "", studentIntegerId)
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

// TODO- Need to go back and handle the END date logic

func CreateClass(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		var newClass Class

		if err := json.NewDecoder(r.Body).Decode(&newClass); err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: fmt.Sprintf("Error Decoding Request: %v", err),
				Code:    http.StatusBadRequest,
			})
			return
		}

		fmt.Printf("New Class Request- Name: %v, Teacher: %v, Day: %v, Time: %v, Description: %v, Month: %v, Capacity: %v, SessionDates: %v",
			newClass.Name, newClass.Teacher, newClass.DayOfWeek, newClass.Time, newClass.Description, newClass.Month, newClass.Capacity, newClass.SessionDates)

		tx, err := myDb.Pool.Begin(ctx)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: fmt.Sprintf("Error Begining transcation: %v", err),
				Code:    http.StatusInternalServerError,
			})
			return
		}

		defer tx.Rollback(ctx)

		if err := dbCreateClass(ctx, tx, &newClass); err != nil {
			log.Printf("Erorr adding class to db: %v", err)
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error adding Class to DB",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		var hasSessionDates bool
		hasSessionDates = false

		if len(newClass.SessionDates) > 0 {
			hasSessionDates = true
			batch := &pgx.Batch{}

			for _, sessionDate := range newClass.SessionDates {
				batch.Queue(
					`INSERT INTO class_schedule (class_id, session_date, month, status)
					VALUES ($1, $2, $3, $4)`,
					newClass.ID, sessionDate, newClass.Month, "scheduled")
			}

			br := tx.SendBatch(ctx, batch)

			if err := br.Close(); err != nil {
				log.Printf("Error batch scheduling session dates: %v", err)
				utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
					Status:  "error",
					Message: "Error batch scheduling session dates",
					Code:    http.StatusInternalServerError,
				})
				return
			}

		}

		if err := tx.Commit(ctx); err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Failed to commit database changes",
				Code:    http.StatusInternalServerError,
			})
			return
		}
		if hasSessionDates {
			utils.WriteJSONResponse(w, http.StatusOK, "Successfully created new class and scheduled session dates")
			log.Printf("Successfully created new class with session dates: %v", newClass)
		} else {
			utils.WriteJSONResponse(w, http.StatusOK, "Successfully created new class")
			log.Printf("Successfully created new class: %v", newClass)
		}
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
