package generator

import (
	"errors"
	"strings"

	"github.com/iancoleman/strcase"
)

var (
	ErrInvalidResourceField = errors.New("invalid resource field")
)

type MigrationTemplateInput struct {
	ResourceUnderscorePlural string
	Fields                   []MigrationTemplateInputField
}

type MigrationTemplateInputField struct {
	Name      string
	Type      string
	Modifiers string
	Default   string
}

type FieldType string

const (
	FieldTypeString     FieldType = "string"
	FieldTypeInt        FieldType = "int"
	FieldTypeBool       FieldType = "bool"
	FieldTypeUUID       FieldType = "uuid"
	FieldTypeReferences FieldType = "references"
	FieldTypeUnknown    FieldType = "unkown"
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
		return "uuid"
	case FieldTypeReferences:
		return "uuid"
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

func ParseField(fieldString string) (Field, error) {
	words := strings.Split(fieldString, ":")

	if len(words) < 2 {
		return Field{}, ErrInvalidResourceField
	}

	name := words[0]
	fieldTypeString := words[1]
	fieldType := FieldTypeUnknown

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
		if strings.HasPrefix(word, "default=") {
			kvWords := strings.Split(word, "=")
			w := "DEFAULT " + kvWords[1]
			defaultValue = &w
		} else if strings.HasPrefix(word, "table=") {
			if modifiers != "" {
				modifiers += " "
			}

			kvWords := strings.Split(word, "=")
			modifiers += "REFERENCES " + kvWords[1] + "(id)"
		} else if word == "unique" {
			if modifiers != "" {
				modifiers += " "
			}

			modifiers += "UNIQUE"
		} else if word == "not_null" {
			if modifiers != "" {
				modifiers += " "
			}

			modifiers += "NOT NULL"
		} else {
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

func (f Field) MigrationTemplateInputField() MigrationTemplateInputField {
	name := strcase.ToSnake(f.Name)

	if f.FieldType == FieldTypeReferences {
		name += "_id"
	}

	return MigrationTemplateInputField{
		Name:      name,
		Type:      f.FieldType.String(),
		Modifiers: f.Modifiers,
		Default: func() string {
			if f.Default == nil {
				return ""
			} else {
				return *f.Default
			}
		}(),
	}
}
