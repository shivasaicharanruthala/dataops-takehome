package etl

import (
	"context"
	"fmt"
	"github.com/shivasaicharanruthala/dataops-takehome/log"
	"github.com/shivasaicharanruthala/dataops-takehome/model"

	"os"
	"strconv"
	"sync"
	"time"
)

type transformer struct {
	logger    *log.CustomLogger
	extractor Extract
	loader    Loader
	wg        *sync.WaitGroup
}

// NewProcessor creates a new instance of the Processor with the given extractor and loader.
func NewProcessor(logger *log.CustomLogger, wg *sync.WaitGroup, extractor Extract, loader Loader) Processor {
	return &transformer{
		logger:    logger,
		extractor: extractor,
		loader:    loader,
		wg:        wg,
	}
}

// Worker is a function that continuously calls the API to fetch data and sends the result to a channel.
func (p *transformer) Worker(ctx context.Context, id int, results chan<- *model.Response) {
	defer p.wg.Done() // Ensure the WaitGroup counter is decremented when the function returns

	// Retrieve maximum allowed empty responses and maximum consecutive empty responses from environment variables
	maxEmptyResponses, _ := strconv.Atoi(os.Getenv("MAX_NO_RESPONSES"))
	maxConsecutiveEmptyResponses, _ := strconv.Atoi(os.Getenv("MAX_CONSECUTIVE_NO_RESPONSES"))
	initialWaitTime := time.Second
	emptyResponseCount := 0
	waitTime := initialWaitTime

	for {
		select {
		case <-ctx.Done():
			lm := log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Worker %d: Stopping.", id)}
			p.logger.Log(&lm)

			return
		default:
			response, err := p.extractor.FetchDataFromSQS()
			if err != nil {
				lm := log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Worker %d: Error fetching data: %v", id, err.Error())}
				p.logger.Log(&lm)

				continue
			}

			// Send the valid response to the results channel if message exists.
			if response != nil && response.MessageId != nil && response.UserID != nil {
				results <- response
				emptyResponseCount = 0     // Reset the empty response counter
				waitTime = initialWaitTime // Reset the wait time
			} else {
				lm := log.Message{Level: "INFO", Msg: fmt.Sprintf("Worker %d: Received empty response", id)}
				p.logger.Log(&lm)

				emptyResponseCount++
				if emptyResponseCount >= maxEmptyResponses {
					lm = log.Message{Level: "INFO", Msg: fmt.Sprintf("Worker %d: Waiting for %v due to consecutive empty responses", id, waitTime)}
					p.logger.Log(&lm)

					time.Sleep(waitTime)        // Wait for the specified time before retrying
					waitTime += initialWaitTime // Increase the wait time linearly
				}

				if emptyResponseCount >= maxConsecutiveEmptyResponses {
					lm = log.Message{Level: "INFO", Msg: fmt.Sprintf("Worker %d: Reached max consecutive empty responses, canceling context", id)}
					p.logger.Log(&lm)

					ctx.Done() // Cancel the context if the maximum consecutive empty responses are reached
					return
				}
			}
		}
	}
}

func (p *transformer) ProcessDataFromWorker(ctx context.Context, results chan *model.Response) {
	batchSize, _ := strconv.Atoi(os.Getenv("BATCH_SIZE"))

	//TODO: not required pointer to model.Response
	var batch []*model.Response
	for {
		select {
		case response := <-results:
			batch = append(batch, response)
			if len(batch) >= batchSize {
				err := p.loader.BatchInsert(ctx, batch)
				if err != nil {
					lm := log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Error inserting batch: %v", err.Error())}
					p.logger.Log(&lm)

				}
				batch = batch[:0] // Reset batch
			}
		case <-ctx.Done():
			// Insert any remaining items before shutting down
			if len(batch) > 0 {
				err := p.loader.BatchInsert(ctx, batch)
				if err != nil {
					lm := log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Error inserting final batch: %v", err.Error())}
					p.logger.Log(&lm)
				}
			}
			return
		}
	}
}
