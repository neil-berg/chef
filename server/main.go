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
	"github.com/rs/cors"
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

	// Unauthenticated routes
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from Chef server!"))
	}).Methods("GET")
	router.HandleFunc("/signup", handler.CreateUser).Methods("POST")
	router.HandleFunc("/signin", handler.SignInUser).Methods("POST")

	// Authenticated GET routes
	authGetRouter := router.Methods("GET").Subrouter()
	authGetRouter.HandleFunc("/recipes/{recipeID}", handler.GetRecipe)
	authGetRouter.HandleFunc("/recipes", handler.GetRecipes)
	authGetRouter.Use(handler.CheckToken)

	// Authenticated POST routes
	authPostRouter := router.Methods("POST").Subrouter()
	authPostRouter.HandleFunc("/auth/me", handler.AuthMe)
	authPostRouter.HandleFunc("/recipes/add", handler.AddRecipe)
	authPostRouter.Use(handler.CheckToken)

	// Authenticated PUT routes
	authPutRouter := router.Methods("PUT").Subrouter()
	authPutRouter.HandleFunc("/recipes/{recipeID}", handler.UpdateRecipe)
	authPutRouter.Use(handler.CheckToken)

	// Authenticated DELETE routes
	authDeleteRouter := router.Methods("DELETE").Subrouter()
	authDeleteRouter.HandleFunc("/recipes/{recipeID}", handler.DeleteRecipe)
	authDeleteRouter.HandleFunc("/me/delete", handler.DeleteMe)
	authDeleteRouter.Use(handler.CheckToken)

	serverAddress := ":" + config.ServerPort

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{config.ClientURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowCredentials: true,
	}).Handler(router)

	server := http.Server{
		Addr:         serverAddress,
		Handler:      corsHandler,
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
