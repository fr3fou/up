package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
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
	// auth is the username:password combination for uploading files
	auth string

	// address is the port or address for up
	address string

	// dir is the directory for uploading files
	dir string

	// static is the handler for static files
	static http.Handler
)

var (
	DB *leveldb.DB
)

func main() {
	auth = env("AUTH", "")
	address = env("ADDRESS", ":8080")
	dir = env("DIR", "files/")

	static = http.StripPrefix("/",
		http.FileServer(
			http.Dir(dir),
		),
	)

	os.Mkdir(dir, 0777)

	var err error

	DB, err = leveldb.OpenFile("./db", nil)
	if err != nil {
		panic(err)
	}

	defer DB.Close()

	http.HandleFunc("/", rootHandler)

	log.Printf("up! ⚡️ is running on %s!", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatal(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method

	// Landing "page"
	if path == "/" && method == "GET" {
		landingPage(w, r)
		return
	}

	// Static file serving
	if path != "/" && method == "GET" {
		static.ServeHTTP(w, r)
		return
	}

	// Disallow any other methods
	if method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}

	// Upload the file
	uploadHandler(w, r)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Handle auth
	if !isAuth(w, r) {
		return
	}

	// Check for file size
	r.ParseMultipartForm(MaxSize)
	f, header, err := r.FormFile("file")
	if err != nil || header.Size > MaxSize {
		fmt.Fprintf(w, "Max file size is 512 MiB")
		return
	}

	// Read the file
	file, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	name, err := uploadFile(file, filepath.Ext(header.Filename))

	fmt.Fprintf(w, "https://%s/%s", r.Host, name)
}

func uploadFile(file []byte, extension string) (string, error) {
	hash := sha256.Sum256(file)
	val, err := DB.Get(hash[:], nil)

	// If the file has been found
	if val != nil && err != leveldb.ErrNotFound {
		file, err := os.Stat(dir + string(val))
		if err != nil {
			log.Println(err)
			return "", err
		}

		// time passed since upload
		timePassed, maxAge := calculateAge(file.ModTime(), file.Size())

		// only reupload if it's >= 95% through its lifespan
		if timePassed/maxAge < 0.95 {
			return string(val), nil
		}

		// Delete it, as it's going to be put anyway
		DB.Delete(hash[:], nil)
		os.Remove(dir + string(val))
	}

	var name string
	for {
		name = generateFileName(5) + extension

		// check if a file exists with the same name
		_, err := os.Stat("files/" + name)
		if err != nil {
			// if the file doesn't exist, we've found a name
			if os.IsNotExist(err) {
				break
			}

			// return if there is any other err
			return "", err
		}
	}

	// Create the file
	f, err := os.Create(dir + name)
	if err != nil {
		return "", err
	}

	f.Write(file)
	f.Close()

	return name, DB.Put(hash[:], []byte(name), nil)
}

func deleteFiles() {
	// Clear files every day
	for range time.Tick(time.Hour * 24) {
		files, err := ioutil.ReadDir(dir)

		if err != nil {
			fmt.Println(err)
			return
		}

		for _, file := range files {
			timePassed, maxAge := calculateAge(file.ModTime(), file.Size())

			if timePassed > maxAge {
				os.Remove(dir + file.Name())
				log.Printf("Deleted %s \n'", file.Name())
			}
		}
	}
}

// calculateAge returns the maximum possible age of a file and the time that has passed since its creation
func calculateAge(mod time.Time, size int64) (float64, float64) {
	// https://0x0.st
	// retention = min_age + (-max_age + min_age) * pow((file_size / max_size - 1), 3)
	return math.Floor(float64(MinAge) + float64(-MaxAge+MinAge)*math.Pow(float64(size)/float64(MaxSize)-float64(1), float64(3))),
		float64(daysBetween(mod, time.Now()))
}

// generateFileName takes in a length and generates a random name
func generateFileName(n int) string {
	// https://medium.com/@kpbird/golang-generate-fixed-size-random-string-dd6dbd5e63c0
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")
	b := make([]rune, n)

	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

func landingPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `%s!
UPLOAD:
	~/ $: curl -F 'file=@your-file' %s
	 %s/fpFx9.png
AUTH:
	Depending on the config of up, you may have to provide a Basic Authorization header
	~/ $: curl -F 'file=@your-file' %s --user username:password 
SOURCE:
	https://github.com/fr3fou/up	
`, r.Host, r.Host, r.Host, r.Host)
}

// https://gist.github.com/nicerobot/4375261#file-server-go
func isAuth(w http.ResponseWriter, r *http.Request) bool {
	// if up doesn't have an auth variable, anyone can upload
	if auth == "" {
		return true
	}
	cred := r.Header.Get("Authorization")

	if !strings.HasPrefix(cred, "Basic ") {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return false
	}

	up, err := base64.StdEncoding.DecodeString(cred[6:])

	if err != nil || string(up) != auth {
		log.Printf("Someone tried accessing with credentials :%s", string(up))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return false
	}

	return true
}

func env(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
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
