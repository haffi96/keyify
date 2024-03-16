package utils

import (
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

func HttpErrorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	log.Fatalf("%s - - Failed with http status code: %d. Error: %s", time.Now().Format(time.RFC3339), statusCode, message)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       message,
	}
}
