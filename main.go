package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/mrobinsn/go-rtorrent/rtorrent"
)

type Stat struct {
	Torrent rtorrent.Torrent
	Status  rtorrent.Status
}

type ByDownRate []Stat

func (x ByDownRate) Len() int {
	return len(x)
}

func (x ByDownRate) Swap(i int, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x ByDownRate) Less(i int, j int) bool {
	return x[i].Status.DownRate < x[j].Status.DownRate
}

type ByAdded []Stat

func (x ByAdded) Len() int {
	return len(x)
}

func (x ByAdded) Swap(i int, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x ByAdded) Less(i int, j int) bool {
	return x[i].Torrent.Started.Before(x[j].Torrent.Started)
}

type ByFinished []Stat

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

func sortStatsByAdded(torrs []Stat, dir SortDir) []Stat {
	// Copy items from torrs to a new slice named sorted
	sorted := make([]Stat, len(torrs))
	copy(sorted, torrs)
	if dir == Ascending {
		sort.Sort(ByFinished(sorted))
	} else {
		sort.Sort(sort.Reverse(ByFinished(sorted)))
	}
	return sorted
}

func main() {
	// conn := rtorrent.New("http://admin:admin@localhost:9080/RPC2", false)
	// name, err := conn.Name()
	// if err != nil {
	// 	log.Panic(err)
	// }
	// log.Printf("My rTorrent's name: %v", name)

	// file, err := os.Create("stats.dat")
	// if err != nil {
	// 	log.Panic(err)
	// }
	// defer file.Close()
	// encoder := gob.NewEncoder(file)

	// torrs, err := conn.GetTorrents(rtorrent.ViewMain)
	// if err != nil {
	// 	log.Panic(err)
	// }

	// wp := workerpool.New(32)
	// done := make(chan Stat)
	// remaining := len(torrs)
	// for _, torr := range torrs {
	// 	torr := torr
	// 	fmt.Println(torr.Pretty())
	// 	wp.Submit(func() {
	// 		stat, err := conn.GetStatus(torr)
	// 		if err != nil {
	// 			log.Panic(err)
	// 		}
	// 		done <- Stat{torr, stat}
	// 		remaining -= 1
	// 		fmt.Println(remaining)
	// 		if remaining == 0 {
	// 			close(done)
	// 		}
	// 	})
	// }

	// stats := []Stat{}
	// for result := range done {
	// 	stats = append(stats, result)
	// 	fmt.Println(result.Torrent.Pretty())
	// }

	// encoder.Encode(stats)
	// fmt.Println(len(stats))

	file, err := os.Open("stats.dat")
	if err != nil {
		log.Panic(err)
	}
	decoder := gob.NewDecoder(file)
	var stats []Stat
	decoder.Decode(&stats)

	newest := sortStatsByAdded(stats, Descending)
	fmt.Println(len(newest))
	for _, stat := range newest[:10] {
		fmt.Println(stat.Torrent.Pretty())
		fmt.Println(stat.Torrent.Started)
	}
}
