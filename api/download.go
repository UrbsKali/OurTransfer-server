package api

import (
	"fmt"
	"net/url"
	"strings"
	"urbskali/file/utils"

	"github.com/gofiber/fiber/v2"
)

func Download(c *fiber.Ctx) error {
	// url decode the file path
	file, _ := url.QueryUnescape(c.Params("*"))
	fmt.Println("[DOWNLOAD] " + file)
	// if the file is a directory, compress it and send it
	if utils.IsDir(fmt.Sprintf("./files/%s", file)) {
		// get the directory name
		dirName := strings.Split(file, "/")
		dirName = dirName[:len(dirName)-1]
		dirName = append(dirName, "compressed")
		// compress the directory
		err := utils.CompressDir(fmt.Sprintf("./files/%s", file), fmt.Sprintf("./tmp/%s.zip", strings.Join(dirName, "/")))
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
}
