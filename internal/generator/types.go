package generator

type Config struct {
	TemplatesFolder string     `json:"templatesFolder"`
	ProjectID       string     `json:"projectId"`
	RepoURL         string     `json:"repoUrl"`
	Resources       []Resource `json:"resources"`
	PreCommands     []string   `json:"preCommands"`
	PostCommands    []string   `json:"postCommands"`
}

type Resource struct {
	Name        string  `json:"name"`
	ServiceName string  `json:"serviceName"`
	Fields      []Field `json:"fields"`
}

type Field struct {
	Name     string    `json:"name"`
	Type     FieldType `json:"type"`
	Nillable bool      `json:"nillable"`
	Unique   bool      `json:"unique"`
	Index    bool      `json:"index"`
	Default  string    `json:"default"`
}

type FieldType string

const (
	FieldString   FieldType = "string"
	FieldInt      FieldType = "int"
	FieldFloat    FieldType = "float"
	FieldBool     FieldType = "bool"
	FieldDate     FieldType = "date"
	FieldDateTime FieldType = "datetime"
)
