package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func Readiness(w http.ResponseWriter, r *http.Request) {
	fmt.Println("readiness status ok")
	w.WriteHeader(http.StatusOK)
}

func Liveness(w http.ResponseWriter, r *http.Request) {
	fmt.Println("liveness status ok")
	w.WriteHeader(http.StatusOK)
}

func GracefulShutdown(server *http.Server, gracefulStop <-chan os.Signal) {
	<-gracefulStop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
}