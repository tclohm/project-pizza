package main

import (
	//"github.com/tclohm/project-pizza/models"
	"net/http"
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
)

func postFile(filename string, targetUrl string) error {
	bodyBuffer := &bytes.Buffer{}
	bodyWriter = multipart.NewWriter(bodyBuffer)

	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	// MARK: -- open file handler
	fileHandler, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}

	defer fileHandler.Close()

	// MARK: -- iocopy
	_, err := io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuffer)
	if err != nil {
		return err
	}

	defer err != nil {
		return err
	}

	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))
	return nil
}

func main() {
	targer_url := "http://localhost:9090/upload"
	filename := "/example.pdg"
	postFile(filename, target_url)
}