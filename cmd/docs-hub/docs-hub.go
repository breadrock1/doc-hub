package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"docs-hub/cmd"
	"docs-hub/internal/cloud/s3minio"
	"docs-hub/internal/server"
	"docs-hub/internal/server/httpserv"
)

func main() {
	servConfig := cmd.Execute()

	cloudService := s3minio.New(&servConfig.Cloud)

	ctx, cancel := context.WithCancel(context.Background())
	go awaitSystemSignals(cancel)

	httpServer := httpserv.Init(&servConfig.Server, cloudService)
	go func() {
		err := httpServer.Server.Start(ctx)
		if err != nil {
			log.Println(err)
			cancel()
		}
	}()

	<-ctx.Done()
	cancel()
	shutdownServices(ctx, httpServer)
}

func awaitSystemSignals(cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	cancel()
}

func shutdownServices(ctx context.Context, httpServ *server.Server) {
	if err := httpServ.Server.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}
