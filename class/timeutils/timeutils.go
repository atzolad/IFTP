package timeutils

import (
	"fmt"
	"time"
)

func FindMonthlyDates(starting_date time.Time) error {
	// Set the starting date
	//Check all the days after that- increment by 7 if the month is the same, add it , if not do not

	// // set the starting date (in any way you wish) - replace with starting_date- if input is a string
	// start, err := time.Parse("2006-1-2", starting_date)
	// if err != nil {
	// 	return fmt.Errorf("error parsing starting date for monthly day / date aggregation: %v", err)
	// }
	// handle error

	// set d to starting date and keep adding 7 days to it as long as month doesn't change
	for d := starting_date; d.Month() == starting_date.Month(); d = d.AddDate(0, 0, 7) {
		date := d.String()
		fmt.Println(date)
	}

	return nil
}

func FirstDay(weekday time.Weekday, year int, month time.Month) int {
	t := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	return (8-int(t.Weekday()))%7 + int(weekday)
}
