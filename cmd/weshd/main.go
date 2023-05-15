package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"berty.tech/weshnet"
	"berty.tech/weshnet/pkg/protocoltypes"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/pkg/errors"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(errors.Wrap(err, "failed to create logger"))
	}

	logger.Info("weshd")

	port := 4242

	grpcServer := grpc.NewServer()

	svc, err := weshnet.NewService(weshnet.Opts{Logger: logger})
	if err != nil {
		panic(errors.Wrap(err, "failed to create weshnet server"))
	}
	// FIXME: svc.Close?

	protocoltypes.RegisterProtocolServiceServer(grpcServer, svc)

	grpclog.SetLogger(log.New(os.Stdout, "exampleserver: ", log.LstdFlags)) // FIXME: adapt

	wrappedServer := grpcweb.WrapServer(grpcServer)
	handler := func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Access-Control-Allow-Origin", "*")
		resp.Header().Set("Access-Control-Allow-Headers", "*")
		logger.Debug(fmt.Sprintf("Request: %v", req))
		wrappedServer.ServeHTTP(resp, req)
	}

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(handler),
	}

	if err := httpServer.ListenAndServe(); err != nil {
		panic(errors.Wrap(err, "failed to start http server"))
	}
}
