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

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form, 10 << 20 specifies a maximum
	// upload of 10MB files
	// retrive file from posted form-data
	fmt.Println("hello?")
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		log.Fatal("error!", err)
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// write temp file within our pizza-image directory that follows
	tempFile, err := ioutil.TempFile("foodImages", "upload-*.png")
	if err != nil {
		log.Fatal("error!", err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("error!", err)
	}

	tempFile.Write(fileBytes)

	fmt.Println("Successfully Uploaded File\n")

	json.NewEncoder(w).Encode("image uploaded")
}

func handleRequests() {
	router := mux.NewRouter()
	router.HandleFunc("/api", Up)
	router.HandleFunc("/upload", uploadFile)
	
	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	fmt.Println("server running...")
	handleRequests()
}