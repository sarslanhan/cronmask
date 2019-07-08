package cronmask

import (
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

var (
	validRanges = []struct {
		start, end int
	}{
		{start: 0, end: 59},
		{start: 0, end: 23},
		{start: 1, end: 31},
		{start: 1, end: 12},
		{start: 0, end: 6},
	}
)

// CronMask interface exposes a method to check whether the
// given time.Time matches the expression CronMask was constructed with.
type CronMask interface {
	Match(t time.Time) bool
}

type cronMask struct {
	Minute     cronField
	Hour       cronField
	DayOfMonth cronField
	Month      cronField
	DayOfWeek  cronField
}

type cronField interface {
	Match(val int) bool
}

type wildcardCronField struct {
}

func (*wildcardCronField) Match(val int) bool {
	return true
}

type constantCronField struct {
	Val int
}

func (f *constantCronField) Match(val int) bool {
	return f.Val == val
}

type rangeCronField struct {
	Start int
	End   int
}

func (f *rangeCronField) Match(val int) bool {
	return val >= f.Start && val <= f.End
}

type listCronField struct {
	Parts []cronField
}

func (f *listCronField) Match(val int) bool {
	for _, p := range f.Parts {
		if p.Match(val) {
			return true
		}
	}

	return false
}

func (c *cronMask) Match(t time.Time) bool {
	return c.Minute.Match(t.Minute()) &&
		c.Hour.Match(t.Hour()) &&
		c.DayOfMonth.Match(t.Day()) &&
		c.Month.Match(int(t.Month())) &&
		c.DayOfWeek.Match(int(t.Weekday()))
}

func parseCronField(fieldIdx int, fieldStr string) (cronField, error) {
	if fieldStr == "*" {
		return &wildcardCronField{}, nil
	}

	validateRange := func(i, val int) error {
		validRange := validRanges[i]
		if validRange.start <= val && validRange.end >= val {
			return nil
		}

		return errors.Errorf("expected %d to be in [%d,%d] range", val, validRange.start, validRange.end)
	}

	parts := strings.Split(fieldStr, ",")
	fields := make([]cronField, 0, len(parts))
	for _, p := range parts {
		if p == "" {
			return nil, errors.Errorf("could not parse the cron field: %s. invalid list item: %s", fieldStr, p)
		}
		possibleRangeFields := strings.Split(p, "-")

		if len(possibleRangeFields) == 1 {
			parsed, err := strconv.Atoi(possibleRangeFields[0])
			if err != nil {
				return nil, errors.Wrapf(err, "could not parse cron field: %s. invalid list item: %s", fieldStr, p)
			}
			if err := validateRange(fieldIdx, parsed); err != nil {
				return nil, err
			}
			fields = append(fields, &constantCronField{Val: parsed})
		} else if len(possibleRangeFields) == 2 {
			start, err := strconv.Atoi(possibleRangeFields[0])
			if err != nil {
				return nil, errors.Wrapf(err, "could not parse cron field: %s. invalid list item: %s", fieldStr, p)
			}
			if err := validateRange(fieldIdx, start); err != nil {
				return nil, err
			}
			end, err := strconv.Atoi(possibleRangeFields[1])
			if err != nil {
				return nil, errors.Wrapf(err, "could not parse cron field: %s. invalid list item: %s", fieldStr, p)
			}
			if err := validateRange(fieldIdx, end); err != nil {
				return nil, err
			}
			fields = append(fields, &rangeCronField{Start: start, End: end})
		} else {
			return nil, errors.Errorf("could not parse cron field: %s. invalid list item: %s", fieldStr, p)
		}
	}

	return &listCronField{Parts: fields}, nil
}

// New constructs a new CronMask instance that can be used to check if a given time.Time
// matches the expression or not.
//
// For CRON expressions (https://en.wikipedia.org/wiki/Cron#CRON_expression):
//
// Expressions are expected to be in the same time zone as the system that generates the time.Time instances.
//
// You can check the tests for what is possible.
//
// Unsupported features:
//
// - Non-standard characters (https://en.wikipedia.org/wiki/Cron#Non-standard_characters)
//
// - Year field
//
// - Command section
//
// - Text representation of the fields "month" and "day of week"
func New(expr string) (CronMask, error) {
	parts := strings.Fields(expr)

	if len(parts) != 5 {
		return nil, errors.New("invalid cron mask expression. expected 5 fields separated by whitespaces")
	}

	var fields [5]cronField
	for i, p := range parts {
		field, err := parseCronField(i, p)
		if err != nil {
			return nil, err
		}
		fields[i] = field
	}

	minute, hour, dayOfMonth, month, dayOfWeek := fields[0], fields[1], fields[2], fields[3], fields[4]

	return &cronMask{
		Minute:     minute,
		Hour:       hour,
		DayOfMonth: dayOfMonth,
		Month:      month,
		DayOfWeek:  dayOfWeek,
	}, nil
}
