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


const directory = "/Users/taylor/Desktop/web/web-projects/pizza-hunter/api/"

type DBClient struct {
	Db *gorm.DB
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
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
		return
	}

	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// write temp file within our uploads directory that follows
	tmpFile, err := ioutil.TempFile("uploads", "upload-*.jpg")
	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode("Error")
		return
	}
	defer tmpFile.Close()

	// read all of the contents of our uploaded file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode("Error on read")
		return
	}

	tmpFile.Write(fileBytes)

	fmt.Println("Successfully Uploaded File\n")
	fmt.Println("Temporary file location", tmpFile.Name())

	img := models.Image{Filename: handler.Filename, Content_type: handler.Header["Content-Type"][0], Location: tmpFile.Name()}
	
	driver.Db.Create(&img)

	res, err := json.Marshal(img)
	if err != nil {
		json.NewEncoder(w).Encode("Error marshal image object")
		return
	}

	fmt.Println("image location saved to the db")

	w.Write(res)

}

func (driver *DBClient) PostVenue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Content-Type", "application/json")

	var venue models.Venue

	err := json.NewDecoder(r.Body).Decode(&venue)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	driver.Db.Create(&venue)

	res, err := json.Marshal(venue)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(res)
}

func (driver *DBClient) GetImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)

	var image models.Image

	driver.Db.First(&image, vars["id"])
	fmt.Println(directory + image.Location)

	http.ServeFile(w, r, directory + image.Location)

}

func (driver *DBClient) PostPizza(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Content-Type", "application/json")

	var pizza models.Pizza

	err := json.NewDecoder(r.Body).Decode(&pizza)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	driver.Db.Create(&pizza)

	res, err := json.Marshal(pizza)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(res)

}


func (driver *DBClient) PostVenuePizza(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Content-Type", "application/json")

	var venuePizza models.VenuePizza

	err := json.NewDecoder(r.Body).Decode(&venuePizza)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	driver.Db.Create(&venuePizza)

	res, err := json.Marshal(venuePizza)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(res)
}

func (driver *DBClient) GetMyPizzas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Content-Type", "application/json")
	var venuePizzas []models.VenuePizza

	driver.Db.Table("venue_pizzas").Select("*").Joins("left join pizzas on venue_pizzas.pizza_id = pizzas.id").Joins("left join venues on venue_pizzas.venue_id = venues.id")

	res, err := json.Marshal(venuePizzas)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(res)

}



func HandleRequests(driver DBClient) {

	router := mux.NewRouter()
	router.HandleFunc("/api", Up)
	router.Use(logging)

	router.HandleFunc("/upload/image", driver.PostImage)
	router.HandleFunc("/post/venue", driver.PostVenue)

	router.HandleFunc("/image/{id}", driver.GetImage)

	router.HandleFunc("/post/pizza", driver.PostPizza)
	router.HandleFunc("/post/venuepizza", driver.PostVenuePizza)

	router.HandleFunc("/get/mypizzas", driver.GetMyPizzas)

	server := &http.Server{
		Handler: router,
		Addr: "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}