package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joecomscience/prom-webhook/pkg/channels"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	const (
		port = ":3000"
	)

	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	server := startServer(port)
	go gracefulShutdown(server, gracefulStop)

	fmt.Printf("server start on port : %s\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Could not listen on %s: %v\n", port, err)
	}
}

func readiness(w http.ResponseWriter, r *http.Request) {
	fmt.Println("readiness status ok")
	w.WriteHeader(http.StatusOK)
}

func liveness(w http.ResponseWriter, r *http.Request) {
	fmt.Println("liveness status ok")
	w.WriteHeader(http.StatusOK)
}

func startServer(port string) *http.Server {
	r := mux.NewRouter()
	r.HandleFunc("/readiness", readiness).Methods("GET")
	r.HandleFunc("/liveness", liveness).Methods("GET")
	r.HandleFunc("/line", channels.LineHandler).Methods("POST")
	r.HandleFunc("/sms", channels.SmsHandler).Methods("POST")

	return &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("%s", port),
	}
}

func gracefulShutdown(server *http.Server, gracefulStop <-chan os.Signal) {
	<-gracefulStop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
}
