package main

import (
	"strings"
)

func containsAllWords(title string, words []string) bool {
	title = strings.ToLower(title)
	for _, word := range words {
		if !strings.Contains(title, word) {
			return false
		}
	}
	return true
}

func Filter(stats []Stat, query string) []Stat {
	words := strings.Fields(strings.ToLower(query))
	var filtered []Stat

	for _, stat := range stats {
		name := strings.ToLower(stat.Torrent.Name)
		if containsAllWords(name, words) {
			filtered = append(filtered, stat)
		}
	}
	return filtered
}
