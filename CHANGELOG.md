# Changelog

All notable changes to this project will be documented in this file.

This project follows Semantic Versioning.

---

## [1.0.0] - 2026-02-20

### Added

- Leap year detection (`IsLeap`)
- Leap year counting (`LeapDays`)
- Month range calculation (`MonthRange`)
- Month matrix generation (`MonthCalendar`)
- Iterator APIs:
  - `IterMonthDays`
  - `IterMonthDates`
- Week header formatting (`WeekHeader`)
- Formatted month rendering (`FormatMonth`)
- Formatted year rendering (`FormatYear`)
- Print helpers:
  - `PrMonth`
  - `PrCalendar`
- Exported weekday/month name lists
- Configurable first weekday
- Basic locale customization
- GitHub Actions CI
- Full GoDoc documentation

### Notes

This release provides functional parity with the core features of Python’s standard `calendar` module while maintaining a minimal, dependency-free Go implementation.
