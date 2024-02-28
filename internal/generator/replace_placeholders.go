package generator

import (
	"context"
	"strings"
)

func (s *Service) ReplacePlaceholders(_ context.Context, templateString string) string {
	finalString := templateString

	for k, v := range s.LookupTable {
		if k != "" {
			finalString = strings.ReplaceAll(finalString, k, v)
		}
	}

	return finalString
}
