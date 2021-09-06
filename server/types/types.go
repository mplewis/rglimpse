package types

import "github.com/mrobinsn/go-rtorrent/rtorrent"

type Stat struct {
	Torrent rtorrent.Torrent
	Status  rtorrent.Status
}
