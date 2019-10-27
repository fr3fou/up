package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
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

func main() {
	auth = env("AUTH", "")
	address = env("ADDRESS", ":8080")
	dir = env("DIR", "files/")

	static = http.StripPrefix("/",
		http.FileServer(
			http.Dir(dir),
		),
	)

	http.HandleFunc("/", rootHandler)

	log.Printf("up! ⚡ is running on %s!", address)
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
	file, header, err := r.FormFile("file")
	if err != nil || header.Size > MaxSize {
		fmt.Fprintf(w, "Max file size is 512 MiB")
		return
	}

	_, err = ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	name := "test"

	// name, err := UploadFile(bytes, header.Size, filepath.Ext(header.Filename))

	// if err != nil {
	// 	fmt.Fprintf(w, err.Error())
	// }

	fmt.Fprintf(w, "https://%s/%s", r.Host, name)
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
