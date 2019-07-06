package conditions

import (
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

type CronMask interface {
	Matches(t time.Time) bool
}

type cronMask struct {
	Minute     cronField
	Hour       cronField
	DayOfMonth cronField
	Month      cronField
	DayOfWeek  cronField
}

type cronField interface {
	Matches(val int) bool
}

type wildcardCronField struct {
}

func (*wildcardCronField) Matches(val int) bool {
	return true
}

type constantCronField struct {
	Val int
}

func (f *constantCronField) Matches(val int) bool {
	return f.Val == val
}

type rangeCronField struct {
	Start int
	End   int
}

func (f *rangeCronField) Matches(val int) bool {
	return val >= f.Start && val <= f.End
}

type listCronField struct {
	Parts []cronField
}

func (f *listCronField) Matches(val int) bool {
	for _, p := range f.Parts {
		if p.Matches(val) {
			return true
		}
	}

	return false
}

func (c *cronMask) Matches(t time.Time) bool {
	utcTS := t.UTC()
	return c.Minute.Matches(utcTS.Minute()) &&
		c.Hour.Matches(utcTS.Hour()) &&
		c.DayOfMonth.Matches(utcTS.Day()) &&
		c.Month.Matches(int(utcTS.Month())) &&
		c.DayOfWeek.Matches(int(utcTS.Weekday()))
}

func parseCronField(fieldStr string) (cronField, error) {
	if fieldStr == "*" {
		return &wildcardCronField{}, nil
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
			fields = append(fields, &constantCronField{Val: parsed})
		} else if len(possibleRangeFields) == 2 {
			start, err := strconv.Atoi(possibleRangeFields[0])
			if err != nil {
				return nil, errors.Wrapf(err, "could not parse cron field: %s. invalid list item: %s", fieldStr, p)
			}
			end, err := strconv.Atoi(possibleRangeFields[1])
			if err != nil {
				return nil, errors.Wrapf(err, "could not parse cron field: %s. invalid list item: %s", fieldStr, p)
			}
			fields = append(fields, &rangeCronField{Start: start, End: end})
		} else {
			return nil, errors.Errorf("could not parse cron field: %s. invalid list item: %s", fieldStr, p)
		}
	}

	return &listCronField{Parts: fields}, nil
}

func New(expr string) (CronMask, error) {
	parts := strings.Fields(expr)

	if len(parts) != 5 {
		return nil, errors.New("invalid cron mask expression. expected 5 fields separated by whitespaces")
	}

	var fields [5]cronField
	for i, p := range parts {
		field, err := parseCronField(p)
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