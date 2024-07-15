package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var r = regexp.MustCompile(`[\wа-я]+(?:[.,:;!?@'"_-]+[\wа-я]+)*|(?:[-]{2,})`)

func countWords(words []string) map[string]int {
	result := make(map[string]int)
	for _, w := range words {
		result[w]++
	}

	return result
}

func sortWords(words map[string]int) []string {
	keys := make([]string, 0, len(words))
	for k := range words {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		if words[keys[i]] == words[keys[j]] {
			return keys[i] < keys[j]
		}
		return words[keys[i]] > words[keys[j]]
	})

	return keys
}

func Top10(str string) []string {
	strLower := strings.ToLower(str)
	words := r.FindAllString(strLower, -1)

	countedWords := countWords(words)

	sortedWords := sortWords(countedWords)

	wordsCount := len(sortedWords)
	var result []string
	if wordsCount > 10 {
		result = sortedWords[:10]
	} else {
		result = sortedWords[:wordsCount]
	}

	return result
}
