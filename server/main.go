package main

import (
	"log"
	"sync"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/mrobinsn/go-rtorrent/rtorrent"
)

type Page struct {
	Total    int         `json:"total"`
	Torrents interface{} `json:"torrents"`
}

type Stat struct {
	Torrent rtorrent.Torrent
	Status  rtorrent.Status
}

type SubscriptionArgs struct {
	Connection      *rtorrent.RTorrent
	Concurrency     int
	RefreshInterval time.Duration
}

func every(duration time.Duration, f func()) {
	go func() {
		f()
		for range time.Tick(duration) {
			f()
		}
	}()
}

func Subscribe(args SubscriptionArgs) (<-chan []Stat, <-chan error) {
	ch := make(chan []Stat)
	er := make(chan error)

	conn := args.Connection
	wp := workerpool.New(args.Concurrency)

	every(args.RefreshInterval, func() {
		start := time.Now()
		torrs, err := conn.GetTorrents(rtorrent.ViewMain)
		if err != nil {
			er <- err
			return
		}

		wg := sync.WaitGroup{}
		stats := []Stat{}
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
		sorted := SortStats(stats)
		ch <- sorted
		log.Printf("Updated %d torrents in %0.2f secs\n", len(torrs), duration.Seconds())
	})

	return ch, er
}

func main() {
	conn := rtorrent.New("http://admin:admin@wintermute:9080/RPC2", false)
	chStats, chErrs := Subscribe(SubscriptionArgs{
		Connection:      conn,
		Concurrency:     32,
		RefreshInterval: time.Second * 15,
	})

	log.Println("Subscribed to rTorrent")
	go func() { Serve(chStats) }()

	for err := range chErrs {
		log.Println(err)
	}
}
