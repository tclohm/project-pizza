package main

import (
	//"github.com/tclohm/project-pizza/models"
	"net/http"
	"fmt"
	"encoding/json"
	stdlog "log"
	"io/ioutil"
	"database/sql"
	"os"
	"runtime/debug"
	"time"

	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/go-kit/kit/log"

)

// MARK: -- middleware! 
// if we needed to pass values between handlers, such as
// the ID of the authenticated user, or a request or a trace
// ID, we can use context.Context attached to the *http.Request
// via the *Request.Context()
// func Middleware(next http.Handler) http.Handler {
// 	// wrap anon func, and cast it to a http.HandlerFunc
// 	// signature matches ServeHTTP(w,r)
// 	return http.HandlerFunc(
// 		func(w http.ResponseWriter, r *http.Request) {
// 			// logic before - reading request values,

// 			// call the 'next' handler in the chain
// 			next.ServeHTTP(w, r)
// 		}
// 	)
// }

// configuredRouter := LoggingMiddleware(OtherMiddleware(YetAnotherMiddleware(router)))
// log.Fatal(http.ListenAndServe(":8000", configuredRouter))

// func NewExampleMiddleware(something string) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		fn := func(w http.ResponseWriter, r *http.Request) {
// 			// Logic here

// 			// Call the next handler
// 			next.ServeHTTP(w, r))
// 		}

// 		return http.HandlerFunc(fn)
// 	}
// }

// Minimal wrapper for http.ResponseWriter that allows the 
// written HTTP status code to be captured for logging
type responseWriter struct {
	http.ResponseWriter
	status int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	return
}

// LoggingMiddleware logs the incoming HTTP request & its duration.
func LoggingMiddleware(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Log(
						"err", err,
						"trace", debug.Stack(),
					)
				}
			}()

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			logger.Log(
				"status", wrapped.status,
				"method", r.Method,
				"path", r.URL.EscapedPath(),
				"duration", time.Since(start),
			)
		}
		return http.HandlerFunc(fn)
	}
}


func Up(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("api endpoint hit: return 'Up'")
	json.NewEncoder(w).Encode(`{"Success": "true", "msg": "Up"}`)
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
		json.NewEncoder(w).Encode(`{ "Success": false, "msg": "File format is not allowed. Please upload a JPEG or PNG image"}`)
		return
	}

	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// write temp file within our pizza-image directory that follows
	tmpFile, err := ioutil.TempFile("uploads", "upload-*.png")
	if err != nil {
		json.NewEncoder(w).Encode(`{"Success": "false", "msg": "Error writing file"}`)
	}
	defer tmpFile.Close()

	// read all of the contents of our uploaded file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		json.NewEncoder(w).Encode(`{"Success": "false", "msg": "Error reading contents"}`)
		return
	}

	tmpFile.Write(fileBytes)

	fmt.Println("Successfully Uploaded File\n")

	json.NewEncoder(w).Encode(`"Success": "true"`)
}

func handleRequests() {
	router := mux.NewRouter()
	router.HandleFunc("/api", Up).Methods("GET")
	router.HandleFunc("/upload", uploadHandler).Methods("POST")
	
	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	err := godotenv.Load()
	pg_connection_string := fmt.Sprintf("port=%s host=%s user=%s "+ 
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("PSQL_PORT"), os.Getenv("HOSTNAME"), os.Getenv("PSQL_USER"), "", os.Getenv("PSQL_DATABASE"))


	db, err := sql.Open("postgres", pg_connection_string)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("DB connected")

	fmt.Println("server running...")
	handleRequests()
}