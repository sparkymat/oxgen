package generator

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

func (*Service) ReplacePlaceholders(ctx context.Context, templateString string, projectName string) string {
	idRegex := regexp.MustCompile(`[^a-z0-9]+`)
	projectIdentifier := idRegex.ReplaceAllString(strings.ToLower(projectName), "")

	replacementMap := map[string]string{
		"__PROJECT__": projectIdentifier,
		"__REPO__":    fmt.Sprintf("github.com/user/%s", projectIdentifier),
	}

	finalString := templateString

	for k, v := range replacementMap {
		finalString = strings.ReplaceAll(finalString, k, v)
	}

	return finalString
}
