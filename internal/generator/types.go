package generator

import (
	"errors"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
)

var ErrInvalidResourceField = errors.New("invalid resource field")

type TemplateInput struct {
	ResourceCamelcaseSingular string
	ResourceCamelcasePlural   string
	ResourceUnderscorePlural  string
	SearchField               string
	Fields                    []TemplateInputField
}

type TemplateInputField struct {
	Name      string
	Type      string
	Modifiers string
	Default   string
}

func TemplateInputFromNameAndFields(name string, fields []Field, searchField string) TemplateInput {
	return TemplateInput{
		SearchField:               searchField,
		ResourceUnderscorePlural:  pluralize.NewClient().Plural(strcase.ToSnake(name)),
		ResourceCamelcaseSingular: pluralize.NewClient().Singular(strcase.ToCamel(name)),
		ResourceCamelcasePlural:   pluralize.NewClient().Plural(strcase.ToCamel(name)),
		Fields: lo.Map(fields, func(f Field, _ int) TemplateInputField {
			return f.TemplateInputField()
		}),
	}
}

type FieldType string

const (
	FieldTypeString     FieldType = "string"
	FieldTypeInt        FieldType = "int"
	FieldTypeBool       FieldType = "bool"
	FieldTypeUUID       FieldType = "uuid"
	FieldTypeReferences FieldType = "references"
	FieldTypeUnknown    FieldType = "unknown"
)

func (f FieldType) String() string {
	switch f {
	case FieldTypeString:
		return "text"
	case FieldTypeInt:
		return "integer"
	case FieldTypeBool:
		return "bool"
	case FieldTypeUUID:
		return "uuid" //nolint:goconst
	case FieldTypeReferences:
		return "uuid"
	case FieldTypeUnknown:
		return "unknown"
	default:
		return "unknown"
	}
}

type Field struct {
	Name      string
	FieldType FieldType
	Modifiers string
	Required  bool
	Default   *string
}

//nolint:funlen,revive,cyclop
func ParseField(fieldString string) (Field, error) {
	words := strings.Split(fieldString, ":")

	//nolint:gomnd
	if len(words) < 2 {
		return Field{}, ErrInvalidResourceField
	}

	name := words[0]
	fieldTypeString := words[1]

	var fieldType FieldType

	switch fieldTypeString {
	case "string":
		fieldType = FieldTypeString
	case "int":
		fieldType = FieldTypeInt
	case "bool":
		fieldType = FieldTypeBool
	case "uuid":
		fieldType = FieldTypeUUID
	case "references":
		fieldType = FieldTypeReferences
	default:
		fieldType = FieldTypeUnknown
	}

	if fieldType == FieldTypeUnknown {
		return Field{}, ErrInvalidResourceField
	}

	modifiers := ""

	var defaultValue *string

	for _, word := range words[2:] {
		switch {
		case strings.HasPrefix(word, "default="):
			kvWords := strings.Split(word, "=")
			w := "DEFAULT "

			if fieldType == FieldTypeString {
				w += "'" + kvWords[1] + "'"
			} else {
				w += kvWords[1]
			}

			defaultValue = &w
		case strings.HasPrefix(word, "table="):
			if modifiers != "" {
				modifiers += " "
			}

			kvWords := strings.Split(word, "=")
			modifiers += "REFERENCES " + kvWords[1] + "(id)"
		case word == "unique":
			if modifiers != "" {
				modifiers += " "
			}

			modifiers += "UNIQUE"
		case word == "not_null":
			if modifiers != "" {
				modifiers += " "
			}

			modifiers += "NOT NULL"
		default:
			return Field{}, ErrInvalidResourceField
		}
	}

	return Field{
		Name:      name,
		FieldType: fieldType,
		Modifiers: modifiers,
		Default:   defaultValue,
	}, nil
}

func (f Field) TemplateInputField() TemplateInputField {
	name := strcase.ToSnake(f.Name)

	if f.FieldType == FieldTypeReferences {
		name += "_id"
	}

	return TemplateInputField{
		Name:      name,
		Type:      f.FieldType.String(),
		Modifiers: f.Modifiers,
		Default: func() string {
			if f.Default == nil {
				return ""
			}

			return *f.Default
		}(),
	}
}
