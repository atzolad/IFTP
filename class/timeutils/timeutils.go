package timeutils

import (
	"fmt"
	"time"
)

func FindMonthlyDates(starting_date time.Time) []time.Time {
	datesList := make([]time.Time, 0)
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
		datesList = append(datesList, d)
		// date := d.String()
	}

	return datesList
}

func FirstDay(weekday time.Weekday, year int, month time.Month) int {
	location, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return 0
	}

	t := time.Date(year, month, 1, 0, 0, 0, 0, location)
	return (8-int(t.Weekday()))%7 + int(weekday)
}

func CreateDatesMap(classDays []time.Weekday, year int, month time.Month) (datesMap map[time.Weekday][]time.Time) {
	datesMap = make(map[time.Weekday][]time.Time)

	location, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return nil
	}

	for _, day := range classDays {

		_, ok := datesMap[day]
		if ok {
			continue
		} else {
			firstDayOfMonth := FirstDay(day, year, month)
			firstDate := time.Date(year, month, firstDayOfMonth, 0, 0, 0, 0, location)
			monthlyDates := FindMonthlyDates(firstDate)
			datesMap[day] = monthlyDates
		}
	}
	return datesMap
}

//createDatesMap:
// go through each day of the week in the list
//For Monday- if Monday does not exists in the map, add it and add the first date.

// Find all of the days that have classes.
// For all the days that have classes- make a map with the weekday as the key and a slice of those days of the month as the values.

// if datesMap['friday'].exists()
// else {
// 	friday = FirstDay(time.Friday, 2025, 11)
// 	datesMap['friday'] = friday
// }

// datesMap map[time.Weekday][]ints

// 'Friday': 3
// 'saturday':4
