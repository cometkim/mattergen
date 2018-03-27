package main

import "fmt"

const (
	TargetTypeScript = "typescript"
	TargetFlow       = "flow"
	TargetGraphQL    = "graphql"
)

type Target struct {
	Id  string
	Ext string
}

func NewTarget(target string) (*Target, error) {
	switch target {
	case TargetTypeScript:
		return &Target{
			Id:  TargetTypeScript,
			Ext: ".d.ts",
		}, nil
	case TargetFlow:
		return &Target{
			Id:  TargetFlow,
			Ext: ".js.flow",
		}, nil
	}
	return nil, fmt.Errorf("%s is not supported", target)
}

func (target *Target) convertType(t string) string {
	switch target.Id {
	case TargetTypeScript:
		fallthrough
	case TargetFlow:
		return map[string]string{
			"bool":        "boolean",
			"int":         "number",
			"int32":       "number",
			"int64":       "number",
			"float":       "number",
			"float32":     "number",
			"float64":     "number",
			"string":      "string",
			"StringArray": "Array<string>",
			"StringMap":   "Map<string, string>",
		}[t]
	}
	return ""
}
