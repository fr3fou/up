package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"go.etcd.io/bbolt"
)

var static = http.StripPrefix("/", http.FileServer(http.Dir("files/")))

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// Landing "page"
	if r.URL.Path == "/" && r.Method == "GET" {
		landingPage(w, r)
		return
	}

	// Static file serving
	if r.URL.Path != "/" && r.Method == "GET" {
		static.ServeHTTP(w, r)
		return
	}

	// Disallow any other methods
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}

	// Handle auth
	if !auth(w, r) {
		return
	}

	r.ParseMultipartForm(MaxSize)

	file, header, err := r.FormFile("file")

	if header.Size > MaxSize {
		fmt.Fprintf(w, "Max file size is 512 MiB")
		return
	}

	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	db, err := bbolt.Open("files.db", 0600, nil)

	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	defer db.Close()

	var bucket *bbolt.Bucket

	db.Update(func(tx *bbolt.Tx) (err error) {
		bucket, err = tx.CreateBucketIfNotExists([]byte("files"))

		if err != nil {
			return err
		}

		name, err := UploadFile(bytes, header.Size, filepath.Ext(header.Filename), bucket)

		if err != nil {
			fmt.Fprintf(w, err.Error())
			return err
		}

		fmt.Fprintf(w, r.Host+"/"+name)

		return nil
	})
}

// https://gist.github.com/nicerobot/4375261#file-server-go
func auth(w http.ResponseWriter, r *http.Request) bool {
	userPass := os.Getenv("AUTH")
	auth := r.Header.Get("Authorization")

	if !strings.HasPrefix(auth, "Basic ") {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return false
	}

	up, err := base64.StdEncoding.DecodeString(auth[6:])

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return false
	}

	if string(up) != userPass {
		log.Print("Someone tried accessing with credentials: ", string(up))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return false
	}

	return true
}

func landingPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `%s!

UPLOAD:
	~/ $: curl -F 'file=@your-file' --user username:password %s
	 %s/fpFx9.png

NOTE:
	Registrations are NOT open.

CONTACT:
	simo at deliriumproducts.me
`, r.Host, r.Host, r.Host)
}
