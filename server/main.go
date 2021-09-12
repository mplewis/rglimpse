package main

import (
	"log"
	"time"

	"github.com/mrobinsn/go-rtorrent/rtorrent"
)

func main() {
	// TODO: args for username, password, host, port, insecure
	conn := rtorrent.New("http://admin:admin@wintermute:9080/RPC2", false)
	// TODO: Show waiting message when torrents haven't yet loaded
	chStats, chErrs := Subscribe(SubscriptionArgs{
		Connection:      conn,
		Concurrency:     16,
		RefreshInterval: time.Second * 30,
	})

	log.Println("Subscribed to rTorrent")
	go func() { Serve(conn, chStats) }()

	for err := range chErrs {
		log.Println(err)
	}
}
