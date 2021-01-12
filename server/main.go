package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/neil-berg/chef/handlers"
)

func main() {
	logger := log.New(os.Stdout, "chef-api", log.LstdFlags)

	// err := godotenv.Load("../.env")
	// if err != nil {
	// 	logger.Fatal(err)
	// }

	router := mux.NewRouter()

	router.HandleFunc("/test", handlers.GetUsers).Methods("GET")

	// port := os.Getenv("SERVER_PORT")
	// address := ":" + port
	address := ":8080"

	server := http.Server{
		Addr:         address,
		Handler:      router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		logger.Println("Server listening on port 8080")
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
