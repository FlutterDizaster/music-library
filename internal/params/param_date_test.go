package params

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newDateParameter(t *testing.T) {
	tests := []struct {
		name   string
		testFn func(t *testing.T)
	}{
		{
			name:   "Exact date",
			testFn: Test_newDateParameter_Exact,
		},
		{
			name:   "After date",
			testFn: Test_newDateParameter_After,
		},
		{
			name:   "Before date",
			testFn: Test_newDateParameter_Before,
		},
		{
			name:   "Between dates ordered",
			testFn: Test_newDateParameter_Between_Ordered,
		},
		{
			name:   "Between dates disordered",
			testFn: Test_newDateParameter_Between_Disordered,
		},
		{
			name:   "Invalid date value",
			testFn: Test_newDateParameter_InvalidDateValue,
		},
		{
			name:   "Invalid date range",
			testFn: Test_newDateParameter_InvalidDateRange,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.testFn)
	}
}

func Test_newDateParameter_Exact(t *testing.T) {
	value := "20.12.2022"
	field := "date"

	p, err := newDateParameter(value, field)
	require.NoError(t, err)

	dateP, ok := p.(*dateParameter)
	require.True(t, ok)

	assert.Equal(t, field, dateP.field)
	assert.Equal(t, value, dateP.value.Format("02.01.2006"))
	assert.Equal(t, exact, dateP.mode)
}

func Test_newDateParameter_After(t *testing.T) {
	value := ">20.12.2022"
	field := "date"

	p, err := newDateParameter(value, field)
	require.NoError(t, err)

	dateP, ok := p.(*dateParameter)
	require.True(t, ok)

	assert.Equal(t, field, dateP.field)
	assert.Equal(t, strings.TrimPrefix(value, ">"), dateP.value.Format("02.01.2006"))
	assert.Equal(t, after, dateP.mode)
}

func Test_newDateParameter_Before(t *testing.T) {
	value := "<20.12.2022"
	field := "date"

	p, err := newDateParameter(value, field)
	require.NoError(t, err)

	dateP, ok := p.(*dateParameter)
	require.True(t, ok)

	assert.Equal(t, field, dateP.field)
	assert.Equal(t, strings.TrimPrefix(value, "<"), dateP.value.Format("02.01.2006"))
	assert.Equal(t, before, dateP.mode)
}

func Test_newDateParameter_Between_Ordered(t *testing.T) {
	value := "20.12.2022-25.12.2022"
	field := "date"

	p, err := newDateParameter(value, field)
	require.NoError(t, err)

	dateP, ok := p.(*dateParameter)
	require.True(t, ok)

	resultValues := dateP.value.Format("02.01.2006") + "-" + dateP.to.Format("02.01.2006")

	assert.Equal(t, field, dateP.field)
	assert.Equal(t, value, resultValues)
	assert.Equal(t, between, dateP.mode)
}

func Test_newDateParameter_Between_Disordered(t *testing.T) {
	value := "25.12.2022-20.12.2022"
	field := "date"

	p, err := newDateParameter(value, field)
	require.NoError(t, err)

	dateP, ok := p.(*dateParameter)
	require.True(t, ok)

	resultValues := dateP.to.Format("02.01.2006") + "-" + dateP.value.Format("02.01.2006")

	assert.Equal(t, field, dateP.field)
	assert.Equal(t, value, resultValues)
	assert.Equal(t, between, dateP.mode)
}

func Test_newDateParameter_InvalidDateValue(t *testing.T) {
	value := "2022.12.20"
	field := "date"

	p, err := newDateParameter(value, field)
	require.ErrorIs(t, err, ErrInvalidDateValue)

	assert.Nil(t, p)
}

func Test_newDateParameter_InvalidDateRange(t *testing.T) {
	value := "2022-12-20-2022-12-20"
	field := "date"

	p, err := newDateParameter(value, field)
	require.ErrorIs(t, err, ErrInvalidDateRange)

	assert.Nil(t, p)
}

func Test_dateParameter_ToSQL(t *testing.T) {
	tests := []struct {
		name   string
		testFn func(t *testing.T)
	}{
		{
			name:   "Exact date",
			testFn: Test_dateParameter_ToSQL_Exact,
		},
		{
			name:   "After date",
			testFn: Test_dateParameter_ToSQL_After,
		},
		{
			name:   "Before date",
			testFn: Test_dateParameter_ToSQL_Before,
		},
		{
			name:   "Between dates",
			testFn: Test_dateParameter_ToSQL_Between,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.testFn)
	}
}

func Test_dateParameter_ToSQL_Exact(t *testing.T) {
	p := &dateParameter{
		field: "date",
		mode:  exact,
	}

	assert.Equal(t, " AND date = $1", p.ToSQL(1))
}

func Test_dateParameter_ToSQL_After(t *testing.T) {
	p := &dateParameter{
		field: "date",
		mode:  after,
	}

	assert.Equal(t, " AND date > $1", p.ToSQL(1))
}

func Test_dateParameter_ToSQL_Before(t *testing.T) {
	p := &dateParameter{
		field: "date",
		mode:  before,
	}

	assert.Equal(t, " AND date < $1", p.ToSQL(1))
}

func Test_dateParameter_ToSQL_Between(t *testing.T) {
	p := &dateParameter{
		field: "date",
		mode:  between,
	}

	assert.Equal(t, " AND date BETWEEN $1 AND $2", p.ToSQL(2))
}

func Test_dateParameter_Values(t *testing.T) {
	tests := []struct {
		name   string
		testFn func(t *testing.T)
	}{
		{
			name:   "One value",
			testFn: Test_dateParameter_Values_OneValue,
		},
		{
			name:   "Between values",
			testFn: Test_dateParameter_Values_Between,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.testFn)
	}
}

func Test_dateParameter_Values_OneValue(t *testing.T) {
	pTime, err := time.Parse("02.01.2006", "20.12.2022")
	require.NoError(t, err)

	p := &dateParameter{
		mode:  exact,
		value: pTime,
	}

	assert.Equal(t, []any{pTime}, p.Values())
}

func Test_dateParameter_Values_Between(t *testing.T) {
	pTime, err := time.Parse("02.01.2006", "20.12.2022")
	require.NoError(t, err)

	pTo, err := time.Parse("02.01.2006", "21.12.2022")
	require.NoError(t, err)

	p := &dateParameter{
		mode:  between,
		to:    pTo,
		value: pTime,
	}

	assert.Equal(t, []any{pTime, pTo}, p.Values())
}
