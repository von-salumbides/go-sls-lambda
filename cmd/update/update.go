package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/von-salumbides/go-sls-lambda/pkg/event"
	"github.com/von-salumbides/go-sls-lambda/pkg/models"
	"github.com/von-salumbides/go-sls-lambda/utils/logger"
	"go.uber.org/zap"
)

func Handler(request events.APIGatewayV2HTTPRequest) (*event.Response, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// create dynamodb client
	svc := dynamodb.New(sess)

	pathParamId := request.PathParameters["id"]

	itemString := request.Body
	itemStruct := models.Item{}
	json.Unmarshal([]byte(itemString), &itemStruct)

	info := models.Item{
		Title:   itemStruct.Title,
		Details: itemStruct.Details,
	}

	zap.L().Info("Updating data...", zap.Any("title", info.Title), zap.Any("details", info.Details))

	// prepare input for update
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {
				S: aws.String(info.Title),
			},
			":d": {
				S: aws.String(info.Details),
			},
		},
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE")),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(pathParamId),
			},
		},
		ReturnValues:     aws.String("UPDATE_NEW"),
		UpdateExpression: aws.String("set title = :t, details = :d"),
	}

	// update item request
	_, err := svc.UpdateItem(input)
	if err != nil {
		zap.L().Fatal("Failed update", zap.Any("error", err))
		return &event.Response{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	return &event.Response{
		StatusCode: http.StatusOK,
	}, err
}
func Init() {
	logger.InitLogger()
}

func main() {
	lambda.Start(Handler)
}
