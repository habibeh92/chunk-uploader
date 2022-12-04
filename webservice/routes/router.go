package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Router struct {
	app *fiber.App
	log *logrus.Logger
}

// CreateRouter Create the router instance with app and logger
func CreateRouter(app *fiber.App) *Router {
	return &Router{app: app, log: logrus.New()}
}

// Routing set up the available routes
func (r *Router) Routing() {
	r.app.Post("image", r.Register)
	r.app.Post("image/:sha256/chunks", r.Upload)
	r.app.Get("image/:sha256", r.Download)
}
