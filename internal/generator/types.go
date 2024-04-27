package generator

import (
	"errors"
	"strings"

	"github.com/iancoleman/strcase"
)

var ErrInvalidResourceField = errors.New("invalid resource field")

type Input struct {
	WorkspaceFolder string
	Service         TemplateName
	Resource        TemplateName
	SearchField     string
	HasSearch       bool
	Fields          []InputField
}

type InputField struct {
	Service    TemplateName
	Resource   TemplateName
	Name       TemplateName
	Type       FieldType
	Required   bool
	Default    string
	Table      string
	Unique     bool
	Updateable bool
	NotNull    bool
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

//nolint:funlen,revive,cyclop
func ParseField(service string, resource string, fieldString string) (InputField, error) {
	field := InputField{}

	words := strings.Split(fieldString, ":")

	//nolint:gomnd
	if len(words) < 2 {
		return InputField{}, ErrInvalidResourceField
	}

	field.Service = TemplateName(service)
	field.Resource = TemplateName(resource)
	field.Name = TemplateName(words[0])
	field.Updateable = false

	fieldTypeString := words[1]

	switch fieldTypeString {
	case "string":
		field.Type = FieldTypeString
	case "int":
		field.Type = FieldTypeInt
	case "bool":
		field.Type = FieldTypeBool
	case "uuid":
		field.Type = FieldTypeUUID
	case "references":
		field.Type = FieldTypeReferences
	case "attachment":
		field.Type = FieldTypeAttachment
		field.Updateable = true
	case "date":
		field.Type = FieldTypeDate
	case "timestamp":
		field.Type = FieldTypeTimestamp
	default:
		return InputField{}, ErrInvalidResourceField
	}

	for _, word := range words[2:] {
		switch {
		case strings.HasPrefix(word, "default="):
			kvWords := strings.Split(word, "=")
			field.Default = kvWords[1]
		case strings.HasPrefix(word, "table="):
			kvWords := strings.Split(word, "=")
			field.Table = kvWords[1]
		case word == "unique":
			field.Unique = true
		case word == "not_null":
			field.NotNull = true
		case word == "updateable":
			field.Updateable = true
		default:
			return InputField{}, ErrInvalidResourceField
		}
	}

	return field, nil
}

func (f FieldType) SQLType() string {
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

func (f FieldType) GoType() string {
	switch f {
	case FieldTypeString:
		return "string"
	case FieldTypeInt:
		return "int32"
	case FieldTypeBool:
		return "bool"
	case FieldTypeUUID:
		return "uuid.UUID" //nolint:goconst
	case FieldTypeReferences:
		return "uuid.UUID"
	case FieldTypeAttachment:
		return "string"
	case FieldTypeDate:
		return "time.Time"
	case FieldTypeTimestamp:
		return "time.Time"
	case FieldTypeUnknown:
		return "unknown"
	default:
		return "unknown"
	}
}

func (f FieldType) PresenterGoType() string {
	switch f {
	case FieldTypeString:
		return "string"
	case FieldTypeInt:
		return "int32"
	case FieldTypeBool:
		return "bool"
	case FieldTypeUUID:
		return "string"
	case FieldTypeReferences:
		return "string"
	case FieldTypeAttachment:
		return "string"
	case FieldTypeDate:
		return "string"
	case FieldTypeTimestamp:
		return "string"
	case FieldTypeUnknown:
		return "unknown"
	default:
		return "unknown"
	}
}

func (f InputField) Initial() bool {
	return f.NotNull || !f.Updateable
}

func (f InputField) JsonName() string {
	return strcase.ToLowerCamel(f.Name.String())
}

func (f InputField) CreateSQLFragment() string {
	fragment := "  " + f.Name.String() + " " + f.Type.SQLType()

	if f.Type == FieldTypeReferences {
		fragment += (" REFERENCES(" + f.Table + ")")
	}

	if f.Default != "" {
		if f.Type == FieldTypeString {
			fragment += (`"` + f.Default + `"`)
		} else {
			fragment += f.Default
		}
	}

	if f.NotNull {
		fragment += " NOT NULL"
	}

	if f.Unique {
		fragment += " UNIQUE"
	}

	return fragment
}

func (f InputField) JsonTag() string {
	return "`json:\"" + f.JsonName() + "\"`"
}

func (f InputField) CreateParamsGoFragment() string {
	fragment := "  " + f.Name.CamelcaseSingular() + " " + f.Type.GoType()

	return fragment
}

func (f InputField) CreateRequestGoFragment() string {
	fragment := "  " + f.Name.CamelcaseSingular() + " " + f.Type.GoType() + " " + f.JsonTag()

	return fragment
}

func (f InputField) PresenterGoFragment() string {
	fragment := "  " + f.Name.CamelcaseSingular() + " "

	if !f.NotNull {
		fragment += "*"
	}

	fragment += (f.Type.PresenterGoType() + " " + f.JsonTag())

	return fragment
}

func (f InputField) CreateAssignParamsGoFragment() string {
	fragment := "  " + f.Name.CamelcaseSingular() + ": " + f.Name.CamelcaseSingular()

	return fragment
}

func (f InputField) UpdateAssignParamGoFragment() string {
	if f.NotNull {
		return f.Name.String() + " = @" + f.Name.String() + "::" + f.Type.SQLType()
	}

	return f.Name.String() + " = sqlc.narg('" + f.Name.String() + "')"
}

func (f InputField) UpdateGoFunctionSignatureParam() string {
	paramString := "value "

	if !f.NotNull {
		paramString = "valuePtr *"
	}

	paramString += f.Type.GoType()

	return paramString
}

func (f InputField) PresenterAssignment() string {
	if f.Type == FieldTypeDate || f.Type == FieldTypeTimestamp {
		str := "if m." + f.Name.CamelcaseSingular() + ".Valid {\n"

		str += "v := m." + f.Name.CamelcaseSingular() + ".Time"

		switch f.Type {
		case FieldTypeDate:
			str += ".Format(\"2006-01-02\")"
		case FieldTypeTimestamp:
			str += ".Format(time.RFC3339)"
		default:
		}

		str += "\n"

		str += "item." + f.Name.CamelcaseSingular() + " = &v\n"
		str += "}\n\n"

		return str
	}

	str := ""

	if !f.NotNull {
		str += "\nif m." + f.Name.CamelcaseSingular() + ".Valid {\n"
	}

	str += "  item." + f.Name.CamelcaseSingular() + " = "

	if !f.NotNull {
		str += "&"
	}

	str += ("m." + f.Name.CamelcaseSingular())
	if !f.NotNull {
		str += ("." + f.PgType())
	}

	str += "\n"

	if !f.NotNull {
		str += "\n}\n\n"
	}

	return str
}

func (f InputField) PgType() string {
	switch f.Type {
	case FieldTypeString:
		return "String"
	case FieldTypeAttachment:
		return "String"
	case FieldTypeInt:
		return "Int32"
	case FieldTypeBool:
		return "Boolean"
	case FieldTypeTimestamp:
		return "Time"
	case FieldTypeDate:
		return "Time"
	}

	return "unknown"
}

func (f InputField) PgZeroValue() string {
	switch f.Type {
	case FieldTypeString:
		return "pgtype.Text{}"
	case FieldTypeInt:
		return "pgtype.Int4{}"
	case FieldTypeDate, FieldTypeTimestamp:
		return "pgtype.Date{}"
	}

	return "unknown"
}

func (f InputField) PgValue() string {
	switch f.Type {
	case FieldTypeString:
		return "pgtype.Text{String: *valuePtr, Valid: true}"
	case FieldTypeInt:
		return "pgtype.Int4{Int32: *valuePtr, Valid: true}"
	case FieldTypeDate, FieldTypeTimestamp:
		return "pgtype.Date{Time: *valuePtr, Valid: true}"
	}

	return "unknown"
}

func (f FieldType) TypescriptType() string {
	switch f {
	case FieldTypeString:
		return "string"
	case FieldTypeInt:
		return "number"
	case FieldTypeBool:
		return "boolean"
	case FieldTypeUUID:
		return "string"
	case FieldTypeReferences:
		return "string"
	case FieldTypeAttachment:
		return "string"
	case FieldTypeDate:
		return "string"
	case FieldTypeTimestamp:
		return "string"
	case FieldTypeUnknown:
		return "unknown"
	default:
		return "unknown"
	}
}

func (f InputField) FrontendModelDeclaration() string {
	str := "  public " + f.Name.LowerCamelcaseSingular()

	if !f.NotNull {
		str += "?"
	}

	str += ": " + f.Type.TypescriptType() + ";"

	return str
}

func (f InputField) FrontendModelAssignment() string {
	str := ""

	if !f.NotNull {
		str += "\n    if (json." + f.Name.LowerCamelcaseSingular() + ") {\n      "
	}

	str += "this." + f.Name.LowerCamelcaseSingular() + " = "

	if f.Type == FieldTypeDate || f.Type == FieldTypeTimestamp {
		str += "dayjs.utc("
	}

	str += "json." + f.Name.LowerCamelcaseSingular()

	if f.Type == FieldTypeDate || f.Type == FieldTypeTimestamp {
		str += ")"
	}

	str += ";"

	if !f.NotNull {
		str += "\n    }"
	}

	return str
}
