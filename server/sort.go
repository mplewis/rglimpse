package main

import (
	"sort"
)

type ByIncompleteThenAddedDesc []Stat

func (x ByIncompleteThenAddedDesc) Len() int {
	return len(x)
}

func (x ByIncompleteThenAddedDesc) Swap(i int, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x ByIncompleteThenAddedDesc) Less(i int, j int) bool {
	if !x[i].Status.Completed && x[j].Status.Completed {
		return true
	}
	if x[i].Status.Completed && !x[j].Status.Completed {
		return false
	}
	return x[i].Torrent.Started.After(x[j].Torrent.Started)
}

// SortStats sorts Stats, incomplete before complete, then by most recent.
func SortStats(torrs []Stat) []Stat {
	sorted := make([]Stat, len(torrs))
	copy(sorted, torrs)
	sort.Sort(ByIncompleteThenAddedDesc(sorted))
	return sorted
}
