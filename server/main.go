package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mplewis/rglimpse/sorts"
	"github.com/mplewis/rglimpse/types"
)

type Page struct {
	Total    int         `json:"total"`
	Torrents interface{} `json:"torrents"`
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
	// // encoder := gob.NewEncoder(file)

	// torrs, err := conn.GetTorrents(rtorrent.ViewMain)
	// if err != nil {
	// 	log.Panic(err)
	// }

	// wp := workerpool.New(32)
	// done := make(chan types.Stat)
	// remaining := len(torrs)
	// for _, torr := range torrs {
	// 	torr := torr
	// 	fmt.Println(torr.Pretty())
	// 	wp.Submit(func() {
	// 		stat, err := conn.GetStatus(torr)
	// 		if err != nil {
	// 			log.Panic(err)
	// 		}
	// 		done <- types.Stat{Torrent: torr, Status: stat}
	// 		remaining -= 1
	// 		fmt.Println(remaining)
	// 		if remaining == 0 {
	// 			close(done)
	// 		}
	// 	})
	// }

	// stats := []types.Stat{}
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
	var stats []types.Stat
	decoder.Decode(&stats)

	newest := sorts.SortStatsByAdded(stats, sorts.Descending)

	r := mux.NewRouter()
	r.HandleFunc("/torrents", func(w http.ResponseWriter, r *http.Request) {
		count, err := strconv.Atoi(r.URL.Query().Get("count"))
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		r.URL.Query().Get("")

		subset := newest[offset : offset+count]
		merged := []map[string]interface{}{}
		for _, stat := range subset {
			one := map[string]interface{}{
				"hash":            stat.Torrent.Hash,
				"name":            stat.Torrent.Name,
				"path":            stat.Torrent.Path,
				"size":            stat.Torrent.Size,
				"label":           stat.Torrent.Label,
				"completed":       stat.Torrent.Completed,
				"ratio":           stat.Torrent.Ratio,
				"created":         stat.Torrent.Created,
				"started":         stat.Torrent.Started,
				"finished":        stat.Torrent.Finished,
				"completed_bytes": stat.Status.CompletedBytes,
				"down_rate":       stat.Status.DownRate,
				"up_rate":         stat.Status.UpRate,
			}
			merged = append(merged, one)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"total":    len(newest),
			"torrents": merged,
		})
	})

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if port == "" {
		port = "9081"
	}

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
