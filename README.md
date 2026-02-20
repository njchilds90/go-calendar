# go-calendar

A simple Go port of Python's `calendar` module for calendar calculations and printing. Useful for generating month matrices, checking leap years, and printing text calendars.

## Installation

```sh
go get github.com/njchilds90/go-calendar
Usage
Import the package:
Goimport "github.com/njchilds90/go-calendar/calendar"
Examples:
Gopackage main

import (
	"fmt"
	"github.com/njchilds90/go-calendar/calendar"
)

func main() {
	fmt.Println(calendar.IsLeap(2024)) // true

	wd, days := calendar.MonthRange(2026, 2)
	fmt.Printf("First weekday: %d, Days: %d\n", wd, days) // First weekday: 6, Days: 28

	cal := calendar.MonthCalendar(2026, 2)
	fmt.Println(cal) // [[0 0 0 0 0 0 1] [2 3 4 5 6 7 8] ...]

	calendar.SetFirstWeekday(calendar.Monday)
	calendar.PrMonth(2026, 2, 2, 0) // Prints the month calendar
}
For full API, see godoc or the source.
Features

Leap year checks (IsLeap, LeapDays)
Weekday calculations
Month range and calendar matrix generation
Text calendar printing (FormatMonth, PrMonth)
Configurable first weekday
Day/month names and abbreviations

Running Tests
go test
License
MIT
