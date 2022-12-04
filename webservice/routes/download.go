package routes

import (
	"chunk-uploader/repository/downloaderpb"
	"chunk-uploader/webservice/rpc"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
)

// Download the uploaded image
func (r *Router) Download(c *fiber.Ctx) error {
	r.log.WithFields(logrus.Fields{
		"sha256": c.Params("sha256"),
	}).Info("downloading image")

	var resp response
	conn := rpc.Connection("downloader:50052")
	defer conn.Close()

	req := &downloaderpb.DownloadRequest{
		Sha256: c.Params("sha256"),
	}

	cl := downloaderpb.NewDownloadServiceClient(conn)

	res, err := cl.Download(c.Context(), req)

	if err != nil {
		r.log.Error(err)
		st, ok := status.FromError(err)
		if !ok {
			return resp.badRequest(c)
		}

		return resp.Send(c, st)
	}

	return c.Status(200).SendString(res.GetData())
}
