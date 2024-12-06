package params

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/FlutterDizaster/music-library/internal/apperrors"
)

type ParamType int

const (
	orderedClauseSize = 2

	orderingFields = "song"

	paramLimit  = "limit"
	paramOffset = "offset"

	TypeLimit ParamType = iota
	TypeOffset
	TypeTextLike
	TypeDate
)

type Parameter interface {
	ToSQL(num int) string
	Values() []any
}

type Params struct {
	paramsList []Parameter

	limit  int
	offset int
}

func newParams(paramsMap map[string][]string, paramsTable map[string]ParamType) (*Params, error) {
	params := &Params{
		paramsList: make([]Parameter, 0, len(paramsMap)),
	}
	var (
		param Parameter
		err   error
	)

	for name, value := range paramsMap {
		param = nil

		switch {
		case paramsTable[name] == TypeTextLike:
			param = newTextLikeParameter(value[0], name)

		case paramsTable[name] == TypeDate:
			param, err = newDateParameter(value[0], name)

		case name == paramLimit:
			params.limit, err = strconv.Atoi(value[0])
			if err != nil {
				err = errors.Join(ErrInvalidLimitValue, err)
			}

		case name == paramOffset:
			params.offset, err = strconv.Atoi(value[0])
			if err != nil {
				err = errors.Join(ErrInvalidOffsetValue, err)
			}

		default:
			err = apperrors.Error{
				Message: "Invalid parameter: " + name,
			}
		}

		if err != nil {
			return nil, errors.Join(ErrFailedToParseParameters, err)
		}

		if param != nil {
			params.paramsList = append(params.paramsList, param)
		}
	}

	return params, nil
}

// ToQuery creates a query string from the base query and slice of values.
func (p Params) ToQuery(baseQuery string) (string, []any) {
	var (
		values = make([]any, 0, len(p.paramsList))

		whereClauseParts = make([]string, 0, len(p.paramsList))
		whereClause      string

		orderedCouseParts = make([]string, 0, orderedClauseSize)
		orderedClause     string
	)

	// Prepare the where clause
	for i, param := range p.paramsList {
		pValues := param.Values()

		whereClauseParts = append(whereClauseParts, param.ToSQL(i+len(pValues)))
		values = append(values, pValues...)
	}

	whereClause = strings.Join(whereClauseParts, "")

	// Prepare the ordered clause
	if p.limit > 0 {
		orderedCouseParts = append(orderedCouseParts, fmt.Sprintf(" LIMIT %d", p.limit))
	}

	if p.offset > 0 {
		orderedCouseParts = append(orderedCouseParts, fmt.Sprintf(" OFFSET %d", p.offset))
	}

	orderedClause = strings.Join(orderedCouseParts, "")

	// Build the final query
	baseQuery = strings.TrimSuffix(baseQuery, ";")

	query := fmt.Sprintf(
		"%s WHERE 1=1%s ORDERED BY %s ASC%s;",
		baseQuery,
		whereClause,
		orderingFields,
		orderedClause,
	)

	return query, values
}

func (p Params) Limit() int {
	return p.limit
}

func (p Params) Offset() int {
	return p.offset
}
