package timeutils

import "time"

// Helper functions for time period comparisons.
func IsToday(t time.Time, now time.Time) bool {
	y1, m1, d1 := t.Date()
	y2, m2, d2 := now.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func IsYesterday(t time.Time, now time.Time) bool {
	yesterday := now.AddDate(0, 0, -1)
	y1, m1, d1 := t.Date()
	y2, m2, d2 := yesterday.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func IsSameWeek(t time.Time, now time.Time) bool {
	year1, week1 := t.ISOWeek()
	year2, week2 := now.ISOWeek()
	return year1 == year2 && week1 == week2
}

func IsLastWeek(t time.Time, now time.Time) bool {
	lastWeek := now.AddDate(0, 0, -7)
	year1, week1 := t.ISOWeek()
	year2, week2 := lastWeek.ISOWeek()
	return year1 == year2 && week1 == week2
}

func IsSameMonth(t time.Time, now time.Time) bool {
	return t.Year() == now.Year() && t.Month() == now.Month()
}

func IsLastMonth(t time.Time, now time.Time) bool {
	lastMonth := now.AddDate(0, -1, 0)
	return t.Year() == lastMonth.Year() && t.Month() == lastMonth.Month()
}

func FormatCQLTimeStamp(t time.Time) string {
	return t.Format("2025-04-11 03:55:15.058+0000")
}

func ParseTime(timeStr string) *time.Time {
	time, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return nil
	}
	return &time
}
