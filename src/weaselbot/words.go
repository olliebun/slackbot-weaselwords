package weaselbot

import (
	"bufio"
	"io"
	"strings"
)

type Words []string

// Case-insensitive test for whether this word is in the set.
func (w Words) Matches(s string) Words {
	found := make(Words, 0)
	s = strings.ToLower(s)
	for _, weaselword := range w {
		if strings.Contains(s, weaselword) {
			found = append(found, weaselword)
		}
	}

	return found
}

// Load line-delimited words from a stream.
func WordsFromReader(input io.Reader) (Words, error) {
	words := make(Words, 0)
	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	return words, scanner.Err()
}
