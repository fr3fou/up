package main

import (
	"log"
	"net/http"
	"os"
)

const (
	// MinAge is the minimum amount of days before a file gets deleted
	MinAge int = 30

	// MaxAge is the maximum amount of days before a file gets deleted
	MaxAge int = 365

	// MaxSize it the maximum size of a file (in bytes)
	MaxSize int64 = 512 * MiB

	// MiB is the size of 1 Mebibyte
	MiB = 1 << 20
)

var (
	// Auth is the username:password combination for uploading files
	Auth string

	// Address is the port or address for up
	Address string
)

func main() {
	Auth = env("AUTH", "")
	Address = env("ADDRESS", ":8080")

	http.HandleFunc("/", rootHandler)

	log.Printf("up! ⚡ is running on %s!", Address)
	if err := http.ListenAndServe(Address, nil); err != nil {
		log.Fatal(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

}

func env(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
