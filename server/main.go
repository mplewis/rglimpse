package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mrobinsn/go-rtorrent/rtorrent"
)

func GetEnvDefault(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}

func MustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatal("Missing required environment variable: " + key)
	}
	return val
}

func BoolEnv(key string, trueVal string, falseVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return falseVal
	}
	return trueVal
}

func IntEnv(key string, defaultValue int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(raw)
	if err != nil {
		log.Fatal(err)
	}
	return val
}

func main() {
	host := MustEnv("RTORRENT_HOST")
	port := GetEnvDefault("RTORRENT_PORT", "9080")
	username := MustEnv("RTORRENT_USERNAME")
	password := MustEnv("RTORRENT_PASSWORD")
	http := BoolEnv("RTORRENT_HTTPS", "https", "http")

	connString := fmt.Sprintf("%s://%s:%s@%s:%s/RPC2", http, username, password, host, port)
	censored := fmt.Sprintf("%s://%s:*****@%s:%s/RPC2", http, username, host, port)
	log.Printf("Connecting to %s\n", censored)

	conn := rtorrent.New(connString, false)
	// TODO: Show waiting message when torrents haven't yet loaded
	chStats, chErrs := Subscribe(SubscriptionArgs{
		Connection:      conn,
		Concurrency:     IntEnv("MAX_CLIENTS", 16),
		RefreshInterval: time.Second * time.Duration(IntEnv("REFRESH_INTERVAL_SECS", 30)),
	})

	go func() { Serve(conn, chStats) }()

	for err := range chErrs {
		log.Println(err)
	}
}
