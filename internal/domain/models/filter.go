package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/FlutterDizaster/music-library/internal/apperrors"
)

type dateMode int

const (
	exact dateMode = 1 + iota
	after
	before
	between

	dateRangeSliceSize = 2

	DateLayout = "02.01.2006"

	orderingField = "title"
)

type RawFilters struct {
	Title       string
	Group       string
	ReleaseDate string
	Text        string
	Link        string
	Limit       string
	Offset      string
}

type Filters struct {
	title    string
	group    string
	dateFrom time.Time
	dateTo   time.Time
	dateMode dateMode
	text     string
	link     string
	limit    int
	offset   int
}

func (r RawFilters) ToFilters() (*Filters, error) {
	f := &Filters{
		title:    r.Title,
		group:    r.Group,
		text:     r.Text,
		link:     r.Link,
		dateFrom: time.Time{},
		dateTo:   time.Time{},
	}

	var err error

	if r.Limit != "" {
		f.limit, err = strconv.Atoi(r.Limit)
		if err != nil {
			return nil, apperrors.ErrInvalidFilters
		}
	}

	if r.Offset != "" {
		f.offset, err = strconv.Atoi(r.Offset)
		if err != nil {
			return nil, apperrors.ErrInvalidFilters
		}
	}

	if r.ReleaseDate != "" {
		err = r.parseDate(f)
		if err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (r RawFilters) parseDate(f *Filters) error {
	var err error
	switch {
	case strings.HasPrefix(r.ReleaseDate, ">"):
		f.dateMode = after
		f.dateFrom, err = time.Parse(DateLayout, strings.TrimPrefix(r.ReleaseDate, ">"))
		if err != nil {
			return apperrors.ErrInvalidDateLayout
		}

	case strings.HasPrefix(r.ReleaseDate, "<"):
		f.dateMode = before
		f.dateFrom, err = time.Parse(DateLayout, strings.TrimPrefix(r.ReleaseDate, "<"))
		if err != nil {
			return apperrors.ErrInvalidDateLayout
		}

	case strings.Contains(r.ReleaseDate, "-"):
		f.dateMode = between
		dates := strings.Split(r.ReleaseDate, "-")
		if len(dates) != dateRangeSliceSize {
			return apperrors.ErrInvalidDateRange
		}

		f.dateFrom, err = time.Parse(DateLayout, dates[0])
		if err != nil {
			return apperrors.ErrInvalidDateLayout
		}

		f.dateTo, err = time.Parse(DateLayout, dates[1])
		if err != nil {
			return apperrors.ErrInvalidDateLayout
		}

		if f.dateFrom.After(f.dateTo) {
			f.dateFrom, f.dateTo = f.dateTo, f.dateFrom
		}

	default:
		f.dateFrom, err = time.Parse(DateLayout, r.ReleaseDate)
		if err != nil {
			return apperrors.ErrInvalidDateLayout
		}
	}
	return nil
}

func (f Filters) ToQueryParams() (string, []any) {
	var (
		query   = " WHERE deleted = false"
		values  = make([]any, 0)
		counter = 1
	)

	// Where clause
	if f.title != "" {
		query += fmt.Sprintf(" AND title LIKE $%d", counter)
		values = append(values, f.title)
		counter++
	}
	if f.group != "" {
		query += fmt.Sprintf(" AND band = $%d", counter)
		values = append(values, f.group)
		counter++
	}
	if f.text != "" {
		query += fmt.Sprintf(" AND text LIKE $%d", counter)
		values = append(values, f.text)
		counter++
	}
	if f.link != "" {
		query += fmt.Sprintf(" AND link LIKE $%d", counter)
		values = append(values, f.link)
		counter++
	}

	switch f.dateMode {
	case exact:
		query += fmt.Sprintf(" AND release_date = $%d", counter)
		values = append(values, f.dateFrom)
	case after:
		query += fmt.Sprintf(" AND release_date > $%d", counter)
		values = append(values, f.dateFrom)
	case before:
		query += fmt.Sprintf(" AND release_date < $%d", counter)
		values = append(values, f.dateFrom)
	case between:
		query += fmt.Sprintf(" AND release_date BETWEEN $%d AND $%d", counter, counter+1)
		values = append(values, f.dateFrom, f.dateTo)
	}

	// Ordering, Limit and offset
	query += fmt.Sprintf(" ORDER BY %s ASC", orderingField)
	if f.limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", f.limit)
	}
	if f.offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", f.offset)
	}

	return query, values
}

func (f Filters) Offset() int {
	return f.offset
}
