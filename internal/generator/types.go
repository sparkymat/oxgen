package generator

import (
	"errors"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
)

var ErrInvalidResourceField = errors.New("invalid resource field")

type TemplateName string

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

type TemplateInput struct {
	Resource    TemplateName
	SearchField string
	HasSearch   bool
	Fields      []TemplateInputField
}

type TemplateInputField struct {
	Resource   TemplateName
	Field      TemplateName
	Type       string
	Modifiers  string
	Default    string
	Updateable bool
}

func TemplateInputFromNameAndFields(name string, fields []Field, searchField string) TemplateInput {
	return TemplateInput{
		SearchField: searchField,
		HasSearch:   searchField != "",
		Resource:    TemplateName(name),
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
	FieldTypeDate       FieldType = "date"
	FieldTypeTimestamp  FieldType = "timestamp"
	FieldTypeUUID       FieldType = "uuid"
	FieldTypeReferences FieldType = "references"
	FieldTypeAttachment FieldType = "attachment"
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
	case FieldTypeAttachment:
		return "text"
	case FieldTypeDate:
		return "date"
	case FieldTypeTimestamp:
		return "timestamp"
	case FieldTypeUnknown:
		return "unknown"
	default:
		return "unknown"
	}
}

type Field struct {
	Resource   string
	Name       string
	FieldType  FieldType
	Modifiers  string
	Required   bool
	Default    *string
	Updateable bool
}

//nolint:funlen,revive,cyclop
func ParseField(resource string, fieldString string) (Field, error) {
	words := strings.Split(fieldString, ":")

	//nolint:gomnd
	if len(words) < 2 {
		return Field{}, ErrInvalidResourceField
	}

	name := words[0]
	fieldTypeString := words[1]
	modifiers := ""
	updateable := false

	var fieldType FieldType

	var defaultValue *string

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
	case "attachment":
		fieldType = FieldTypeAttachment
		updateable = true
	case "date":
		fieldType = FieldTypeDate
	case "timestamp":
		fieldType = FieldTypeTimestamp
	default:
		fieldType = FieldTypeUnknown
	}

	if fieldType == FieldTypeUnknown {
		return Field{}, ErrInvalidResourceField
	}

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
		case word == "updateable":
			updateable = true
		default:
			return Field{}, ErrInvalidResourceField
		}
	}

	return Field{
		Resource:   resource,
		Name:       name,
		Updateable: updateable,
		FieldType:  fieldType,
		Modifiers:  modifiers,
		Default:    defaultValue,
	}, nil
}

func (f Field) TemplateInputField() TemplateInputField {
	name := strcase.ToSnake(f.Name)

	switch f.FieldType {
	case FieldTypeReferences:
		name += "_id"
	case FieldTypeAttachment:
		name += "_path"
	case FieldTypeDate:
	case FieldTypeTimestamp:
	case FieldTypeUUID:
	case FieldTypeBool:
	case FieldTypeInt:
	case FieldTypeString:
	case FieldTypeUnknown:
	default:
	}

	return TemplateInputField{
		Resource:   TemplateName(f.Resource),
		Field:      TemplateName(name),
		Updateable: f.Updateable,
		Type:       f.FieldType.String(),
		Modifiers:  f.Modifiers,
		Default: func() string {
			if f.Default == nil {
				return ""
			}

			return *f.Default
		}(),
	}
}
