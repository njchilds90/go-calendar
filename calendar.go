// Package calendar provides functions for calendar calculations and printing,
// similar to Python's calendar module.
package calendar

import (
	"fmt"
	"strings"
	"time"
)

// Constants for weekdays (matching time.Weekday)
const (
	Sunday    = 0
	Monday    = 1
	Tuesday   = 2
	Wednesday = 3
	Thursday  = 4
	Friday    = 5
	Saturday  = 6
)

// DayNames are full day names.
var DayNames = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

// DayAbbrs are abbreviated day names.
var DayAbbrs = []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

// MonthNames are full month names (index 0 unused).
var MonthNames = []string{"", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}

// MonthAbbrs are abbreviated month names (index 0 unused).
var MonthAbbrs = []string{"", "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

// firstWeekday is the starting day of the week (default Monday, like Python).
var firstWeekday = Monday

// SetFirstWeekday sets the first day of the week (0=Sunday, 1=Monday, etc.).
func SetFirstWeekday(wd int) {
	if wd < 0 || wd > 6 {
		panic("invalid weekday")
	}
	firstWeekday = wd
}

// FirstWeekday returns the current first day of the week.
func FirstWeekday() int {
	return firstWeekday
}

// IsLeap returns true if year is a leap year.
func IsLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// LeapDays returns the number of leap years between y1 and y2 (exclusive y2).
func LeapDays(y1, y2 int) int {
	f := func(y int) int {
		return y/4 - y/100 + y/400
	}
	return f(y2) - f(y1)
}

// Weekday returns the weekday (0=Sunday, ..., 6=Saturday) for the given date.
func Weekday(year, month, day int) int {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return int(t.Weekday())
}

// MonthRange returns the weekday of the first day and number of days in the month.
func MonthRange(year, month int) (weekday, numDays int) {
	if month < 1 || month > 12 {
		panic("invalid month")
	}
	numDays = 31
	switch month {
	case 4, 6, 9, 11:
		numDays = 30
	case 2:
		if IsLeap(year) {
			numDays = 29
		} else {
			numDays = 28
		}
	}
	weekday = Weekday(year, month, 1)
	return
}

// MonthCalendar returns a matrix representing a month's calendar (weeks as rows, days as columns, 0 for padding).
func MonthCalendar(year, month int) [][]int {
	wd, days := MonthRange(year, month)
	cal := make([][]int, 6) // Max 6 weeks
	for i := range cal {
		cal[i] = make([]int, 7)
	}
	day := 1
	shift := (wd - firstWeekday + 7) % 7
	for w := 0; w < 6; w++ {
		for d := 0; d < 7; d++ {
			idx := w*7 + d
			if idx < shift || day > days {
				cal[w][d] = 0
			} else {
				cal[w][d] = day
				day++
			}
		}
		if day > days {
			break
		}
	}
	// Trim empty weeks
	for i := len(cal) - 1; i >= 0; i-- {
		empty := true
		for _, d := range cal[i] {
			if d != 0 {
				empty = false
				break
			}
		}
		if empty {
			cal = cal[:i]
		} else {
			break
		}
	}
	return cal
}

// WeekHeader returns a string of weekday abbreviations (width per name).
func WeekHeader(width int) string {
	var sb strings.Builder
	for i := 0; i < 7; i++ {
		day := (i + firstWeekday) % 7
		abbr := DayAbbrs[day]
		if len(abbr) > width {
			abbr = abbr[:width]
		}
		sb.WriteString(fmt.Sprintf("%-*s ", width, abbr))
	}
	return sb.String()[:sb.Len()-1] // Trim trailing space
}

// FormatMonth returns a multi-line string for the month calendar.
func FormatMonth(year, month int, w, l int) string {
	if w < 1 {
		w = 3
	}
	if l < 1 {
		l = 0
	}
	var sb strings.Builder
	header := fmt.Sprintf("%s %d", MonthNames[month], year)
	sb.WriteString(fmt.Sprintf("%*s\n", (7*(w+1)-1+len(header))/2, header))
	sb.WriteString(WeekHeader(w) + "\n")
	cal := MonthCalendar(year, month)
	for _, week := range cal {
		for _, day := range week {
			if day == 0 {
				sb.WriteString(fmt.Sprintf("%*s ", w, ""))
			} else {
				sb.WriteString(fmt.Sprintf("%*d ", w, day))
			}
		}
		sb.WriteString("\n")
		if l > 0 {
			sb.WriteString(strings.Repeat("\n", l))
		}
	}
	return sb.String()
}

// PrMonth prints the month calendar to stdout.
func PrMonth(year, month int, w, l int) {
	fmt.Print(FormatMonth(year, month, w, l))
}
