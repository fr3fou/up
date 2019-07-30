package main

import (
	"net/http"
	"text/template"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, struct{}{})
}

func uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	// if r.Method == "POST" {
	// 	fmt.Fprintf(w, "/up only accepts POST")
	// 	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	// 	if err := r.ParseMultipartForm(MB * 512); err != nil {
	// 		fmt.Fprintf(w, "Max file size is 512MB")
	// 	}

	// 	file, header, err := r.FormFile("file")

	// 	if err != nil {
	// 		fmt.Fprintf(w, err.Error())
	// 		return
	// 	}

	// 	fmt.Printf(header.Filename, header.Size)

	// 	bytes, err := ioutil.ReadAll(file)

	// 	if err != nil {
	// 		fmt.Fprintf(w, err.Error())
	// 		return
	// 	}

	// 	name, err := uploadImage(bytes)

	// 	if err != nil {
	// 		fmt.Fprintf(w, err.Error())
	// 		return
	// 	}

	// 	fmt.Fprintf(w, "")
	// } else if r.Method == "GET" {
	// 	// render template
	// }

}

// func uploadImage(bytes []byte) (string, error) {

// }
