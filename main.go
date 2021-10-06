package main

import (
	//"github.com/tclohm/project-pizza/models"
	"net/http"
	"fmt"
	"io/ioutil"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form, 10 << 20 specifies a maximum
	// upload of 10MB files
	// retrive file from posted form-data
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// write temp file within our pizza-image directory that follows
	tempFile, err := ioutil.TempFile("pizza-images", "upload-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	tempFile.Write(fileBytes)
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func setRoutes() {
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8080", nil)
}

func main() {
	setRoutes()
}