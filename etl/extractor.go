package etl

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/shivasaicharanruthala/dataops-takehome/log"
	"github.com/shivasaicharanruthala/dataops-takehome/model"
	"io"
	"net/http"
)

type extractor struct {
	httpClient    *http.Client
	logger        *log.CustomLogger
	sqsEndpoint   string
	encryptionKey string
}

// NewExtractor creates a new instance of the Extractor and initializes it with the SQS endpoint from environment variables.
func NewExtractor(logger *log.CustomLogger, sqsEndpoint string, encryptionKey string) Extract {
	return &extractor{
		httpClient:    new(http.Client),
		logger:        logger,
		sqsEndpoint:   sqsEndpoint,
		encryptionKey: encryptionKey,
	}
}

// FetchDataFromSQS fetches data from the SQS endpoint, processes the response, and returns a model.Response.
func (ex extractor) FetchDataFromSQS() (*model.Response, error) {
	// Create a new GET request to the SQS endpoint.
	req, err := http.NewRequest("GET", ex.sqsEndpoint, nil)
	if err != nil {
		lm := log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Error creating request to sqs enpoint: %v", err.Error())}
		ex.logger.Log(&lm)

		return nil, err
	}

	// Send the request and receive the response.
	resp, err := ex.httpClient.Do(req)
	if err != nil {
		lm := log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Error sending request to sqs enpoint: %v", err.Error())}
		ex.logger.Log(&lm)

		return nil, err
	}

	defer resp.Body.Close()

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		lm := log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Error reading response from sqs enpoint : %v", err.Error())}
		ex.logger.Log(&lm)

		return nil, err
	}

	// Unmarshal the XML response into the ReceiveMessageResponse struct.
	var sqsMessageResponse model.ReceiveMessageResponse
	err = xml.Unmarshal(body, &sqsMessageResponse)
	if err != nil {
		lm := log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Error unmarshalling XML response from sqs enpoint : %v", err.Error())}
		ex.logger.Log(&lm)

		return nil, err
	}

	// Initialize a new Response struct.
	var res model.Response
	if sqsMessageResponse.ReceiveMessageResult.Message != nil {
		// Unmarshal the JSON body of the SQS message into the Response struct.
		err = json.Unmarshal([]byte(sqsMessageResponse.ReceiveMessageResult.Message.Body), &res)
		if err != nil {
			lm := log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Error unmarshalling JSON body from XML response from sqs enpoint : %v", err.Error())}
			ex.logger.Log(&lm)

			fmt.Println("Error unmarshalling JSON body:", err)
			return nil, err
		}

		// Set additional data from the SQS message response into the Response struct.
		res.SetData(sqsMessageResponse)
	}

	// Mask sensitive data in the Response struct.
	err = res.MaskBody(ex.encryptionKey)
	if err != nil {
		return nil, err
	}

	// Validation of fields
	//_ = res.Validation()

	// Return the populated Response struct.
	return &res, nil
}
