// Package calendar provides calendar utilities inspired by Python's calendar module.
// It supports leap year checks, weekday calculations, month calendar matrices,
// text-based printing, and configurable first weekday of the week.
//
// All functions use Go's time package internally and are timezone-agnostic (UTC-based).
//
// Example usage:
//
//	calendar.PrMonth(2025, 12, 3, 0) // prints December calendar
package calendar

import (
	"fmt"
	"strings"
	"time"
)

// Weekday constants (matches time.Weekday: 0=Sunday ... 6=Saturday)
const (
	Sunday    = 0
	Monday    = 1
	Tuesday   = 2
	Wednesday = 3
	Thursday  = 4
	Friday    = 5
	Saturday  = 6
)

// DayNames are full English weekday names (index matches constants).
var DayNames = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

// DayAbbrs are abbreviated weekday names.
var DayAbbrs = []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

// MonthNames are full English month names (index 0 unused).
var MonthNames = []string{"", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}

// MonthAbbrs are abbreviated month names (index 0 unused).
var MonthAbbrs = []string{"", "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

// firstWeekday defines the week start (default Monday).
var firstWeekday = Monday

// SetFirstWeekday sets the first day of the week (0=Sunday to 6=Saturday).
func SetFirstWeekday(wd int) {
	if wd < 0 || wd > 6 {
		panic("weekday must be 0-6")
	}
	firstWeekday = wd
}

// FirstWeekday returns the current first weekday.
func FirstWeekday() int {
	return firstWeekday
}

// IsLeap reports whether year is a leap year.
func IsLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// LeapDays counts leap years in [y1, y2) range.
func LeapDays(y1, y2 int) int {
	f := func(y int) int { return y/4 - y/100 + y/400 }
	return f(y2) - f(y1)
}

// Weekday returns the day of week for the date (0=Sunday ... 6=Saturday).
func Weekday(year, month, day int) int {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return int(t.Weekday())
}

// MonthRange returns (first_weekday, days_in_month) for the given year/month.
func MonthRange(year, month int) (int, int) {
	if month < 1 || month > 12 {
		panic("month must be 1-12")
	}
	days := 31
	switch month {
	case 4, 6, 9, 11:
		days = 30
	case 2:
		days = 28
		if IsLeap(year) {
			days = 29
		}
	}
	return Weekday(year, month, 1), days
}

// MonthCalendar returns a [][]int matrix for the month (rows=weeks, 0 = padding).
func MonthCalendar(year, month int) [][]int {
	wd, days := MonthRange(year, month)
	cal := make([][]int, 6)
	for i := range cal {
		cal[i] = make([]int, 7)
	}
	day := 1
	shift := (wd - firstWeekday + 7) % 7
	for w := 0; w < 6; w++ {
		for d := 0; d < 7; d++ {
			if w*7+d < shift || day > days {
				cal[w][d] = 0
			} else {
				cal[w][d] = day
				day++
			}
		}
		if day > days {
			return cal[:w+1]
		}
	}
	return cal
}

// weekHeader builds the weekday abbreviation line.
func weekHeader(width int) string {
	var sb strings.Builder
	for i := 0; i < 7; i++ {
		wd := (i + firstWeekday) % 7
		abbr := DayAbbrs[wd]
		if len(abbr) > width {
			abbr = abbr[:width]
		}
		fmt.Fprintf(&sb, "%*s ", width, abbr)
	}
	s := sb.String()
	return s[:len(s)-1]
}

// FormatMonth returns a formatted text calendar string for one month.
func FormatMonth(year, month, width, lines int) string {
	if width < 2 {
		width = 2
	}
	header := fmt.Sprintf("%s %d", MonthNames[month], year)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%*s\n", (7*(width+1)-1+len(header))/2, header))
	sb.WriteString(weekHeader(width) + "\n")
	cal := MonthCalendar(year, month)
	for _, week := range cal {
		for _, d := range week {
			if d == 0 {
				fmt.Fprintf(&sb, "%*s ", width, "")
			} else {
				fmt.Fprintf(&sb, "%*d ", width, d)
			}
		}
		sb.WriteString("\n")
		if lines > 0 {
			sb.WriteString(strings.Repeat("\n", lines))
		}
	}
	return sb.String()
}

// PrMonth prints the month calendar to stdout.
func PrMonth(year, month, width, lines int) {
	fmt.Print(FormatMonth(year, month, width, lines))
}
