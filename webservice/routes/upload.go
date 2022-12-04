package routes

import (
	"chunk-uploader/repository/uploaderpb"
	"chunk-uploader/webservice/rpc"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type chunk struct {
	ID   *int64  `json:"id,omitempty"`
	Data *string `json:"data,omitempty"`
	Size float32 `json:"size,omitempty"`
}

// Upload the chunked image
func (r *Router) Upload(c *fiber.Ctx) error {
	var resp response
	var ch chunk

	r.log.WithFields(logrus.Fields{
		"sha256": c.Params("sha256"),
	}).Info("uploading image chunk")

	err := c.BodyParser(&ch)
	if err != nil {
		r.log.Error(err)
		return resp.badRequest(c)
	}

	if ch.ID == nil {
		r.log.Error(errors.New("chunk ID field is missing"))
		return resp.invalidArgument(c, "Chunk ID field is missing")
	}

	if ch.Data == nil {
		r.log.Error(errors.New("data field is missing"))
		return resp.invalidArgument(c, "Data field is missing")
	}

	conn := rpc.Connection("uploader_service:50051")
	defer conn.Close()

	_, err = ch.chunkUpload(conn, c)

	if err != nil {
		r.log.Error(err)
		st, ok := status.FromError(err)
		if !ok {
			return resp.badRequest(c)
		}

		return resp.Send(c, st)
	}

	return resp.created(c, "Chunk successfully uploaded")
}


// chunkUpload upload the chunked file to the uploader service via grpc
func (ch chunk) chunkUpload(conn *grpc.ClientConn, c *fiber.Ctx) (*uploaderpb.ChunkResponse, error) {
	cl := uploaderpb.NewChunkServiceClient(conn)

	req := &uploaderpb.ChunkRequest{
		Id:     *ch.ID,
		Data:   *ch.Data,
		Size:   ch.Size,
		Sha256: c.Params("sha256"),
	}

	return cl.ChunkUpload(c.Context(), req)
}
