package generator

import (
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

type TemplateName string

func (t TemplateName) LowerCamelcaseSingular() string {
	return pluralize.NewClient().Singular(strcase.ToLowerCamel(string(t)))
}

func (t TemplateName) LowerCamelcasePlural() string {
	return pluralize.NewClient().Plural(strcase.ToLowerCamel(string(t)))
}

func (t TemplateName) CamelcaseSingular() string {
	return pluralize.NewClient().Singular(strcase.ToCamel(string(t)))
}

func (t TemplateName) CamelcasePlural() string {
	return pluralize.NewClient().Plural(strcase.ToCamel(string(t)))
}

func (t TemplateName) UnderscoreSingular() string {
	return pluralize.NewClient().Singular(strcase.ToSnake(string(t)))
}

func (t TemplateName) UnderscorePlural() string {
	return pluralize.NewClient().Plural(strcase.ToSnake(string(t)))
}

func (t TemplateName) Upcase() string {
	return strings.ToUpper(string(t))
}

func (t TemplateName) Downcase() string {
	return strings.ToLower(string(t))
}

func (t TemplateName) Capitalize() string {
	return strings.ToUpper(string(t)[0:1]) + string(t)[1:]
}

func (t TemplateName) String() string {
	return string(t)
}
