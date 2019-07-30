package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	// MB is the size of 1 Megabyte
	MB = 1 << 20
)

func main() {
	fs := http.FileServer(http.Dir("assets/"))

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/upload", uploadImageHandler)
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("up is running on port :8080!")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
