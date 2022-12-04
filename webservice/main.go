package main

import (
	"chunk-uploader/webservice/routes"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	app := fiber.New()
	routes.CreateRouter(app).Routing()
	log.Fatal(app.Listen(":4444"))
}
