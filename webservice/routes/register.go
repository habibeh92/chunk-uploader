package routes

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"os"
)

type register struct {
	Sha256    *string `json:"sha256"`
	Size      float64 `json:"size,omitempty"`
	ChunkSize float64 `json:"chunk_size,omitempty"`
}

// Register register the image sha256
func (r *Router) Register(c *fiber.Ctx) error {
	var resp response
	var re register
	err := c.BodyParser(&re)
	if err != nil {
		r.log.Error(err)
		return resp.badRequest(c)
	}

	r.log.WithFields(logrus.Fields{
		"sha256": *re.Sha256,
	}).Info("registering image")

	if re.Sha256 == nil {
		r.log.Error(errors.New("sha256 field is missing"))
		return resp.invalidArgument(c, "Sha256 field is missing")
	}

	err = re.createImagesDir()
	if err != nil {
		r.log.Error(err)
		return resp.badRequest(c)
	}

	if re.checkImage() {
		r.log.Error(errors.New("image already exists"))
		return resp.conflict(c, "Image already exists")
	}

	// create image directory
	err = os.Mkdir(re.path(), 0755)
	if err != nil {
		r.log.Error(err)
		return resp.unsupportedFormat(c)
	}

	return resp.created(c, "Image successfully registered")
}

// checkImage check if the image is already exists
func (r *register) checkImage() bool {
	if _, err := os.Stat(r.path()); !os.IsNotExist(err) {
		return true
	}

	return false
}

// path get the path of the image
func (r *register) path() string {
	return fmt.Sprintf("./repository/images/%s", *r.Sha256)
}

// createImagesDir create images directory if it does not exist
func (r *register) createImagesDir() error {
	if _, err := os.Stat("./repository/images"); os.IsNotExist(err) {
		return os.Mkdir("./repository/images", 0755)
	}

	return nil
}
