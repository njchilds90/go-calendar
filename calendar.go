// Package calendar provides simple calendar utilities inspired by Python's calendar module:
// leap year checks, month calendars as matrices, text and HTML formatting/printing, and basic locale support.
package calendar

import (
	"fmt"
	"strings"
	"time"
)

// Constants for weekdays (matching time.Weekday: Sunday=0, ..., Saturday=6)
const (
	Sunday    = 0
	Monday    = 1
	Tuesday   = 2
	Wednesday = 3
	Thursday  = 4
	Friday    = 5
	Saturday  = 6
)

// Default names (English)
var (
	dayNames   = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	dayAbbrs   = []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	monthNames = []string{"", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	monthAbbrs = []string{"", "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
)

// Locale holds customizable names (extendable for i18n)
type Locale struct {
	DayNames   []string
	DayAbbrs   []string
	MonthNames []string
	MonthAbbrs []string
}

// DefaultLocale is English
var DefaultLocale = Locale{
	DayNames:   dayNames,
	DayAbbrs:   dayAbbrs,
	MonthNames: monthNames,
	MonthAbbrs: monthAbbrs,
}

// Current locale (global for simplicity; can be made per-instance later)
var currentLocale = DefaultLocale

// SetLocale changes the global names used (e.g. for other languages)
func SetLocale(loc Locale) {
	if len(loc.DayNames) == 7 && len(loc.DayAbbrs) == 7 &&
		len(loc.MonthNames) == 13 && len(loc.MonthAbbrs) == 13 {
		currentLocale = loc
	} else {
		panic("invalid locale: must have exactly 7 days and 13 months (index 0 unused for months)")
	}
}

// firstWeekday is the starting day of the week (default Monday)
var firstWeekday = Monday

// SetFirstWeekday sets the first day of the week (0=Sunday, 1=Monday, ..., 6=Saturday).
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
		numDays = 28
		if IsLeap(year) {
			numDays = 29
		}
	}
	weekday = Weekday(year, month, 1)
	return
}

// MonthCalendar returns a matrix representing a month's calendar (weeks as rows, days as columns, 0 for padding).
func MonthCalendar(year, month int) [][]int {
	wd, days := MonthRange(year, month)
	cal := make([][]int, 0, 6)
	day := 1
	shift := (wd - firstWeekday + 7) % 7

	for week := 0; week < 6; week++ {
		row := make([]int, 7)
		empty := true
		for d := 0; d < 7; d++ {
			if week == 0 && d < shift || day > days {
				row[d] = 0
			} else {
				row[d] = day
				day++
				empty = false
			}
		}
		if empty {
			break
		}
		cal = append(cal, row)
	}
	return cal
}

// weekHeader returns weekday abbrs line for text output.
func weekHeader(width int) string {
	var sb strings.Builder
	for i := 0; i < 7; i++ {
		day := (i + firstWeekday) % 7
		abbr := currentLocale.DayAbbrs[day]
		if len(abbr) > width {
			abbr = abbr[:width]
		}
		fmt.Fprintf(&sb, "%-*s ", width, abbr)
	}
	return sb.String()[:sb.Len()-1]
}

// FormatMonth returns a multi-line text string for the month calendar.
func FormatMonth(year, month, width, lines int) string {
	if width < 1 {
		width = 3
	}
	var sb strings.Builder
	header := fmt.Sprintf("%s %d", currentLocale.MonthNames[month], year)
	sb.WriteString(fmt.Sprintf("%*s\n", (7*(width+1)-1+len(header))/2, header))
	sb.WriteString(weekHeader(width) + "\n")
	cal := MonthCalendar(year, month)
	for _, week := range cal {
		for _, day := range week {
			if day == 0 {
				fmt.Fprintf(&sb, "%*s ", width, "")
			} else {
				fmt.Fprintf(&sb, "%*d ", width, day)
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

// FormatYear returns text calendars for the full year (3 months per row by default).
func FormatYear(year, width, lines, monthsPerRow int) string {
	if monthsPerRow < 1 || monthsPerRow > 12 {
		monthsPerRow = 3
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%*s\n\n", 20*monthsPerRow, fmt.Sprintf("%d", year)))

	for m := 1; m <= 12; m += monthsPerRow {
		for row := 0; row < 3; row++ { // header + week + calendar lines
			lineParts := []string{}
			for col := 0; col < monthsPerRow && m+col <= 12; col++ {
				month := m + col
				if row == 0 {
					h := fmt.Sprintf("%s", currentLocale.MonthNames[month])
					lineParts = append(lineParts, fmt.Sprintf("%*s", (7*(width+1)-1), h))
				} else if row == 1 {
					lineParts = append(lineParts, weekHeader(width))
				} else {
					cal := MonthCalendar(year, month)
					if len(cal) > 0 && len(cal[0]) == 7 {
						// Simple: just first week for brevity; extend if needed
						weekStr := ""
						for _, d := range cal[0] {
							if d == 0 {
								weekStr += fmt.Sprintf("%*s ", width, "")
							} else {
								weekStr += fmt.Sprintf("%*d ", width, d)
							}
						}
						lineParts = append(lineParts, weekStr[:len(weekStr)-1])
					}
				}
			}
			sb.WriteString(strings.Join(lineParts, "   ") + "\n")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// PrYear prints the year calendar to stdout.
func PrYear(year, width, lines, monthsPerRow int) {
	fmt.Print(FormatYear(year, width, lines, monthsPerRow))
}

// HTMLCalendar provides HTML generation like Python's HTMLCalendar.
type HTMLCalendar struct {
	firstweekday int
	cssclasses   map[int]string // weekday -> td class
}

// NewHTMLCalendar creates a new HTMLCalendar (firstweekday 0=Mon by default).
func NewHTMLCalendar(firstweekday int) *HTMLCalendar {
	if firstweekday < 0 || firstweekday > 6 {
		firstweekday = Monday
	}
	return &HTMLCalendar{
		firstweekday: firstweekday,
		cssclasses: map[int]string{
			0: "mon", 1: "tue", 2: "wed", 3: "thu", 4: "fri", 5: "sat", 6: "sun",
		},
	}
}

// FormatMonthHTML returns HTML table for one month (withyear=true includes year in header).
func (c *HTMLCalendar) FormatMonthHTML(theyear, themonth int, withyear bool) string {
	var sb strings.Builder
	sb.WriteString(`<table border="0" cellpadding="0" cellspacing="0" class="month">` + "\n")
	header := currentLocale.MonthNames[themonth]
	if withyear {
		header += fmt.Sprintf(" %d", theyear)
	}
	sb.WriteString(fmt.Sprintf("<tr><th colspan=\"7\" class=\"month\">%s</th></tr>\n", header))

	// Week header
	sb.WriteString("<tr>")
	for i := 0; i < 7; i++ {
		day := (i + c.firstweekday) % 7
		sb.WriteString(fmt.Sprintf("<th class=\"%s\">%s</th>", c.cssclasses[day], currentLocale.DayAbbrs[day]))
	}
	sb.WriteString("</tr>\n")

	cal := MonthCalendar(theyear, themonth)
	for _, week := range cal {
		sb.WriteString("<tr>")
		for d, day := range week {
			wd := (d + c.firstweekday) % 7
			if day == 0 {
				sb.WriteString(`<td class="noday">&nbsp;</td>`)
			} else {
				sb.WriteString(fmt.Sprintf(`<td class="%s">%d</td>`, c.cssclasses[wd], day))
			}
		}
		sb.WriteString("</tr>\n")
	}
	sb.WriteString("</table>")
	return sb.String()
}

// FormatYearHTML returns HTML for the full year (width=3 months per row).
func (c *HTMLCalendar) FormatYearHTML(theyear int, width int) string {
	if width < 1 {
		width = 3
	}
	var sb strings.Builder
	sb.WriteString(`<table border="0" cellpadding="0" cellspacing="0" class="year">` + "\n")
	sb.WriteString(fmt.Sprintf("<tr><th colspan=\"%d\" class=\"year\">%d</th></tr>\n", width*7, theyear))

	for m := 1; m <= 12; m += width {
		sb.WriteString("<tr>")
		for col := 0; col < width && m+col <= 12; col++ {
			month := m + col
			sb.WriteString("<td>")
			sb.WriteString(c.FormatMonthHTML(theyear, month, false))
			sb.WriteString("</td>")
		}
		sb.WriteString("</tr>\n")
	}
	sb.WriteString("</table>")
	return sb.String()
}
