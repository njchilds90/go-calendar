// Package calendar provides calendar utilities inspired by Python's calendar module.
// It supports leap year checks, month matrices, iterators over month days/dates, text and HTML formatting,
// year-at-a-glance views, configurable first weekday, basic locale support, and a simple holiday registry.
//
// All operations use Go's time package (UTC-based).
package calendar

import (
	"fmt"
	"strings"
	"time"
)

// Weekday constants (matches time.Weekday)
const (
	Sunday    = 0
	Monday    = 1
	Tuesday   = 2
	Wednesday = 3
	Thursday  = 4
	Friday    = 5
	Saturday  = 6
)

// Exported name lists (can be overridden directly or via Locale)
var (
	DayNames   = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	DayAbbrs   = []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	MonthNames = []string{"", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	MonthAbbrs = []string{"", "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
)

// Locale for name overrides
type Locale struct {
	DayNames   []string // len=7
	DayAbbrs   []string // len=7
	MonthNames []string // len=13 (0 unused)
	MonthAbbrs []string // len=13 (0 unused)
}

var currentLocale = Locale{
	DayNames:   DayNames,
	DayAbbrs:   DayAbbrs,
	MonthNames: MonthNames,
	MonthAbbrs: MonthAbbrs,
}

// SetLocale updates global names (panics on invalid lengths)
func SetLocale(l Locale) {
	if len(l.DayNames) != 7 || len(l.DayAbbrs) != 7 ||
		len(l.MonthNames) != 13 || len(l.MonthAbbrs) != 13 {
		panic("invalid locale: must have 7 days and 13 months")
	}
	currentLocale = l
}

var firstWeekday = Monday

func SetFirstWeekday(wd int) {
	if wd < 0 || wd > 6 {
		panic("invalid weekday")
	}
	firstWeekday = wd
}

func FirstWeekday() int {
	return firstWeekday
}

// ── Core Utils ───────────────────────────────────────────────────────────────

func IsLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func LeapDays(y1, y2 int) int {
	f := func(y int) int { return y/4 - y/100 + y/400 }
	return f(y2) - f(y1)
}

func Weekday(year, month, day int) int {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return int(t.Weekday())
}

func MonthRange(year, month int) (firstWd, days int) {
	if month < 1 || month > 12 {
		panic("invalid month")
	}
	days = 31
	switch month {
	case 4, 6, 9, 11:
		days = 30
	case 2:
		days = 28
		if IsLeap(year) {
			days = 29
		}
	}
	firstWd = Weekday(year, month, 1)
	return
}

func MonthCalendar(year, month int) [][]int {
	firstWd, days := MonthRange(year, month)
	cal := make([][]int, 0, 6)
	day := 1
	shift := (firstWd - firstWeekday + 7) % 7
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

// ── Iterators (full Python equivalents) ──────────────────────────────────────

// IterWeekdays returns an iterator over the 7 weekday numbers starting from firstWeekday.
func IterWeekdays() <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 0; i < 7; i++ {
			ch <- (firstWeekday + i) % 7
		}
	}()
	return ch
}

// IterMonthDays yields (day_number, weekday) pairs for the month, including padding zeros.
func IterMonthDays(year, month int) <-chan [2]int {
	ch := make(chan [2]int)
	go func() {
		defer close(ch)
		firstWd, days := MonthRange(year, month)
		shift := (firstWd - firstWeekday + 7) % 7
		wd := firstWeekday
		for i := 0; i < shift+days; i++ {
			day := 0
			if i >= shift {
				day = i - shift + 1
			}
			ch <- [2]int{day, wd}
			wd = (wd + 1) % 7
		}
	}()
	return ch
}

// IterMonthDates yields time.Time values for every day in the month + padding days to fill weeks.
func IterMonthDates(year, month int) <-chan time.Time {
	ch := make(chan time.Time)
	go func() {
		defer close(ch)
		firstWd, days := MonthRange(year, month)
		shift := (firstWd - firstWeekday + 7) % 7
		startDay := 1 - shift
		for i := 0; i < 42; i++ { // max 6 weeks
			d := startDay + i
			if d < 1 || d > days {
				// Padding: yield zero time or skip; here we yield valid dates outside month
				ch <- time.Date(year, time.Month(month), d, 0, 0, 0, 0, time.UTC)
			} else {
				ch <- time.Date(year, time.Month(month), d, 0, 0, 0, 0, time.UTC)
			}
			if i >= shift+days-1 && (i+1)%7 == 0 {
				break
			}
		}
	}()
	return ch
}

// ── Formatting ───────────────────────────────────────────────────────────────

func WeekHeader(width int) string {
	var sb strings.Builder
	for i := 0; i < 7; i++ {
		wd := (firstWeekday + i) % 7
		abbr := currentLocale.DayAbbrs[wd]
		if len(abbr) > width {
			abbr = abbr[:width]
		}
		fmt.Fprintf(&sb, "%*s ", width, abbr)
	}
	s := sb.String()
	return s[:len(s)-1]
}

func FormatMonth(year, month, width, lines int) string {
	if width < 2 {
		width = 2
	}
	header := fmt.Sprintf("%s %d", currentLocale.MonthNames[month], year)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%*s\n", (7*(width+1)-1+len(header))/2, header))
	sb.WriteString(WeekHeader(width) + "\n")
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

func PrMonth(year, month, width, lines int) {
	fmt.Print(FormatMonth(year, month, width, lines))
}

func FormatYear(year, width, lines, monthsPerRow int) string {
	// ... (your existing compact year text implementation - keep as-is)
	// For brevity, assume it's already in your file; if not, use previous versions.
	return "" // placeholder - replace with your impl
}

func PrCalendar(w fmt.Writer, year int, monthsPerRow int) {
	// ... (your print helper - keep as-is)
}

// ── Holiday Support (simple registry) ────────────────────────────────────────

var holidays = make(map[time.Time]string) // date → holiday name

// RegisterHoliday adds a holiday (date → name). Overwrites if exists.
func RegisterHoliday(t time.Time, name string) {
	// Normalize to date-only (zero time)
	d := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	holidays[d] = name
}

// IsHoliday checks if the date is registered as a holiday.
func IsHoliday(t time.Time) (bool, string) {
	d := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	name, ok := holidays[d]
	return ok, name
}

// ClearHolidays removes all registered holidays.
func ClearHolidays() {
	holidays = make(map[time.Time]string)
}

// Example usage in iterators: extend IterMonthDates to yield holiday info if desired
// (you can add a variant like IterMonthWithHolidays later)

// ── HTML Support (keep your existing if present) ─────────────────────────────

type HTMLCalendar struct {
	Firstweekday int
	// ... your existing fields/methods
}

func NewHTMLCalendar(fw int) *HTMLCalendar {
	// ... your impl
	return nil // placeholder
}

// In FormatMonthHTML / FormatYearHTML, you can now check IsHoliday for each day and add e.g. class="holiday"
