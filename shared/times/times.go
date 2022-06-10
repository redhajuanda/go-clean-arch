package times

import (
	"time"
	_ "time/tzdata"

	"github.com/pkg/errors"
)

const (
	DATE_FORMAT = "2006-01-02"
)

func Now() time.Time {
	return time.Now().UTC()
}

func NowJkt() time.Time {
	return time.Now().In(LocJkt())
}

func ParseDateString(oldFormat string, newFormat string, date string) (string, error) {
	dateTime, err := time.Parse(oldFormat, date)
	if err != nil {
		return "", errors.Wrap(err, "cannot parse date string")
	}
	return dateTime.Format(newFormat), nil
}

func LocJkt() *time.Location {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err)
	}
	return loc
}

// MonthsCountUntil calculates the months between now
// and the until time.Time value passed
func MonthsCountUntil(until time.Time) int {
	now := time.Now()
	months := 0
	month := until.Month()
	for now.Before(until) {
		now = now.Add(time.Hour * 24)
		nextMonth := now.Month()
		if nextMonth != month {
			months++
		}
		month = nextMonth
	}

	if months == 0 {
		months = 1
	}

	return months
}
