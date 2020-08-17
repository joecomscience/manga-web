package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joecomscience/prom-webhook/pkg/channels/sms"
	"github.com/joecomscience/prom-webhook/pkg/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	const (
		port = ":3000"
	)

	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	s := startServer(port)
	go server.GracefulShutdown(s, gracefulStop)

	fmt.Printf("server start on port : %s\n", port)
	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("Could not listen on %s: %v\n", port, err)
	}
}

func startServer(port string) *http.Server {
	r := mux.NewRouter()

	r.HandleFunc("/readiness", server.Readiness).Methods("GET")
	r.HandleFunc("/liveness", server.Liveness).Methods("GET")
	r.HandleFunc("/sms", sms.Handler).Methods("POST")

	return &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("%s", port),
	}
}
