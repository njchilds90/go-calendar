package calendar

import (
	"testing"
)

func TestIsLeap(t *testing.T) {
	tests := []struct {
		year int
		want bool
	}{
		{2000, true},
		{1900, false},
		{2004, true},
		{2005, false},
	}
	for _, tt := range tests {
		if got := IsLeap(tt.year); got != tt.want {
			t.Errorf("IsLeap(%d) = %v, want %v", tt.year, got, tt.want)
		}
	}
}

func TestLeapDays(t *testing.T) {
	tests := []struct {
		y1, y2 int
		want   int
	}{
		{2000, 2005, 2},
		{1900, 2000, 25},
	}
	for _, tt := range tests {
		if got := LeapDays(tt.y1, tt.y2); got != tt.want {
			t.Errorf("LeapDays(%d, %d) = %v, want %v", tt.y1, tt.y2, got, tt.want)
		}
	}
}

func TestWeekday(t *testing.T) {
	tests := []struct {
		y, m, d int
		want    int
	}{
		{2000, 1, 1, Saturday},
		{2026, 2, 20, Friday},
	}
	for _, tt := range tests {
		if got := Weekday(tt.y, tt.m, tt.d); got != tt.want {
			t.Errorf("Weekday(%d,%d,%d) = %v, want %v", tt.y, tt.m, tt.d, got, tt.want)
		}
	}
}

func TestMonthRange(t *testing.T) {
	tests := []struct {
		y, m     int
		wantWd   int
		wantDays int
	}{
		{2000, 2, Tuesday, 29},
		{2026, 2, Friday, 28},
	}
	for _, tt := range tests {
		wd, days := MonthRange(tt.y, tt.m)
		if wd != tt.wantWd || days != tt.wantDays {
			t.Errorf("MonthRange(%d,%d) = %v,%v want %v,%v", tt.y, tt.m, wd, days, tt.wantWd, tt.wantDays)
		}
	}
}

func TestMonthCalendar(t *testing.T) {
	SetFirstWeekday(Monday)
	cal := MonthCalendar(2026, 2)
	want := [][]int{
		{0, 0, 0, 0, 0, 0, 1},
		{2, 3, 4, 5, 6, 7, 8},
		{9, 10, 11, 12, 13, 14, 15},
		{16, 17, 18, 19, 20, 21, 22},
		{23, 24, 25, 26, 27, 28, 0},
	}
	if len(cal) != len(want) {
		t.Fatalf("len = %d, want %d", len(cal), len(want))
	}
	for i := range cal {
		for j := range cal[i] {
			if cal[i][j] != want[i][j] {
				t.Errorf("cal[%d][%d] = %d, want %d", i, j, cal[i][j], want[i][j])
			}
		}
	}
}

func TestFormatMonth(t *testing.T) {
	SetFirstWeekday(Monday)
	got := FormatMonth(2026, 2, 2, 0)
	want := `   February 2026
Mon Tue Wed Thu Fri Sat Sun
                    1
  2   3   4   5   6   7   8
  9  10  11  12  13  14  15
 16  17  18  19  20  21  22
 23  24  25  26  27  28   

`
	if got != want {
		t.Errorf("FormatMonth = %q, want %q", got, want)
	}
}
