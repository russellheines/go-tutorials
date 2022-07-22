package main

import (
	"strings"

	"golang.org/x/tour/wc"
)

func WordCount(s string) map[string]int {
	m := make(map[string]int)

	for _, elem := range strings.Fields(s) {
		if _, ok := m[elem]; ok == false {
			m[elem] = 1
		} else {
			m[elem] = m[elem] + 1
		}
	}

	return m
}

func main() {
	wc.Test(WordCount)
}
