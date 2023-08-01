package api

import (
	"fmt"
	"os"
	"strings"
	"urbskali/file/models"

	"github.com/gofiber/fiber/v2"
)

func FileInfo(c *fiber.Ctx) error {
	// URL decode the path
	file := c.FormValue("path")
	// get the file info
	fileInfo, err := os.Stat(fmt.Sprintf("./files/%s", file))
	if err != nil {
		return c.Status(404).SendString("File not found")
	}
	// create a new file object
	fileTypeArr := strings.Split(fileInfo.Name(), ".")
	fileType := fileTypeArr[len(fileTypeArr)-1]
	newFile := models.File{
		Name:  fileInfo.Name(),
		Size:  fileInfo.Size(),
		Date:  fileInfo.ModTime().String(),
		Type:  fileType,
		Url:   file,
		IsDir: fileInfo.IsDir(),
	}
	return c.JSON(newFile)
}
