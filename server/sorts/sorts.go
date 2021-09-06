package sorts

import (
	"sort"

	"github.com/mplewis/rglimpse/types"
)

type ByDownRate []types.Stat

func (x ByDownRate) Len() int {
	return len(x)
}

func (x ByDownRate) Swap(i int, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x ByDownRate) Less(i int, j int) bool {
	return x[i].Status.DownRate < x[j].Status.DownRate
}

type ByAdded []types.Stat

func (x ByAdded) Len() int {
	return len(x)
}

func (x ByAdded) Swap(i int, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x ByAdded) Less(i int, j int) bool {
	return x[i].Torrent.Started.Before(x[j].Torrent.Started)
}

type ByFinished []types.Stat

func (x ByFinished) Len() int {
	return len(x)
}

func (x ByFinished) Swap(i int, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x ByFinished) Less(i int, j int) bool {
	return x[i].Torrent.Finished.Before(x[j].Torrent.Finished)
}

type SortAttr int

const (
	DownRate SortAttr = iota
	UpRate
	TimeRemaining
	Size
	Added
	Finished
)

type SortDir int

const (
	Ascending SortDir = iota
	Descending
)

func SortStatsByAdded(torrs []types.Stat, dir SortDir) []types.Stat {
	// Copy items from torrs to a new slice named sorted
	sorted := make([]types.Stat, len(torrs))
	copy(sorted, torrs)
	if dir == Ascending {
		sort.Sort(ByFinished(sorted))
	} else {
		sort.Sort(sort.Reverse(ByFinished(sorted)))
	}
	return sorted
}
