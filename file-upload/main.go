package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// file handler for route "/upload"
func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Uploading file...")
	r.ParseMultipartForm(10 << 20) // "10 << 20" means that the maximum file size is 10MB

	file, handler, err := r.FormFile("fileUploadForm")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Uploaded file: %+v\n", handler.Filename)
	fmt.Printf("File size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// make sure we close the file
	defer file.Close()

	tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	tempFile.Write(fileBytes)

	fmt.Println("File upload successfull")
}

func setupRoutes() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	// add route
	http.HandleFunc("/upload", uploadFile)
	fmt.Println("Server is up and running!")
	// host the server on port 8080
	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("Setting up file-upload server...")

	// init routes
	setupRoutes()
}
