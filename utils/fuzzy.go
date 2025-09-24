package utils

import "github.com/agnivade/levenshtein"

func SuggestStations(input string, stations map[string][]string, max int) []string {
	suggestions := []string{}
	for name := range stations {
		dist := levenshtein.ComputeDistance(input, name)
		if dist <= 5 {
			suggestions = append(suggestions, name)
		}
	}
	if len(suggestions) > max {
		return suggestions[:max]
	}
	return suggestions
}
