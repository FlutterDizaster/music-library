package params

import "fmt"

type textLikeParameter struct {
	field string
	value string
}

func newTextLikeParameter(value, field string) Parameter {
	return &textLikeParameter{
		value: value,
		field: field,
	}
}

func (p textLikeParameter) ToSQL(num int) string {
	return fmt.Sprintf(" AND %s LIKE $%d", p.field, num)
}

func (p textLikeParameter) Values() []any {
	return []any{p.value}
}
