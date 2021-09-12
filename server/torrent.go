package main

import (
	"log"
	"sync"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/mrobinsn/go-rtorrent/rtorrent"
)

// Stats combine the Torrent and Status structs for full detailed info.
type Stat struct {
	Torrent rtorrent.Torrent
	Status  rtorrent.Status
}

// SubscriptionArgs is the arguments for a new RTorrent subscription.
type SubscriptionArgs struct {
	Connection      *rtorrent.RTorrent
	Concurrency     int
	RefreshInterval time.Duration
}

// Every runs a function immediately, then at every following interval.
func Every(duration time.Duration, f func()) {
	go func() {
		go f()
		for range time.Tick(duration) {
			go f()
		}
	}()
}

// Subscribe connects to RTorrent and returns channels for periodic updates and any errors.
func Subscribe(args SubscriptionArgs) (<-chan []Stat, <-chan error) {
	ch := make(chan []Stat)
	er := make(chan error)

	conn := args.Connection
	wp := workerpool.New(args.Concurrency)

	Every(args.RefreshInterval, func() {
		start := time.Now()
		log.Println("Fetching torrents...")
		torrs, err := conn.GetTorrents(rtorrent.ViewMain)
		if err != nil {
			er <- err
			return
		}

		wg := sync.WaitGroup{}
		stats := []Stat{}
		log.Println("Fetching detailed torrent info...")
		for _, torr := range torrs {
			torr := torr
			wg.Add(1)
			wp.Submit(func() {
				st, err := conn.GetStatus(torr)
				if err != nil {
					er <- err
					return
				}
				stats = append(stats, Stat{Torrent: torr, Status: st})
				wg.Done()
			})
		}
		wg.Wait()

		duration := time.Since(start)
		log.Printf("Updated %d torrents in %0.2f secs\n", len(torrs), duration.Seconds())
		ch <- SortStats(stats)
	})

	return ch, er
}
