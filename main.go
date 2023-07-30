// Make a webserver that can serve file and handle requests

package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mholt/archiver/v3"
)

type File struct {
	Name  string
	Size  int64
	Date  string
	Type  string
	Url   string
	IsDir bool
}

func setup() {
	fmt.Println("Initializing...")
	// create ./files directory
	if _, err := os.Stat("./files"); os.IsNotExist(err) {
		err := os.Mkdir("./files", 0755)
		if err != nil {
			fmt.Println("Failed to create files directory")
		}
		fmt.Println("Files directory created successfully")
	} else {
		fmt.Println("Files directory already exists")
	}
	// create a PASSWORD file and ask for your secret
	if _, err := os.Stat("./PASSWORD"); os.IsNotExist(err) {
		file, err := os.Create("./PASSWORD")
		if err != nil {
			fmt.Println("Failed to create PASSWORD file")
		}
		// ask for your secret
		fmt.Print("Enter your secret: ")
		var secret string
		fmt.Scanln(&secret)
		// write the secret to the file
		file.WriteString(secret)
		defer file.Close()
		fmt.Println("PASSWORD created successfully")
	} else {
		fmt.Println("PASSWORD file already exists")
	}
	fmt.Println("Initialization complete")
}

func server() {
	//load password from PASSWORD file
	password, err := os.ReadFile("PASSWORD")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Password:", string(password))
	h := sha256.New()
	h.Write([]byte(string(password)))
	// Convert the sha256 hash to a string
	secret := fmt.Sprintf("%x", h.Sum(nil))
	fmt.Println("Secret:", secret)

	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024 * 1024, // 100 GB
	})

	// serve the ui
	app.Static("/", "./ui")

	// Start the routing
	app.Get("/api/get_files/*", func(c *fiber.Ctx) error {
		files, err := GetFiles(c.Params("*"))
		if err != nil {
			return c.Status(404).SendString("Directory not found")
		}
		return c.JSON(files)
	})

	app.Get("/api/file_info/*", func(c *fiber.Ctx) error {
		// URL decode the path
		file, _ := url.QueryUnescape(c.Params("*"))
		// get the file info
		fileInfo, err := os.Stat(fmt.Sprintf("./files/%s", file))
		if err != nil {
			return c.Status(404).SendString("File not found")
		}
		// create a new file object
		fileTypeArr := strings.Split(fileInfo.Name(), ".")
		fileType := fileTypeArr[len(fileTypeArr)-1]
		newFile := File{
			Name:  fileInfo.Name(),
			Size:  fileInfo.Size(),
			Date:  fileInfo.ModTime().String(),
			Type:  fileType,
			Url:   file,
			IsDir: fileInfo.IsDir(),
		}
		return c.JSON(newFile)
	})

	app.Get("/download/*", func(c *fiber.Ctx) error {
		return c.SendFile("./ui/index.html")
	})

	app.Get("/api/download/*", func(c *fiber.Ctx) error {
		// URL decode the path
		file, _ := url.QueryUnescape(c.Params("*"))
		// if the file is a directory, compress it and send it
		if isDir(fmt.Sprintf("./files/%s", file)) {
			// get the directory name
			dirName := strings.Split(file, "/")
			dirName = dirName[:len(dirName)-1]
			dirName = append(dirName, "compressed")
			// compress the directory
			err := CompressDir(fmt.Sprintf("./files/%s", file), fmt.Sprintf("./tmp/%s.zip", strings.Join(dirName, "/")))
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to compress the folder",
				})
			}
			// send the file
			return c.SendFile(fmt.Sprintf("./tmp/%s.zip", strings.Join(dirName, "/")))
		}
		path := fmt.Sprintf("./files/%s", file)
		return c.Download(path)
	})

	app.Post("/api/delete/*", func(c *fiber.Ctx) error {
		if c.FormValue("secret") != secret {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}
		fmt.Println("[DELETE] " + c.Params("*"))
		// URL decode the path
		file, _ := url.QueryUnescape(c.Params("*"))
		// delete the file
		err := os.RemoveAll(fmt.Sprintf("./files/%s", file))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete the file",
			})
		}
		return c.JSON(fiber.Map{
			"message": "File deleted successfully",
		})
	})

	app.Post("/api/create_dir/*", func(c *fiber.Ctx) error {
		if c.FormValue("secret") != secret {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}
		// URL decode the path
		file, _ := url.QueryUnescape(c.Params("*"))
		// create the directory
		fmt.Println("[CREATE] " + file + c.FormValue("name"))
		err := os.MkdirAll(fmt.Sprintf("./files/%s/%s", file, c.FormValue("name")), 0755)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create the directory",
			})
		}
		return c.JSON(fiber.Map{
			"message": "Directory created successfully",
		})
	})

	app.Post("/api/upload/*", func(c *fiber.Ctx) error {
		fmt.Println("[UPLOAD] " + c.Params("*"))
		// Check is the secret is correct
		if c.FormValue("secret") != secret {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}
		// URL decode the path
		file_path, _ := url.QueryUnescape(c.Params("*"))
		// get the files from the request
		form, err := c.MultipartForm()
		if err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get the files",
			})
		}
		// get the files
		files := form.File["files"]
		// loop through the files
		for _, file := range files {
			// save the file
			fmt.Println("[Upload] Received: ", file.Filename)
			err := c.SaveFile(file, fmt.Sprintf("./files/%s/%s", file_path, file.Filename))
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to save the file",
				})
			}
		}
		return c.JSON(fiber.Map{
			"message": "Files uploaded successfully",
		})
	})

	app.Get("/api/check_secret/*", func(c *fiber.Ctx) error {
		fmt.Println("[Check Secret] Checking secret from IP:", c.IP())
		// URL decode the path
		input_secret, _ := url.QueryUnescape(c.Params("*"))
		if input_secret == secret {
			return c.JSON(fiber.Map{
				"message": true,
			})
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Incorrect secret",
		})
	})
	app.Listen(":80")
	//app.ListenTLS(":443", "./cert/cert.crt", "./cert/cert.key")
}

func CompressDir(input string, output string) error {
	err := archiver.Archive([]string{input}, output)
	if err != nil {
		return err
	}
	return nil
}

func isDir(path string) bool {
	// Get the FileInfo for the given path
	fileInfo, err := os.Stat(path)
	if err != nil {
		// If an error occurs, it means the path does not exist or is inaccessible
		return false
	}

	// Check if the FileInfo represents a directory
	return fileInfo.IsDir()
}

func GetFiles(path string) ([]File, error) {
	// get the list of files in the directory
	// print the url requested
	if path == "./files/" {
		path = "./files"
	}
	files, err := os.ReadDir(fmt.Sprintf("./files/%s", path))
	data := []File{}
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
		newFile := File{
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

func main() {
	if len(os.Args) > 1 && os.Args[1] == "setup" {
		setup()
	} else {
		server()
	}
}
