package main

import (
	"fmt"
	"multipart-upload-to-s3-using-presign-url/server/services/aws"
	"multipart-upload-to-s3-using-presign-url/server/services/hello"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

func main() {
	server := rpc.NewServer()
	rtr := mux.NewRouter()
	server.RegisterCodec(json.NewCodec(), "application/json")

	server.RegisterService(new(hello.HelloService), "")
	server.RegisterService(new(aws.AWSService), "")
	rtr.Handle("/api/", server)

	headersOk := handlers.AllowedHeaders([]string{"Content-Type", "X-Requested-With"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	originsOk := handlers.AllowedOrigins([]string{"*"})

	fmt.Println("Started server at :10021")
	err := http.ListenAndServe(":10021", handlers.CORS(originsOk, headersOk, methodsOk)(rtr))
	fmt.Println(err)
}
