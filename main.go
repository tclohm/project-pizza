package main

import (
	//"github.com/tclohm/project-pizza/models"
	"net/http"
	"fmt"
	"encoding/json"
	"log"
	"io/ioutil"
	"github.com/gorilla/mux"
)

func Up(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("api endpoint hit: return 'Up'")
	json.NewEncoder(w).Encode("Up")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// Parse multipart form, 10 << 20 specifies a maximum
	// upload of 10MB files
	// retrive file from posted form-data
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		json.NewEncoder(w).Encode("Error reading buffer")
		return
	}

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		json.NewEncoder(w).Encode("Error reading buffer")
		return
	}

	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" && filetype != "image/png" {
		json.NewEncoder(w).Encode("File format is not allowed. Please upload a JPEG or PNG image")
		return
	}

	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// write temp file within our pizza-image directory that follows
	tmpFile, err := ioutil.TempFile("uploads", "upload-*.png")
	if err != nil {
		json.NewEncoder(w).Encode("Error writing file")
	}
	defer tmpFile.Close()

	// read all of the contents of our uploaded file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		json.NewEncoder(w).Encode("Error reading contents")
		return
	}

	tmpFile.Write(fileBytes)

	fmt.Println("Successfully Uploaded File\n")

	json.NewEncoder(w).Encode("image uploaded")
}

func handleRequests() {
	router := mux.NewRouter()
	router.HandleFunc("/api", Up)
	router.HandleFunc("/upload", uploadHandler)
	
	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	fmt.Println("server running...")
	handleRequests()
}