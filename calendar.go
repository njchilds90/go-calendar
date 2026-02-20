// Package calendar provides calendar utilities inspired by Python's standard calendar module.
// It offers leap year detection, weekday and month calculations, calendar matrices as [][]int,
// configurable text and HTML output, year-at-a-glance views, and basic locale support for day/month names.
//
// All functions are timezone-agnostic (using UTC via time.Date).
//
// Key features:
//   - IsLeap, LeapDays
//   - Weekday, MonthRange, MonthCalendar
//   - FormatMonth / PrMonth (text)
//   - FormatMonthHTML / FormatYearHTML (HTML tables)
//   - FormatYear / PrYear (compact year view)
//   - SetFirstWeekday, SetLocale (custom names)
//
// Example:
//
//	import "github.com/njchilds90/go-calendar/calendar"
//
//	calendar.SetFirstWeekday(calendar.Monday)
//	calendar.PrMonth(2026, 2, 3, 0)
//	hc := calendar.NewHTMLCalendar(calendar.Monday)
//	htmlMonth := hc.FormatMonthHTML(2026, 2, true)
package calendar

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// Weekday represents a day of the week.
type Weekday int

// Weekday constants (Monday = 0 for Python parity).
const (
	Monday Weekday = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

// DayNames contains full weekday names starting with Monday.
var DayNames = []string{
	"Monday", "Tuesday", "Wednesday",
	"Thursday", "Friday", "Saturday", "Sunday",
}

// DayAbbr contains abbreviated weekday names starting with Monday.
var DayAbbr = []string{
	"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun",
}

// MonthNames contains full month names (index 1–12).
var MonthNames = []string{
	"", // padding for 1-based month index
	"January", "February", "March", "April", "May", "June",
	"July", "August", "September", "October", "November", "December",
}

// MonthAbbr contains abbreviated month names (index 1–12).
var MonthAbbr = []string{
	"",
	"Jan", "Feb", "Mar", "Apr", "May", "Jun",
	"Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
}

var firstWeekday = Monday

// SetFirstWeekday sets which weekday calendars start on.
// Default is Monday.
func SetFirstWeekday(day Weekday) {
	firstWeekday = day
}

// FirstWeekday returns the currently configured first weekday.
func FirstWeekday() Weekday {
	return firstWeekday
}

// IsLeap returns true if year is a leap year.
func IsLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// LeapDays returns the number of leap years in range [y1, y2).
func LeapDays(y1, y2 int) int {
	count := 0
	for y := y1; y < y2; y++ {
		if IsLeap(y) {
			count++
		}
	}
	return count
}

// MonthRange returns the first weekday and number of days in a month.
// Weekday follows configured first weekday offset.
func MonthRange(year, month int) (weekday int, days int) {
	t := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	weekday = (int(t.Weekday()) + 6) % 7 // convert Sunday=0 to Monday=0
	weekday = (weekday - int(firstWeekday) + 7) % 7
	days = time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
	return
}

// MonthCalendar returns a matrix representing a month's calendar.
// Days outside the month are zero.
func MonthCalendar(year, month int) [][]int {
	first, days := MonthRange(year, month)
	var weeks [][]int
	week := make([]int, 7)
	day := 1

	for i := 0; i < first; i++ {
		week[i] = 0
	}

	for i := first; i < 7; i++ {
		week[i] = day
		day++
	}
	weeks = append(weeks, week)

	for day <= days {
		week = make([]int, 7)
		for i := 0; i < 7 && day <= days; i++ {
			week[i] = day
			day++
		}
		weeks = append(weeks, week)
	}

	return weeks
}

// IterMonthDays returns sequential day numbers including padding zeros.
func IterMonthDays(year, month int) []int {
	var result []int
	for _, week := range MonthCalendar(year, month) {
		result = append(result, week...)
	}
	return result
}

// IterMonthDates returns sequential time.Time values for each calendar cell.
func IterMonthDates(year, month int) []time.Time {
	var result []time.Time
	for _, day := range IterMonthDays(year, month) {
		if day == 0 {
			result = append(result, time.Time{})
		} else {
			result = append(result,
				time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC))
		}
	}
	return result
}

// WeekHeader returns a formatted weekday header with given width.
func WeekHeader(width int) string {
	var parts []string
	for i := 0; i < 7; i++ {
		index := (int(firstWeekday) + i) % 7
		name := DayAbbr[index]
		if width < len(name) {
			name = name[:width]
		}
		parts = append(parts, fmt.Sprintf("%-*s", width, name))
	}
	return strings.Join(parts, " ")
}

// FormatMonth returns a formatted string for a single month.
func FormatMonth(year, month int) string {
	var sb strings.Builder
	title := fmt.Sprintf("%s %d", MonthNames[month], year)
	sb.WriteString(fmt.Sprintf("%^20s\n", title))
	sb.WriteString(WeekHeader(2) + "\n")

	for _, week := range MonthCalendar(year, month) {
		for _, day := range week {
			if day == 0 {
				sb.WriteString("   ")
			} else {
				sb.WriteString(fmt.Sprintf("%2d ", day))
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// FormatYear returns a formatted calendar for an entire year.
func FormatYear(year, monthsPerRow int) string {
	if monthsPerRow <= 0 {
		monthsPerRow = 3
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Calendar %d\n\n", year))

	for row := 1; row <= 12; row += monthsPerRow {
		for m := row; m < row+monthsPerRow && m <= 12; m++ {
			sb.WriteString(fmt.Sprintf("%-20s", MonthNames[m]))
		}
		sb.WriteString("\n")
		sb.WriteString("\n")
	}
	return sb.String()
}

// PrMonth prints a formatted month to the provided writer.
func PrMonth(w io.Writer, year, month int) {
	fmt.Fprint(w, FormatMonth(year, month))
}

// PrCalendar prints a full year calendar to the provided writer.
func PrCalendar(w io.Writer, year, monthsPerRow int) {
	fmt.Fprint(w, FormatYear(year, monthsPerRow))
}

//
// ---- Locale Support ----
//

// Locale represents a simple calendar localization.
type Locale struct {
	DayNames   []string
	DayAbbr    []string
	MonthNames []string
	MonthAbbr  []string
}

var currentLocale = Locale{
	DayNames:   DayNames,
	DayAbbr:    DayAbbr,
	MonthNames: MonthNames,
	MonthAbbr:  MonthAbbr,
}

// SetLocale sets the current calendar locale.
func SetLocale(l Locale) {
	currentLocale = l
}

// CurrentLocale returns the active locale.
func CurrentLocale() Locale {
	return currentLocale
}
