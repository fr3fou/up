package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"
)

const (
	// MB is the size of 1 Megabyte
	MB = 1 << 20
)

var static = http.StripPrefix("/", http.FileServer(http.Dir("files/")))

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// Landing "page"
	if r.URL.Path == "/" && r.Method == "GET" {
		fmt.Fprintf(w, `up.simo.sh!

UPLOAD:
	~/ $: curl -F 'file=@your-file' --user username:password up.simo.sh
	 up.simo.sh/fpFx9.png
`)
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

	if err := r.ParseMultipartForm(MB * 512); err != nil {
		fmt.Fprintf(w, "Max file size is 512MB")
	}

	file, header, err := r.FormFile("file")

	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	fmt.Printf(header.Filename, header.Size)

	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	name, err := UploadFile(bytes, time.Hour*24*30, filepath.Ext(header.Filename))

	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	fmt.Fprintf(w, "https://up.simo.sh/"+name)
}
