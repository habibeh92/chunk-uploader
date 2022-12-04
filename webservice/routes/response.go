package routes

import (
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type response struct {
	Code    string
	Message string
}

// badRequest send http response with status 400
func (r response) badRequest(c *fiber.Ctx) error {
	r.Code = http.StatusText(http.StatusBadRequest)
	r.Message = "Malformed request"
	return c.Status(http.StatusBadRequest).JSON(r)
}

// conflict send http response with status 409
func (r response) conflict(c *fiber.Ctx, msg string) error {
	r.Code = http.StatusText(http.StatusConflict)
	r.Message = msg
	return c.Status(http.StatusConflict).JSON(r)
}

// notFound send http response with status 404
func (r response) notFound(c *fiber.Ctx, msg string) error {
	r.Code = http.StatusText(http.StatusNotFound)
	r.Message = msg
	return c.Status(http.StatusNotFound).JSON(r)
}

// unsupportedFormat send http response with status 415
func (r response) unsupportedFormat(c *fiber.Ctx) error {
	r.Code = http.StatusText(http.StatusUnsupportedMediaType)
	r.Message = "Unsupported payload format"
	return c.Status(http.StatusUnsupportedMediaType).JSON(r)
}

// created send http response with status 201
func (r response) created(c *fiber.Ctx, msg string) error {
	r.Code = http.StatusText(http.StatusCreated)
	r.Message = msg
	return c.Status(http.StatusCreated).JSON(r)
}

// internal send http response with status 500
func (r response) internal(c *fiber.Ctx, msg string) error {
	r.Code = http.StatusText(http.StatusInternalServerError)
	r.Message = msg
	return c.Status(http.StatusInternalServerError).JSON(r)
}

// invalidArgument send http response with status 400
func (r response) invalidArgument(c *fiber.Ctx, msg string) error {
	r.Code = http.StatusText(http.StatusBadRequest)
	r.Message = msg
	return c.Status(http.StatusBadRequest).JSON(r)
}

// Send http response based on status.Status
func (r response) Send(c *fiber.Ctx, st *status.Status) error {
	msg := st.Message()
	switch st.Code() {
	case codes.NotFound:
		return r.notFound(c, msg)
	case codes.AlreadyExists:
		return r.conflict(c, msg)
	case codes.InvalidArgument:
		return r.invalidArgument(c, msg)
	case codes.Internal:
		return r.internal(c, msg)
	default:
		return r.badRequest(c)
	}
}
