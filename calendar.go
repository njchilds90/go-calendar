// Package calendar provides calendar utilities inspired by Python's calendar module.
// It includes leap year checks, weekday calculations, month calendar matrices,
// configurable text and HTML formatting, year-at-a-glance views, and basic locale support
// for month and day names.
//
// All operations are timezone-agnostic (using UTC via time.Date).
//
// Example:
//
//	calendar.SetFirstWeekday(calendar.Monday)
//	calendar.PrMonth(2026, 2, 3, 0) // prints text calendar
//	hc := calendar.NewHTMLCalendar(calendar.Monday)
//	html := hc.FormatMonthHTML(2026, 2, true)
package calendar

import (
	"fmt"
	"strings"
	"time"
)

// Weekday constants (matches time.Weekday values)
const (
	Sunday    = 0
	Monday    = 1
	Tuesday   = 2
	Wednesday = 3
	Thursday  = 4
	Friday    = 5
	Saturday  = 6
)

// Default English names
var (
	defaultDayNames   = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	defaultDayAbbrs   = []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	defaultMonthNames = []string{"", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	defaultMonthAbbrs = []string{"", "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
)

// Locale defines customizable names for days and months (index 0 unused for months).
type Locale struct {
	DayNames   []string // 7 elements
	DayAbbrs   []string // 7 elements
	MonthNames []string // 13 elements (0 unused)
	MonthAbbrs []string // 13 elements (0 unused)
}

// DefaultLocale is English.
var DefaultLocale = Locale{
	DayNames:   defaultDayNames,
	DayAbbrs:   defaultDayAbbrs,
	MonthNames: defaultMonthNames,
	MonthAbbrs: defaultMonthAbbrs,
}

// currentLocale is the active set of names.
var currentLocale = DefaultLocale

// SetLocale updates the global locale names. Panics on invalid lengths.
func SetLocale(loc Locale) {
	if len(loc.DayNames) != 7 || len(loc.DayAbbrs) != 7 ||
		len(loc.MonthNames) != 13 || len(loc.MonthAbbrs) != 13 {
		panic("Locale must have 7 days and 13 months (index 0 unused for months)")
	}
	currentLocale = loc
}

// firstWeekday is the week starting day (default Monday).
var firstWeekday = Monday

// SetFirstWeekday sets the starting weekday (0=Sunday to 6=Saturday).
func SetFirstWeekday(wd int) {
	if wd < 0 || wd > 6 {
		panic("weekday must be 0-6")
	}
	firstWeekday = wd
}

// FirstWeekday returns the current starting weekday.
func FirstWeekday() int {
	return firstWeekday
}

// IsLeap reports whether the year is a leap year.
func IsLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// LeapDays counts leap years in the half-open range [y1, y2).
func LeapDays(y1, y2 int) int {
	f := func(y int) int { return y/4 - y/100 + y/400 }
	return f(y2) - f(y1)
}

// Weekday returns the weekday for the date (0=Sunday ... 6=Saturday).
func Weekday(year, month, day int) int {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return int(t.Weekday())
}

// MonthRange returns (first_weekday, days_in_month) for year/month.
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

// MonthCalendar returns a matrix (up to 6 rows × 7 cols) for the month; 0 = padding.
func MonthCalendar(year, month int) [][]int {
	wd, days := MonthRange(year, month)
	cal := make([][]int, 0, 6)
	day := 1
	shift := (wd - firstWeekday + 7) % 7
	for w := 0; w < 6; w++ {
		row := make([]int, 7)
		empty := true
		for d := 0; d < 7; d++ {
			if w*7+d < shift || day > days {
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

// weekHeader returns the weekday abbr line for text calendars.
func weekHeader(width int) string {
	var sb strings.Builder
	for i := 0; i < 7; i++ {
		wd := (i + firstWeekday) % 7
		abbr := currentLocale.DayAbbrs[wd]
		if len(abbr) > width {
			abbr = abbr[:width]
		}
		fmt.Fprintf(&sb, "%*s ", width, abbr)
	}
	s := sb.String()
	return s[:len(s)-1]
}

// FormatMonth returns formatted text for one month.
func FormatMonth(year, month, width, lines int) string {
	if width < 2 {
		width = 2
	}
	header := fmt.Sprintf("%s %d", currentLocale.MonthNames[month], year)
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

// FormatYear returns text year-at-a-glance (compact: one line per month header + week).
func FormatYear(year, width, lines, monthsPerRow int) string {
	if monthsPerRow < 1 || monthsPerRow > 12 {
		monthsPerRow = 3
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%*s\n\n", monthsPerRow*20, fmt.Sprintf("%d", year)))
	for m := 1; m <= 12; m += monthsPerRow {
		for row := 0; row < 3; row++ {
			parts := []string{}
			for col := 0; col < monthsPerRow && m+col <= 12; col++ {
				month := m + col
				if row == 0 {
					h := fmt.Sprintf("%s", currentLocale.MonthNames[month])
					parts = append(parts, fmt.Sprintf("%*s", 7*(width+1)-1, h))
				} else if row == 1 {
					parts = append(parts, weekHeader(width))
				} else {
					cal := MonthCalendar(year, month)
					if len(cal) > 0 {
						wstr := ""
						for _, d := range cal[0] {
							if d == 0 {
								wstr += fmt.Sprintf("%*s ", width, "")
							} else {
								wstr += fmt.Sprintf("%*d ", width, d)
							}
						}
						parts = append(parts, wstr[:len(wstr)-1])
					}
				}
			}
			sb.WriteString(strings.Join(parts, "   ") + "\n")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// PrYear prints the year calendar to stdout.
func PrYear(year, width, lines, monthsPerRow int) {
	fmt.Print(FormatYear(year, width, lines, monthsPerRow))
}

// HTMLCalendar generates HTML tables like Python's HTMLCalendar.
type HTMLCalendar struct {
	firstweekday int
	cssclasses   map[int]string
}

// NewHTMLCalendar creates one with given firstweekday (default Monday).
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

// FormatMonthHTML returns HTML <table> for one month (withyear includes year in title).
func (c *HTMLCalendar) FormatMonthHTML(year, month int, withyear bool) string {
	var sb strings.Builder
	sb.WriteString(`<table border="0" cellpadding="0" cellspacing="0" class="month">` + "\n")
	title := currentLocale.MonthNames[month]
	if withyear {
		title += fmt.Sprintf(" %d", year)
	}
	sb.WriteString(fmt.Sprintf(`<tr><th colspan="7" class="month">%s</th></tr>`+"\n", title))
	sb.WriteString("<tr>")
	for i := 0; i < 7; i++ {
		wd := (i + c.firstweekday) % 7
		sb.WriteString(fmt.Sprintf(`<th class="%s">%s</th>`, c.cssclasses[wd], currentLocale.DayAbbrs[wd]))
	}
	sb.WriteString("</tr>\n")
	cal := MonthCalendar(year, month)
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

// FormatYearHTML returns full-year HTML (monthsPerRow columns).
func (c *HTMLCalendar) FormatYearHTML(year int, monthsPerRow int) string {
	if monthsPerRow < 1 {
		monthsPerRow = 3
	}
	var sb strings.Builder
	sb.WriteString(`<table border="0" cellpadding="0" cellspacing="0" class="year">` + "\n")
	sb.WriteString(fmt.Sprintf(`<tr><th colspan="%d" class="year">%d</th></tr>`+"\n", monthsPerRow*7, year))
	for m := 1; m <= 12; m += monthsPerRow {
		sb.WriteString("<tr>")
		for col := 0; col < monthsPerRow && m+col <= 12; col++ {
			sb.WriteString("<td valign=\"top\">")
			sb.WriteString(c.FormatMonthHTML(year, m+col, false))
			sb.WriteString("</td>")
		}
		sb.WriteString("</tr>\n")
	}
	sb.WriteString("</table>")
	return sb.String()
}
