package cmd

import (
	"strings"
)

const ListSeparator string = ","

func parseStringToList(
	list, sep string,
	narmalizationFn func(string) string,
) []string {
	var sanitizedItems []string

	items := strings.Split(list, sep)
	for _, i := range items {
		sanitizedItems = append(sanitizedItems, narmalizationFn(i))
	}

	return sanitizedItems
}

func noSpaceLowerCase(s string) string {
	return strings.ToLower(
		strings.TrimSpace(s),
	)
}
