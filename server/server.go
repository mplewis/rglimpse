package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func Clamp(n int, min int, max int) int {
	if min > max {
		min, max = max, min
	}
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}

// TODO: write tests
func Subset(stats []Stat, offset int, count int) []Stat {
	left := Clamp(offset, 0, len(stats))
	right := Clamp(offset+count, left, len(stats))
	return stats[left:right]
}

func Structure(stats []Stat) []map[string]interface{} {
	merged := []map[string]interface{}{}
	for _, stat := range stats {
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
	return merged
}

func Serve(newStats <-chan []Stat) {
	stats := []Stat{}
	go func() {
		for incoming := range newStats {
			stats = incoming
		}
	}()

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
		subset := Subset(stats, offset, count)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"total":    len(stats),
			"torrents": Structure(subset),
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
	log.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}