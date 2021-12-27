package cmd

import (
	"strings"
)

const ListSeparator string = ","

// A minor utility function used to convert a string into a list of strings by
// splitting the string by the proveded `sep`arator and then passing it to a
// normalizing function (to do some post-processing) before the element is
// finally appended to the list.
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

// As the name (hopefully) suggest, it downcases a string and gets rid of its
// leading and trailing spaces.
func noSpaceLowerCase(s string) string {
	return strings.ToLower(
		strings.TrimSpace(s),
	)
}
