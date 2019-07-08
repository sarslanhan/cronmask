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
	minute     field
	hour       field
	dayOfMonth field
	month      field
	dayOfWeek  field
}

type field interface {
	match(val int) bool
}

type wildcard struct {
}

func (*wildcard) match(val int) bool {
	return true
}

type constant struct {
	val int
}

func (f *constant) match(val int) bool {
	return f.val == val
}

// range is a reserved keyword
type rangeF struct {
	start int
	end   int
}

func (f *rangeF) match(val int) bool {
	return val >= f.start && val <= f.end
}

type list struct {
	parts []field
}

func (f *list) match(val int) bool {
	for _, p := range f.parts {
		if p.match(val) {
			return true
		}
	}

	return false
}

func (c *cronMask) Match(t time.Time) bool {
	return c.minute.match(t.Minute()) &&
		c.hour.match(t.Hour()) &&
		c.dayOfMonth.match(t.Day()) &&
		c.month.match(int(t.Month())) &&
		c.dayOfWeek.match(int(t.Weekday()))
}

func parseCronField(fieldIdx int, fieldStr string) (field, error) {
	if fieldStr == "*" {
		return &wildcard{}, nil
	}

	validateRange := func(i, val int) error {
		validRange := validRanges[i]
		if validRange.start <= val && validRange.end >= val {
			return nil
		}

		return errors.Errorf("expected %d to be in [%d,%d] range", val, validRange.start, validRange.end)
	}

	parts := strings.Split(fieldStr, ",")
	fields := make([]field, 0, len(parts))
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
			fields = append(fields, &constant{val: parsed})
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
			fields = append(fields, &rangeF{start: start, end: end})
		} else {
			return nil, errors.Errorf("could not parse cron field: %s. invalid list item: %s", fieldStr, p)
		}
	}

	return &list{parts: fields}, nil
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

	var fields [5]field
	for i, p := range parts {
		field, err := parseCronField(i, p)
		if err != nil {
			return nil, err
		}
		fields[i] = field
	}

	minute, hour, dayOfMonth, month, dayOfWeek := fields[0], fields[1], fields[2], fields[3], fields[4]

	return &cronMask{
		minute:     minute,
		hour:       hour,
		dayOfMonth: dayOfMonth,
		month:      month,
		dayOfWeek:  dayOfWeek,
	}, nil
}
