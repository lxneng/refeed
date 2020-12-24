package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/lxneng/refeed/config"
	"github.com/lxneng/refeed/handlers"
)

func main() {
	config.Init()
	r := mux.NewRouter()
	handler := &handlers.Handler{Router: r}
	r.HandleFunc("/", handler.IndexHandler).Methods("GET")
	r.HandleFunc("/{slug}", handler.FeedHandler).Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         ":5000",
		ReadTimeout:  300 * time.Second,
		WriteTimeout: 300 * time.Second,
		IdleTimeout:  300 * time.Second,
	}

	log.Printf(`Listening on %q`, srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(`Server failed to start: %v`, err)
	}
}
