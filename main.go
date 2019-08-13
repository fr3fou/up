package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"go.etcd.io/bbolt"
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

var db *bbolt.DB

func main() {
	http.HandleFunc("/", rootHandler)

	fmt.Println("up! ⚡ is running on port :8080!")

	var err error

	db, err = bbolt.Open("files.db", 0600, nil)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	go func() {
		// Clear files every day
		for range time.Tick(time.Hour * 24) {
			files, err := ioutil.ReadDir("files/")

			if err != nil {
				fmt.Println(err)
				return
			}

			for _, file := range files {
				timePassed := float64(daysBetween(file.ModTime(), time.Now()))
				maxAge := math.Floor(calculateAge(MinAge, MaxAge, file.Size(), MaxSize))

				if timePassed > maxAge {
					os.Remove("files/" + file.Name())
					fmt.Println("Deleted " + file.Name())
				}
			}
		}
	}()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// https://0x0.st
// retention = min_age + (-max_age + min_age) * pow((file_size / max_size - 1), 3)
func calculateAge(minAge, maxAge int, fileSize, maxSize int64) float64 {
	return float64(MinAge) + float64(-MaxAge+MinAge)*math.Pow(float64(fileSize)/float64(MaxSize)-float64(1), float64(3))
}

func daysBetween(a, b time.Time) int {
	if a.After(b) {
		a, b = b, a

	}

	days := -a.YearDay()
	for year := a.Year(); year < b.Year(); year++ {
		days += time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC).YearDay()

	}
	days += b.YearDay()

	return days

}
