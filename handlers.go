package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `up.simo.sh!

UPLOAD:
	~/ $: curl -F 'file=@your-file' --user username:password up.simo.sh
	up.simo.sh/fpFx9.png
`)
		return
	} else if r.Method == "POST" {
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

		name, err := UploadFile(bytes, time.Hour*24*30)

		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		fmt.Fprintf(w, name)
	}
}
