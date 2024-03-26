package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

const (
	topLimit    = 10
	dashSymbol  = "-"
	spaceSymbol = " "
	empty       = ""
)

type Frequency struct {
	Word  string
	Count int
}

func Top10(input string) []string {
	reg := regexp.MustCompile("[a-zA-zа-яА-я-,.]*[a-zA-zа-яА-я-]")
	resultSlice := make([]string, 0, topLimit)

	strimedInput := strings.Trim(input, spaceSymbol)
	if strimedInput == empty {
		return resultSlice
	}

	wordList := strings.Fields(strimedInput)
	frequencyMap := make(map[string]Frequency, 0)
	for _, word := range wordList {
		if word == dashSymbol {
			continue
		}

		regWord := reg.FindString(strings.ToLower(word))
		frequency, ok := frequencyMap[regWord]
		if !ok {
			frequency = Frequency{
				Word:  regWord,
				Count: 0,
			}
		}

		frequency.Count++
		frequencyMap[regWord] = frequency
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

	for _, word := range frequencySlice[:topLimit] {
		resultSlice = append(resultSlice, word.Word)
	}

	return resultSlice
}
