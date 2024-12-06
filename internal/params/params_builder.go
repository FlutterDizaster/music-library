package params

import "github.com/FlutterDizaster/music-library/internal/abstraction"

type Builder struct {
	paramsTable map[string]ParamType
}

func NewBuilder(paramsTable map[string]ParamType) *Builder {
	return &Builder{
		paramsTable: paramsTable,
	}
}

func (b *Builder) Build(paramsMap map[string][]string) (abstraction.Params, error) {
	return newParams(paramsMap, b.paramsTable)
}
