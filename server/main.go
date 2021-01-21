package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	cfg "github.com/neil-berg/chef/config"
	"github.com/neil-berg/chef/database"
	"github.com/neil-berg/chef/handlers"
)

func main() {
	logger := log.New(os.Stdout, "chef-api", log.LstdFlags)

	config := cfg.Get()
	logger.Println("Retrieved configuration")

	db, err := database.Connect(config.DBHost, config.DBPort, config.DBName, config.DBUser, config.DBPassword)
	if err != nil {
		logger.Fatal("Unable to connect to database")
	}
	logger.Println("Connected to DB!")

	err = database.Migrate(db)
	if err != nil {
		logger.Fatal("Unable to apply DB migrations")
	}
	logger.Println("Successfully applied DB migrations")

	handler := handlers.CreateHandler(logger, db, config)
	router := mux.NewRouter()

	router.HandleFunc("/signup", handler.CreateUser).Methods("POST")

	serverAddress := ":" + config.ServerPort

	server := http.Server{
		Addr:         serverAddress,
		Handler:      router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		logger.Println("Server listening on port: ", config.ServerPort)
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
