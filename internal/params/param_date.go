package params

import (
	"fmt"
	"strings"
	"time"
)

type dateMode int

const (
	dateLayout = "02.01.2006"

	dateRangeSliceSize = 2

	after dateMode = iota
	before
	between
	exact
)

type dateParameter struct {
	field string
	mode  dateMode
	value time.Time
	to    time.Time
}

func newDateParameter(value, field string) (Parameter, error) {
	p := &dateParameter{
		field: field,
	}

	var (
		err  error
		date time.Time
		to   time.Time
	)

	switch {
	case strings.HasPrefix(value, ">"):
		p.mode = after
		date, err = p.parseDate(value, ">")

	case strings.HasPrefix(value, "<"):
		p.mode = before
		date, err = p.parseDate(value, "<")

	case strings.Contains(value, "-"):
		p.mode = between
		rawValues := strings.Split(value, "-")
		if len(rawValues) != dateRangeSliceSize {
			err = ErrInvalidDateRange
			return nil, err
		}

		date, err = p.parseDate(rawValues[0], "")
		if err != nil {
			return nil, err
		}

		to, err = p.parseDate(rawValues[1], "")

		if date.After(to) {
			date, to = to, date
		}

	default:
		p.mode = exact
		date, err = p.parseDate(value, "")
	}

	if err != nil {
		return nil, err
	}

	p.value = date
	p.to = to

	return p, nil
}

func (p dateParameter) ToSQL(num int) string {
	var query string

	switch p.mode {
	case after:
		query = fmt.Sprintf(" AND %s > $%d", p.field, num)

	case before:
		query = fmt.Sprintf(" AND %s < $%d", p.field, num)

	case between:
		query = fmt.Sprintf(" AND %s BETWEEN $%d AND $%d", p.field, num-1, num)

	case exact:
		query = fmt.Sprintf(" AND %s = $%d", p.field, num)
	}

	return query
}

func (p dateParameter) Values() []any {
	if p.mode == between {
		return []any{p.value, p.to}
	}
	return []any{p.value}
}

func (p dateParameter) parseDate(value, prefix string) (time.Time, error) {
	value = strings.TrimPrefix(value, prefix)

	date, err := time.Parse(dateLayout, value)
	if err != nil {
		return time.Time{}, ErrInvalidDateValue
	}

	return date, nil
}
