package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/neil-berg/chef/db"
	"github.com/neil-berg/chef/handlers"
)

func main() {
	logger := log.New(os.Stdout, "chef-api", log.LstdFlags)

	_, err := db.Connect()
	if err != nil {
		logger.Fatal("Unable to connect to database")
	}
	logger.Println("Connected to DB!")

	router := mux.NewRouter()

	router.HandleFunc("/test", handlers.GetUsers).Methods("GET")

	serverPort := os.Getenv("SERVER_PORT")
	serverAddress := ":" + serverPort

	server := http.Server{
		Addr:         serverAddress,
		Handler:      router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		logger.Println("Server listening on port: ", serverPort)
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	logger.Printf("Terminating server [%s], gracefully shutting down...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
