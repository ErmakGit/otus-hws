package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type Frequency struct {
	Word  string
	Count int
}

func Top10(input string) []string {
	strimedInput := strings.Trim(input, " ")
	resultSlice := make([]string, 0, 10)
	if strimedInput == "" {
		return resultSlice
	}

	wordList := strings.Fields(strimedInput)
	frequencyMap := make(map[string]Frequency, 0)
	for _, word := range wordList {
		frequency, ok := frequencyMap[word]
		if !ok {
			frequency = Frequency{
				Word:  word,
				Count: 1,
			}
		} else {
			frequency.Count++
		}
		frequencyMap[word] = frequency
	}

	frequencySlice := make([]Frequency, 0, len(frequencyMap))
	for _, wordFrequency := range frequencyMap {
		frequencySlice = append(frequencySlice, wordFrequency)
	}

	sort.Slice(frequencySlice, func(i, j int) bool {
		if frequencySlice[i].Count == frequencySlice[j].Count {
			return frequencySlice[i].Word < frequencySlice[j].Word
		}

		return frequencySlice[i].Count > frequencySlice[j].Count
	})

	for _, word := range frequencySlice[:10] {
		resultSlice = append(resultSlice, word.Word)
	}

	return resultSlice
}
