package params

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParams_newParams(t *testing.T) {
	tests := []struct {
		name   string
		testFn func(t *testing.T)
	}{
		{
			name:   "Success",
			testFn: TestParams_newParams_Success,
		},
		{
			name:   "InvalidDateValue",
			testFn: TestParams_newParams_InvalidDateValue,
		},
		{
			name:   "InvalidDateRange",
			testFn: TestParams_newParams_InvalidDateRange,
		},
		{
			name:   "InvalidLimitValue",
			testFn: TestParams_newParams_InvalidLimitValue,
		},
		{
			name:   "InvalidOffsetValue",
			testFn: TestParams_newParams_InvalidOffsetValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.testFn)
	}
}

func TestParams_newParams_Success(t *testing.T) {
	paramsTable := map[string]ParamType{
		"limit":  TypeLimit,
		"offset": TypeOffset,
		"title":  TypeTextLike,
		"date":   TypeDate,
	}
	paramsMap := map[string][]string{
		"limit":  {"10"},
		"offset": {"20"},
		"title":  {"song"},
		"date":   {"20.12.2022"},
	}

	params, err := newParams(paramsMap, paramsTable)

	require.NoError(t, err)

	assert.Equal(t, 10, params.limit)
	assert.Equal(t, 20, params.offset)

	assert.Len(t, params.paramsList, 2)
}

func TestParams_newParams_InvalidDateValue(t *testing.T) {
	paramsTable := map[string]ParamType{
		"date": TypeDate,
	}
	paramsMap := map[string][]string{
		"date": {"2022.12.20"},
	}

	params, err := newParams(paramsMap, paramsTable)

	require.ErrorIs(t, err, ErrFailedToParseParameters)
	require.ErrorIs(t, err, ErrInvalidDateValue)

	assert.Nil(t, params)
}

func TestParams_newParams_InvalidDateRange(t *testing.T) {
	paramsTable := map[string]ParamType{
		"date": TypeDate,
	}
	paramsMap := map[string][]string{
		"date": {"2022-12-20 2022-12-21"},
	}

	params, err := newParams(paramsMap, paramsTable)

	require.ErrorIs(t, err, ErrFailedToParseParameters)
	require.ErrorIs(t, err, ErrInvalidDateRange)

	assert.Nil(t, params)
}

func TestParams_newParams_InvalidOffsetValue(t *testing.T) {
	paramsTable := map[string]ParamType{
		"offset": TypeOffset,
	}
	paramsMap := map[string][]string{
		"offset": {"f"},
	}

	params, err := newParams(paramsMap, paramsTable)

	require.ErrorIs(t, err, ErrFailedToParseParameters)
	require.ErrorIs(t, err, ErrInvalidOffsetValue)

	assert.Nil(t, params)
}

func TestParams_newParams_InvalidLimitValue(t *testing.T) {
	paramsTable := map[string]ParamType{
		"limit": TypeLimit,
	}
	paramsMap := map[string][]string{
		"limit": {"f"},
	}

	params, err := newParams(paramsMap, paramsTable)

	require.ErrorIs(t, err, ErrFailedToParseParameters)
	require.ErrorIs(t, err, ErrInvalidLimitValue)

	assert.Nil(t, params)
}

func TestParams_ToQuery(t *testing.T) {
	tests := []struct {
		name   string
		testFn func(t *testing.T)
	}{
		{
			name:   "Exact date",
			testFn: TestParams_ToQuery_OneDate,
		},
		{
			name:   "Between dates",
			testFn: TestParams_ToQuery_BetweenDates,
		},
		{
			name:   "Text like",
			testFn: TestParams_ToQuery_TextLike,
		},
		{
			name:   "Limit",
			testFn: TestParams_ToQuery_Limit,
		},
		{
			name:   "Offset",
			testFn: TestParams_ToQuery_Offset,
		},
		{
			name:   "All field types",
			testFn: TestParams_ToQuery_AllFieldTypes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFn)
	}
}

func TestParams_ToQuery_OneDate(t *testing.T) {
	paramsTable := map[string]ParamType{
		"date": TypeDate,
	}
	paramsMap := map[string][]string{
		"date": {"20.12.2022"},
	}

	params, err := newParams(paramsMap, paramsTable)
	require.NoError(t, err)

	query, values := params.ToQuery("")
	assert.Equal(t, " WHERE 1=1 AND date = $1 ORDERED BY song ASC;", query)

	assert.Len(t, values, 1)

	expTime, err := time.Parse("02.01.2006", "20.12.2022")
	require.NoError(t, err)
	assert.Equal(t, expTime, values[0])
}

func TestParams_ToQuery_BetweenDates(t *testing.T) {
	paramsTable := map[string]ParamType{
		"date": TypeDate,
	}
	paramsMap := map[string][]string{
		"date": {"20.12.2022-25.12.2022"},
	}

	params, err := newParams(paramsMap, paramsTable)
	require.NoError(t, err)

	query, values := params.ToQuery("")
	assert.Equal(t, " WHERE 1=1 AND date BETWEEN $1 AND $2 ORDERED BY song ASC;", query)

	assert.Len(t, values, 2)

	expTime, err := time.Parse("02.01.2006", "20.12.2022")
	require.NoError(t, err)

	expToTime, err := time.Parse("02.01.2006", "25.12.2022")
	require.NoError(t, err)

	assert.Equal(t, expTime, values[0])
	assert.Equal(t, expToTime, values[1])
}

func TestParams_ToQuery_TextLike(t *testing.T) {
	paramsTable := map[string]ParamType{
		"title": TypeTextLike,
	}
	paramsMap := map[string][]string{
		"title": {"song"},
	}

	params, err := newParams(paramsMap, paramsTable)
	require.NoError(t, err)

	query, values := params.ToQuery("")
	assert.Equal(t, " WHERE 1=1 AND title LIKE $1 ORDERED BY song ASC;", query)

	assert.Len(t, values, 1)

	assert.Equal(t, "song", values[0])
}

func TestParams_ToQuery_Limit(t *testing.T) {
	paramsTable := map[string]ParamType{
		"limit": TypeLimit,
	}
	paramsMap := map[string][]string{
		"limit": {"10"},
	}

	params, err := newParams(paramsMap, paramsTable)
	require.NoError(t, err)

	query, values := params.ToQuery("")
	assert.Equal(t, " WHERE 1=1 ORDERED BY song ASC LIMIT 10;", query)

	assert.Empty(t, values)
}

func TestParams_ToQuery_Offset(t *testing.T) {
	paramsTable := map[string]ParamType{
		"offset": TypeOffset,
	}
	paramsMap := map[string][]string{
		"offset": {"10"},
	}

	params, err := newParams(paramsMap, paramsTable)
	require.NoError(t, err)

	query, values := params.ToQuery("")
	assert.Equal(t, " WHERE 1=1 ORDERED BY song ASC OFFSET 10;", query)

	assert.Empty(t, values)
}

//nolint:lll // too long query
func TestParams_ToQuery_AllFieldTypes(t *testing.T) {
	paramsTable := map[string]ParamType{
		"limit":        TypeLimit,
		"offset":       TypeOffset,
		"title":        TypeTextLike,
		"date":         TypeDate,
		"date-between": TypeDate,
	}
	paramsMap := map[string][]string{
		"limit":        {"10"},
		"offset":       {"20"},
		"title":        {"song"},
		"date":         {"20.12.2022"},
		"date-between": {"20.12.2022-25.12.2022"},
	}

	params, err := newParams(paramsMap, paramsTable)
	require.NoError(t, err)

	query, values := params.ToQuery("")

	assert.Equal(
		t,
		" WHERE 1=1 AND title LIKE $1 AND date = $2 AND date-between BETWEEN $3 AND $4 ORDERED BY song ASC LIMIT 10 OFFSET 20;",
		query,
	)

	assert.Len(t, values, 4)

	expTime, err := time.Parse("02.01.2006", "20.12.2022")
	require.NoError(t, err)

	expToTime, err := time.Parse("02.01.2006", "25.12.2022")
	require.NoError(t, err)

	assert.Equal(t, "song", values[0])
	assert.Equal(t, expTime, values[1])
	assert.Equal(t, expTime, values[2])
	assert.Equal(t, expToTime, values[3])
}
