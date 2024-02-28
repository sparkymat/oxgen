package generator

func generateLookupTable(c Config) map[string]string {
	lookupTable := map[string]string{}

	// Project ID
	lookupTable["__PROJECT__"] = c.ProjectID

	// Repo
	lookupTable["__REPO__"] = c.RepoURL

	return lookupTable
}
