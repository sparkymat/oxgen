package generator

type Config struct {
	TemplatesFolder string   `json:"templatesFolder"`
	ProjectID       string   `json:"projectId"`
	RepoURL         string   `json:"repoUrl"`
	PreCommands     []string `json:"preCommands"`
	PostCommands    []string `json:"postCommands"`
}
