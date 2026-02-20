# go-calendar

[![Go Reference](https://pkg.go.dev/badge/github.com/njchilds90/go-calendar)](https://pkg.go.dev/github.com/njchilds90/go-calendar)

A simple, lightweight Go library inspired by Python’s built-in `calendar` module — providing calendar math, month matrices, and formatted text calendars.

Go’s standard library provides powerful date/time types (`time.Time`), but it lacks utilities to generate calendar grids, text calendars, and flexible weekday routines — exactly what this package gives you.

---

## 🚀 Features

✔ Determine leap years  
✔ Get first weekday and number of days in a month  
✔ Generate calendar matrices (weeks × days)  
✔ Format month calendars as text  
✔ Configurable first weekday (e.g., Monday or Sunday)  
✔ Weekday and month names / abbreviations

This library is designed to feel familiar to Python developers while being idiomatic Go for CLI tools, backend services, scheduling utilities, and more.

---

## 📦 Installation

```sh
go get github.com/njchilds90/go-calendar
```

Import it anywhere in your project:

```go
import "github.com/njchilds90/go-calendar/calendar"
```

---

## 💡 Quick Usage

```go
package main

import (
	"fmt"
	"github.com/njchilds90/go-calendar/calendar"
)

func main() {
	// Check for leap year
	fmt.Println("2024 leap year?", calendar.IsLeap(2024))

	// Month info
	weekday, days := calendar.MonthRange(2026, 2)
	fmt.Printf("Feb 2026 starts on weekday %d and has %d days\n", weekday, days)

	// Generate 2D calendar matrix
	grid := calendar.MonthCalendar(2026, 2)
	fmt.Println("Calendar grid:", grid)

	// Format text calendar
	calendar.SetFirstWeekday(calendar.Monday)
	text := calendar.FormatMonth(2026, 2)
	fmt.Println(text)
}
```

---

## 📖 API Overview

| Function / Method | Description |
|------------------|-------------|
| `IsLeap(year int) bool` | Returns true if `year` is a leap year |
| `LeapDays(start, end int) int` | Number of leap years between two years |
| `MonthRange(year, month int) (weekday, days int)` | First weekday and number of days in a month |
| `MonthCalendar(year, month int) [][]int` | Returns a matrix of weeks (0 padding for days outside month) |
| `FormatMonth(year, month int) string` | Returns human-readable text calendar |
| `SetFirstWeekday(day Weekday)` | Set which day the weeks start on (default: Monday) |

*(See GoDoc for full list and examples.)*

---

## 📈 Why This Library Matters

Go’s `time` package is great, but it doesn’t provide:

* Calendar grids (matrix of weeks × days) like Python’s `calendar.monthcalendar`  
* Preformatted text month views (like *ncal* or Python’s `prmonth`)  
* Easy control over first weekday and weekday names

This fills that gap with minimal dependencies and familiar patterns for Python users working in Go.

---

## 📚 Testing

Run:

```sh
go test ./...
```

All tests should pass on standard Go tooling.

---

## 🛠 Release Checklist

Before making a release (e.g., `v0.1.0`):

✔ Add semantic version tags (e.g., `v0.1.0`, `v0.2.0`)  
✔ Ensure all public functions have GoDoc comments  
✔ Add more usage examples or an `examples/` folder  
✔ Add badges: Go Reference, CI status, coverage  
✔ Consider documenting behavior for edge cases (e.g., negative years)

---

## 🧠 Roadmap & Enhancements

Future additions could include:

* ISO week numbers  
* Locale-specific weekday/month names  
* Holiday support (workday calendars, business calendars)  
* ICS / iCal import/export utilities

---

## 🌍 Inspiration

This library is inspired by Python’s `calendar` module, which offers similar utilities in the Python standard library — functions for printing calendars and working with month/week layouts. :contentReference[oaicite:0]{index=0}

---

## 📬 Contributing

Contributions, issues, and feature requests are welcome! Please open them on the GitHub repository.

---

## 📜 License

This project is open-source under the MIT license.
