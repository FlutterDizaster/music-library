package abstraction

type Params interface {
	ToQuery(baseQuery string) (string, []any)
	Limit() int
	Offset() int
}

type ParamsBuilder interface {
	Build(paramsMap map[string][]string) (Params, error)
}
