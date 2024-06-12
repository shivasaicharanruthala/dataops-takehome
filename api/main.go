package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/shivasaicharanruthala/dataops-takehome/api/handler"
	"github.com/shivasaicharanruthala/dataops-takehome/api/store"
	"github.com/shivasaicharanruthala/dataops-takehome/database"
	"github.com/shivasaicharanruthala/dataops-takehome/log"
	"net/http"
	"os"
)

// init function runs before the main function. It loads environment variables from a .env file.
func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}
func main() {
	encryptionKey := os.Getenv("ENCRYPTION_SECRET")

	// Initialize Logger
	logger, err := log.NewCustomLogger("../../app_logs")
	if err != nil {
		lm := log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Initiating logger with error %v", err.Error())}
		logger.Log(&lm)
	}

	lm := log.Message{Level: "INFO", Msg: "Logger initialized successfully"}
	logger.Log(&lm)

	// Initialize a new database connection.
	db := database.New(logger)
	dbConn, err := db.Open()
	defer dbConn.Close()
	if err != nil {
		lm = log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Initiating database failed with error %v", err.Error())}
		logger.Log(&lm)

		return
	}

	loginStore := store.New(dbConn, encryptionKey)
	loginHandler := handler.New(loginStore)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/login-data", loginHandler.Get).Methods("GET")

	// Start the server
	port := os.Getenv("PORT")
	server := fmt.Sprintf(":%s", port)

	lm = log.Message{Level: "INFO", Msg: fmt.Sprintf("Server starting to listen on port %v", port)}
	logger.Log(&lm)

	err = http.ListenAndServe(server, router)
	if err != nil {
		lm = log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Initializing weapp server to listen on port %v with error %v", port, err.Error())}
		logger.Log(&lm)
	}
}
