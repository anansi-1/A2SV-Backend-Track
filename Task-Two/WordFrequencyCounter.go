package main

import (
	"strings"
	"unicode"
)

func WordFrequencyCounter(s string)map[string]int{

	 count := make(map[string]int)

	 var cleaned string
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			cleaned += string(r)
		}
	}

	cleaned = strings.ToLower(cleaned)
	words := strings.Fields(cleaned)

	for _,word := range words {
		count[word] += 1
	}

	return count

}


