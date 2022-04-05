package main

import (
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/von-salumbides/go-sls-lambda/pkg/event"
	"github.com/von-salumbides/go-sls-lambda/utils/logger"
	"go.uber.org/zap"
)

func Handler(request events.APIGatewayV2HTTPRequest) (*event.Response, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// create dynamodb service
	svc := dynamodb.New(sess)

	pathParamId := request.PathParameters["id"]

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(pathParamId),
			},
		},
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE")),
	}
	//delete item request
	_, err := svc.DeleteItem(input)
	if err != nil {
		zap.L().Fatal("Got error calling Delete Item", zap.Any("error", err.Error()))
		return &event.Response{
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	zap.L().Info("Successfully deleted item", zap.Any("result", pathParamId))
	return &event.Response{
		StatusCode: http.StatusOK,
	}, nil
}

func init() {
	logger.InitLogger()
}

func main() {
	lambda.Start(Handler)
}
