// Make a webserver that can serve file and handle requests

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type File struct {
	Name  string
	Size  int64
	Date  string
	Type  string
	Url   string
	IsDir bool
}

func GetFiles(w http.ResponseWriter, r *http.Request) {
	// get the list of files in the directory
	// print the url requested
	newPath := fmt.Sprintf("./files/%s", strings.Join(strings.Split(r.URL.Path, "/")[3:], "/"))
	if newPath == "./files/" {
		newPath = "./files"
	}
	files, err := os.ReadDir(newPath)
	data := []File{}
	if err != nil {
		// Warn and send 404 if the directory is not found
		log.Printf("Directory not found: %s\n", newPath)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	for _, file := range files {
		// get the file info
		fileInfo, err := file.Info()
		if err != nil {
			log.Fatal(err)
		}
		// create a new file object
		fileTypeArr := strings.Split(file.Name(), ".")
		fileType := fileTypeArr[len(fileTypeArr)-1]
		newFile := File{
			Name:  file.Name(),
			Size:  fileInfo.Size(),
			Date:  fileInfo.ModTime().String(),
			Type:  fileType,
			Url:   fmt.Sprintf("%s/%s", newPath, file.Name()),
			IsDir: fileInfo.IsDir(),
		}
		// append the file to the data
		data = append(data, newFile)
	}

	// format in json and send
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(data)
}

func main() {
	fileServer := http.FileServer(http.Dir("../file/dist"))
	http.Handle("/", fileServer)

	// handle files download requests from the client
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("./files"))))

	http.HandleFunc("/api/get_files/", GetFiles)

	fmt.Printf("Starting server at port 8080\nhttp://localhost:8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
