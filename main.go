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
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/upload", uploadImageHandler)

	fmt.Println("up is running on port :8080!")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
