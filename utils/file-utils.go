package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"urbskali/file/models"

	"github.com/mholt/archiver"
)

func CompressDir(input string, output string) error {
	err := archiver.Archive([]string{input}, output)
	if err != nil {
		return err
	}
	return nil
}

func IsDir(path string) bool {
	// Get the FileInfo for the given path
	fileInfo, err := os.Stat(path)
	if err != nil {
		// If an error occurs, it means the path does not exist or is inaccessible
		return false
	}

	// Check if the FileInfo represents a directory
	return fileInfo.IsDir()
}

func GetFiles(path string) ([]models.File, error) {
	// get the list of files in the directory
	// print the url requested
	if path == "./files/" {
		path = "./files"
	}
	files, err := os.ReadDir(fmt.Sprintf("./files/%s", path))
	data := []models.File{}
	if err != nil {
		// Warn and send 404 if the directory is not found
		log.Printf("Directory not found: %s\n", path)
		return data, err
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
		newFile := models.File{
			Name:  file.Name(),
			Size:  fileInfo.Size(),
			Date:  fileInfo.ModTime().String(),
			Type:  fileType,
			Url:   fmt.Sprintf("%s/%s", path, file.Name()),
			IsDir: fileInfo.IsDir(),
		}
		// append the file to the data
		data = append(data, newFile)
	}
	return data, nil
}

func LoadConfig() (models.Config, error) {
	config := models.Config{}
	f, err := os.Open("./config.json")
	if err != nil {
		return config, err
	}
	defer f.Close()
	json.NewDecoder(f).Decode(&config)
	return config, nil
}
