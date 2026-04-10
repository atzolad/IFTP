package class

import (
	"IFTP/db"
	"IFTP/timeutils"
	"IFTP/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

	var monthStr string
	if c.Month != nil {
		monthStr = c.Month.Format("2006-01-02")
	} else {
		monthStr = ""
	}

	return json.Marshal(&struct {
		Month        string   `json:"month"`
		SessionDates []string `json:"session_dates"`
		EndDate      string   `json:"endDate"`
		Time         string   `json:"time"`
		Alias
	}{
		Month:        monthStr,
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

type ActiveClass struct {
	ID        int    `db:"id"`
	DayOfWeek string `db:"day_of_week"`
}

type MonthlyClassScheduleApproval struct {
	ID           string      `db:"id" json:"id"`
	ClassID      int         `db:"class_id" json:"class_id"`
	ClassName    string      `db:"class_name" json:"class_name"`
	Teacher      string      `db:"teacher" json:"teacher"`
	DayOfWeek    DayOfWeek   `db:"day_of_week" json:"day_of_week"`
	Time         string      `db:"time" json:"time"`
	Capacity     int         `db:"capacity" json:"capacity"`
	Month        *time.Time  `db:"month" json:"month"`
	Status       string      `db:"status" json:"status"`
	PendingDates []time.Time `db:"pending_dates" json:"pending_dates"`
}

type ScheduleApprovalInput struct {
	Status string      `json:"status"` // Approved or Rejected
	Dates  []time.Time `json:"dates"`  // Returned dates after calendar edits
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
		utils.WriteJSONResponse(w, http.StatusOK, classes)
	}
}

func TriggerScheduleApprovals(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if err := GenerateScheduleApprovals(ctx, myDb); err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: "Error generating schedule approvals",
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Error generating schedule approvals: %v", err)
			return
		}
		utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
			Status:  "success",
			Message: "Successfully generated schedule approvals",
			Code:    http.StatusOK,
		})
	}
}

func GenerateScheduleApprovals(ctx context.Context, myDb *db.MyDatabase) error {
	// Having this automatically calculate the next month from now so cron job doesn't need to handle times. Can change this later if needed.
	now := time.Now()
	month := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())

	myDb.Logger.Printf("Generating schedule approvals for month: %v", month)
	rows, err := myDb.Pool.Query(ctx, `
	SELECT id, day_of_week FROM classes WHERE active = true`)
	if err != nil {
		return err
	}
	defer rows.Close()

	classes, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[ActiveClass])
	if err != nil {
		return fmt.Errorf("Error collecting active classes: %v", err)
	}

	myDb.Logger.Printf("Found %v active classes to process", len(classes))

	for _, class := range classes {
		if err := generateClassApproval(ctx, myDb, class.ID, class.DayOfWeek, month); err != nil {
			return fmt.Errorf("Error generating approval for class %v: %w", class.ID, err)
		}
	}
	myDb.Logger.Printf("Successfully generated schedule approvals for %v classes", len(classes))
	return nil
}

func generateClassApproval(ctx context.Context, myDb *db.MyDatabase, classId int, dayOfWeek string, month time.Time) error {
	tx, err := myDb.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Check if this approval already exists
	var exists bool

	err = tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM schedule_approvals
			WHERE class_id = $1 and month = $2)`,
		classId, month).Scan(&exists)

	if err != nil {
		return err
	}
	if exists {
		myDb.Logger.Printf("Approval already exists for class %v month %v, skipping", classId, month)
		return nil // already exists- no need to process
	}

	weekday, err := timeutils.ParseWeekday(dayOfWeek)
	if err != nil {
		return err
	}

	dates := timeutils.CreateDatesMap([]time.Weekday{weekday}, month.Year(), month.Month())
	pendingDates := dates[weekday]

	approvalId, err := dbInsertScheduleApproval(ctx, tx, classId, month)
	if err != nil {
		return err
	}

	err = dbInsertPendingDates(ctx, tx, approvalId, pendingDates)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	myDb.Logger.Printf("Successfully created approval %v for class %v", approvalId, classId)
	return nil

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
			myDb.Logger.Printf("Erorr adding class to db: %v", err)
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error adding Class to DB",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		hasSessionDates := false

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
				myDb.Logger.Printf("Error batch scheduling session dates: %v", err)
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
			myDb.Logger.Printf("Successfully created new class with session dates: %v", newClass)
		} else {
			utils.WriteJSONResponse(w, http.StatusOK, "Successfully created new class")
			myDb.Logger.Printf("Successfully created new class: %v", newClass)
		}
	}
}

// Update Class updates the class details based on the JSON received in the request body.
func UpdateClass(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		id := r.PathValue("class_id")
		integerID, err := strconv.Atoi(id)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: fmt.Sprintf("Invalid Class id- must be an integer: %v", err),
				Code:    http.StatusBadRequest,
			})
			myDb.Logger.Printf("Invalid Class id- must be an integer: %v", err)
			return
		}

		var updateRequest Class

		if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
				Status:  "error",
				Message: fmt.Sprintf("Error Decoding Request: %v", err),
				Code:    http.StatusBadRequest,
			})
			myDb.Logger.Printf("Error Decoding Request: %v", err)
			return
		}

		updateRequest.ID = integerID

		// Validate that the time provided is in the correct format
		if updateRequest.Time != "" {
			parsedTime, err := time.Parse("15:04", updateRequest.Time)
			if err != nil {
				parsedTime, err = time.Parse("15:04:05", updateRequest.Time)
			}
			if err != nil {
				utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
					Status:  "error",
					Message: "Invalid time format, expected HH:MM or HH:MM:SS",
					Code:    http.StatusBadRequest,
				})
				myDb.Logger.Printf("Error updating class- invalid time format: %v", err)
				return
			}
			updateRequest.Time = parsedTime.Format("15:04:05")
		}

		myDb.Logger.Printf("updateRequest Request- Name: %v, Teacher: %v, Day: %v, Time: %v, Description: %v, Month: %v, Capacity: %v, SessionDates: %v",
			updateRequest.Name, updateRequest.Teacher, updateRequest.DayOfWeek, updateRequest.Time, updateRequest.Description, updateRequest.Month, updateRequest.Capacity, updateRequest.SessionDates)

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

		returnedClass, err := dbUpdateClass(ctx, tx, integerID, &updateRequest)
		if returnedClass == nil {

		}
		if err != nil {
			if errors.Is(err, utils.ErrNoFieldsToUpdate) {
				utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ResponseData{
					Status:  "error",
					Message: "Error updating class in database- not fields provided to update",
					Code:    http.StatusBadRequest,
				})

			} else {

				utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
					Status:  "error",
					Message: "Error updating Class in DB",
					Code:    http.StatusInternalServerError,
				})

			}

			myDb.Logger.Printf("Error updating class in database: %v", err)
			return
		}

		hasSessionDates := false

		if updateRequest.SessionDates != nil {
			hasSessionDates = true
			err := dbDeleteFromClassSchedule(ctx, tx, updateRequest.ID)
			if err != nil {
				utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
					Status:  "error",
					Message: "Error removing old session dates from db",
					Code:    http.StatusInternalServerError,
				})
				myDb.Logger.Printf("Error removing old session dates from db: %v", err)
				return
			}

			if len(updateRequest.SessionDates) > 0 {
				batch := &pgx.Batch{}

				for _, sessionDate := range updateRequest.SessionDates {
					batch.Queue(
						`INSERT INTO class_schedule (class_id, session_date, month, status)
					VALUES ($1, $2, $3, $4)`,
						updateRequest.ID, sessionDate, updateRequest.Month, "scheduled")
				}

				br := tx.SendBatch(ctx, batch)

				if err := br.Close(); err != nil {
					myDb.Logger.Printf("Error batch scheduling session dates: %v", err)
					utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
						Status:  "error",
						Message: "Error batch scheduling session dates",
						Code:    http.StatusInternalServerError,
					})
					return
				}
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
			utils.WriteJSONResponse(w, http.StatusOK, updateRequest)
			myDb.Logger.Printf("Successfully updated class with session dates: %v", updateRequest)
		} else {
			utils.WriteJSONResponse(w, http.StatusOK, updateRequest)
			myDb.Logger.Printf("Successfully updated class: %v", updateRequest)
		}
	}
}

func GetPendingScheduleApprovals(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		approvals, err := dbGetPendingScheduleApprovals(ctx, myDb)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: fmt.Sprintf("Error fetching pending schedule approvals from db: %v", err),
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Print("Error fetching pending schedule approvals from db")
			return
		}
		utils.WriteJSONResponse(w, http.StatusOK, approvals)
	}
}

func ConfirmScheduleApproval(myDb *db.MyDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		approvalId := r.PathValue("approval_id")

		var input ScheduleApprovalInput

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error decoding confirmed schedule dates for approval",
				Code:    http.StatusBadRequest,
			})
			myDb.Logger.Printf("Error decoding confirmed schedule dates for approval: %v", err)
			return
		}

		if input.Status != "Approved" && input.Status != "Rejected" {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Status must be Approved or Rejected",
				Code:    http.StatusBadRequest,
			})
			return
		}

		tx, err := myDb.Pool.Begin(ctx)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error confirming schedule dates for approval",
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Error confirming schedule dates for approval: %v", err)
			return
		}

		defer tx.Rollback(ctx)

		classId, month, err := dbUpdateScheduleApprovalStatus(ctx, tx, approvalId, input.Status)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error updating schedule approval status",
				Code:    http.StatusInternalServerError,
			})
			myDb.Logger.Printf("Error updating schedule status: %v", err)
			return
		}

		if input.Status == "Approved" {
			if err := dbInsertClassScheduleRows(ctx, tx, classId, month, input.Dates); err != nil {
				utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
					Status:  "error",
					Message: "Error inserting into class schedule",
					Code:    http.StatusInternalServerError,
				})
				myDb.Logger.Printf("Error inserting into class schedule: %v", err)
				return
			}
		}

		if err := tx.Commit(ctx); err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.ResponseData{
				Status:  "error",
				Message: "Error committing transaction",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		utils.WriteJSONResponse(w, http.StatusOK, utils.ResponseData{
			Status:  "success",
			Message: fmt.Sprintf("Schedule approval %v successfully", input.Status),
			Code:    http.StatusOK,
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
