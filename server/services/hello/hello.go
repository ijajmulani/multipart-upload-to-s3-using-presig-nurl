package hello

import (
	helloctx "multipart-upload-to-s3-using-presign-url/server/context/hello"
	"net/http"
)

type HelloService struct{}

func (h *HelloService) Greet(r *http.Request, args *helloctx.HelloArgs, reply *helloctx.HelloResponse) error {
	reply.Message = "Hi " + args.Name
	return nil
}
