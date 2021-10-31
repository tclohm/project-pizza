package handlers

import (
	"net/http"
	"fmt"
	"encoding/json"
	"log"
	"io/ioutil"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"github.com/tclohm/project-pizza/models"
)


type DBClient struct {
	Db *gorm.DB
}

func Up(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("api endpoint hit: return 'Up'")

	response := make(map[string]string)
	response["status"] = "OK"
	response["message"] = "Up and Running"


	res, err := json.Marshal(response)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(res)
}

func (driver *DBClient) PostImage(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Content-Type", "application/json")
	// Parse multipart form, 10 << 20 specifies a maximum
	// upload of 10MB files
	// retrive file from posted form-data
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		json.NewEncoder(w).Encode("Error reading buffer")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" && filetype != "image/png" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// write temp file within our uploads directory that follows
	tmpFile, err := ioutil.TempFile("uploads", "upload-*.png")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer tmpFile.Close()

	// read all of the contents of our uploaded file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	tmpFile.Write(fileBytes)

	fmt.Println("Successfully Uploaded File\n")
	fmt.Println("Temporary file location", tmpFile.Name())

	img := models.Image{Filename: handler.Filename, Content_type: handler.Header["Content-Type"][0], Location: tmpFile.Name()}
	
	driver.Db.Create(&img)

	res, err := json.Marshal(img)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(res)

}

func (driver *DBClient) GetImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println(r)
	// fs := http.FileServer(http.Dir("/uploads"))
	// http.Handle("/uploads/", http.StripPrefix("/uploads/", fs))

}

func (driver *DBClient) PostVenue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	fmt.Println(r)

}

func HandleRequests(driver DBClient) {

	router := mux.NewRouter()
	router.HandleFunc("/api", Up).Methods("GET")

	router.HandleFunc("/upload/image", driver.PostImage).Methods("POST")
	router.HandleFunc("/image/{id}", driver.GetImage).Methods("GET")

	router.HandleFunc("/post/venue", driver.PostVenue).Methods("POST")
	server := &http.Server{
		Handler: router,
		Addr: "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}