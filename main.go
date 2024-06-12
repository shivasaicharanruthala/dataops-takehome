package main

import (
	"context"
	"fmt"
	"github.com/shivasaicharanruthala/dataops-takehome/database"
	"github.com/shivasaicharanruthala/dataops-takehome/etl"
	"github.com/shivasaicharanruthala/dataops-takehome/log"
	"github.com/shivasaicharanruthala/dataops-takehome/model"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	_ "time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// init function runs before the main function. It loads environment variables from a .env file.
func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func main() {
	// Retrieve the number of workers and batch size from environment variables and convert them to integers.
	noOfWorkers, _ := strconv.Atoi(os.Getenv("NO_OF_WORKERS"))
	batchSize, _ := strconv.Atoi(os.Getenv("BATCH_SIZE"))
	sqsEndpoint := os.Getenv("SQS_ENDPOINT")
	encryptionKey := os.Getenv("ENCRYPTION_SECRET")

	// Initialize Logger
	logger, err := log.NewCustomLogger("logs")
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

	lm = log.Message{Level: "INFO", Msg: fmt.Sprintf("Database initilized sucessfully.")}
	logger.Log(&lm)

	// WaitGroup to synchronize goroutines
	var wg sync.WaitGroup
	// Channel to collect results from workers
	results := make(chan *model.Response, noOfWorkers*batchSize)
	// Context to handle singling to go routines to terminate
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize the ETL components.
	extractor := etl.NewExtractor(logger, sqsEndpoint, encryptionKey)
	loader := etl.NewLoader(logger, dbConn)
	processor := etl.NewProcessor(logger, &wg, extractor, loader)

	lm = log.Message{Level: "INFO", Msg: fmt.Sprintf("Extractor, Loader, Processor initilized sucessfully.")}
	logger.Log(&lm)

	// Start worker goroutines
	for idx := 0; idx < noOfWorkers; idx++ {
		wg.Add(1)

		lm = log.Message{Level: "INFO", Msg: fmt.Sprintf("Worker %v assigned to extract data.", idx)}
		logger.Log(&lm)

		go processor.Worker(ctx, idx, results)
	}

	// Channel to listen for termination signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Goroutine to handle batching and inserting data from workers
	go processor.ProcessDataFromWorker(ctx, results)

	// Block until a signal is received
	sig := <-sigChan
	lm = log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Received signal: %v. Shutting down gracefully...", sig)}
	logger.Log(&lm)
	cancel() // Cancel the context to stop worker goroutines

	// Wait for all workers to finish
	wg.Wait()
	lm = log.Message{Level: "INFO", Msg: fmt.Sprintf("All workers have finished.")}
	logger.Log(&lm)
}
