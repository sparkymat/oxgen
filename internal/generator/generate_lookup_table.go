package generator

import (
	pluralize "github.com/gertd/go-pluralize"
	"github.com/martinusso/inflect"
)

func generateLookupTableForProject(c Config) map[string]string {
	lookupTable := map[string]string{}

	// Project ID
	lookupTable["__PROJECT__"] = c.ProjectID

	// Repo
	lookupTable["__REPO__"] = c.RepoURL

	return lookupTable
}

func generateLookupTableForResource(c Config, resource string) map[string]string {
	lookupTable := generateLookupTableForProject(c)

	p := pluralize.NewClient()

	resourcePlural := p.Plural(resource)
	resourceSingular := p.Singular(resource)

	lookupTable["__RESOURCE_P_LC__"] = inflect.LowerCamelize(resourcePlural)
	lookupTable["__RESOURCE_S_LC__"] = inflect.LowerCamelize(resourceSingular)

	lookupTable["__RESOURCE_P_C__"] = inflect.Camelize(resourcePlural)
	lookupTable["__RESOURCE_S_C__"] = inflect.Camelize(resourceSingular)

	lookupTable["__RESOURCE_P_U__"] = inflect.Underscore(resourcePlural)
	lookupTable["__RESOURCE_S_U__"] = inflect.Underscore(resourceSingular)

	return lookupTable
}
