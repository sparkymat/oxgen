package generator

import (
	"context"
	"strings"
)

func replacePlaceholders(_ context.Context, lookupTable map[string]string, templateString string) string {
	finalString := templateString

	for k, v := range lookupTable {
		if k != "" {
			finalString = strings.ReplaceAll(finalString, k, v)
		}
	}

	return finalString
}
