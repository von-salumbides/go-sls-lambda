package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
type Response events.APIGatewayProxyResponse
type DefaultResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (Response, error) {
	var resp Response
	l, _ := zap.NewProduction()
	defer l.Sync()

	l.Info("event received", zap.Any("method", event.HTTPMethod), zap.Any("path", event.Path), zap.Any("body", event.Body))
	if event.Path == "/hello" {
		body, _ := json.Marshal(&DefaultResponse{
			Status:  http.StatusOK,
			Message: "hello lambda",
		})
		resp = Response{
			StatusCode: http.StatusOK,
			Body:       string(body),
		}
	} else {
		body, _ := json.Marshal(&DefaultResponse{
			Status:  http.StatusOK,
			Message: "Default path",
		})
		resp = Response{
			StatusCode: http.StatusOK,
			Body:       string(body),
		}
	}
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
