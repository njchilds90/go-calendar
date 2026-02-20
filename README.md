# go-calendar

[![Go Reference](https://pkg.go.dev/badge/github.com/njchilds90/go-calendar)](https://pkg.go.dev/github.com/njchilds90/go-calendar)
[![CI](https://github.com/njchilds90/go-calendar/actions/workflows/go.yml/badge.svg)](https://github.com/njchilds90/go-calendar/actions)

A lightweight, dependency-free Go port inspired by Python’s built-in `calendar` module.

This library provides calendar math, month matrices, formatted month/year rendering, iterators, weekday utilities, and optional locale customization — functionality not included in Go’s standard `time` package.

---

## Why This Exists

Go’s `time` package handles date arithmetic well but does not provide:

- Month calendar matrices (weeks × days)
- Text-formatted calendars
- Year-wide formatted calendars
- Iterator-style month traversal
- Built-in weekday name lists

This project fills that gap in a minimal, idiomatic way.

---

## Features

- ✔ Leap year detection (`IsLeap`)
- ✔ Leap year counting (`LeapDays`)
- ✔ Month range calculation (`MonthRange`)
- ✔ Month matrix generation (`MonthCalendar`)
- ✔ Month iterators (`IterMonthDays`, `IterMonthDates`)
- ✔ Configurable first weekday
- ✔ Week header formatting (`WeekHeader`)
- ✔ Formatted month output (`FormatMonth`)
- ✔ Formatted full year output (`FormatYear`)
- ✔ Print helpers (`PrMonth`, `PrCalendar`)
- ✔ Exported weekday/month name lists
- ✔ Basic locale customization support
- ✔ No external dependencies

---

## Installation

```bash
go get github.com/njchilds90/go-calendar
```

Import:

```go
import "github.com/njchilds90/go-calendar/calendar"
```

---

## Quick Example

```go
package main

import (
	"fmt"
	"os"

	"github.com/njchilds90/go-calendar/calendar"
)

func main() {
	// Leap year check
	fmt.Println(calendar.IsLeap(2024))

	// Month matrix
	grid := calendar.MonthCalendar(2026, 2)
	fmt.Println(grid)

	// Formatted month
	fmt.Println(calendar.FormatMonth(2026, 2))

	// Full year
	calendar.PrCalendar(os.Stdout, 2026, 3)
}
```

---

## Python Parity

This library provides equivalents of:

| Python | Go |
|--------|----|
| `isleap()` | `IsLeap()` |
| `leapdays()` | `LeapDays()` |
| `monthrange()` | `MonthRange()` |
| `monthcalendar()` | `MonthCalendar()` |
| `itermonthdays()` | `IterMonthDays()` |
| `itermonthdates()` | `IterMonthDates()` |
| `weekheader()` | `WeekHeader()` |
| `calendar(year)` | `FormatYear()` |
| `prcal()` | `PrCalendar()` |

---

## Locale Support

You can override default names:

```go
calendar.SetLocale(calendar.Locale{
    DayNames:   []string{"Lun", "Mar", "Mer", "Jeu", "Ven", "Sam", "Dim"},
    DayAbbr:    []string{"Lu", "Ma", "Me", "Je", "Ve", "Sa", "Di"},
    MonthNames: calendar.MonthNames,
    MonthAbbr:  calendar.MonthAbbr,
})
```

---

## Design Principles

- Minimal API
- Zero dependencies
- Python familiarity
- AI-friendly documentation
- Stable semantic versioning

---

## Testing

```
go test ./...
```

---

## Versioning

This project follows semantic versioning.

- `v1.x.x` — Stable API
- `v2.x.x` — Breaking changes

---

## Roadmap (Optional Future Additions)

- ISO week numbers
- Business/trading calendars
- Holiday registry support
- ICS export
- Exchange calendar integration

---

## License

MIT License
