package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/upload", uploadController)

	fmt.Println("App is running @ 4050")
	log.Fatalln(http.ListenAndServe(":4050", nil))
}

func uploadController(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, "uploading file")

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("picture")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}
	defer file.Close()

	fmt.Println(handler.Filename, handler.Header, handler.Size)

	tempFile, err := ioutil.TempFile("temp-images", "*.png")
	if err != nil {
		fmt.Fprintf(w, "something went wrong")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	if _, err := tempFile.Write(fileBytes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}

	fmt.Fprintln(w, "successfully uploaded")
}
