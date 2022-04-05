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
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"github.com/von-salumbides/go-sls-lambda/pkg/event"
	"github.com/von-salumbides/go-sls-lambda/pkg/models"
	"github.com/von-salumbides/go-sls-lambda/utils/logger"
	"go.uber.org/zap"
)

func Handler(request events.APIGatewayV2HTTPRequest) (*event.Response, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// svc create dynamodb client
	svc := dynamodb.New(sess)

	// itemUuid is a new uuid for item id
	itemUuid := uuid.New().String()
	zap.L().Info("Generate new item uuid", zap.Any("itemId", itemUuid))

	// itemString unmarshal to item to access object properties
	itemString := request.Body
	itemStruct := models.Item{}
	json.Unmarshal([]byte(itemString), &itemStruct)
	if itemStruct.Title == "" {
		return &event.Response{
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	// Create of new item of type item
	item := models.Item{
		Id:      itemUuid,
		Title:   itemStruct.Title,
		Details: itemStruct.Details,
	}

	// Marshal to dynamodb item
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		zap.L().Fatal("Error marshalling item", zap.Any("msg", err.Error()))
		return &event.Response{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	tableName := os.Getenv("DYNAMODB_TABLE")

	// Build put item input
	zap.L().Info("Putting Item:", zap.Any("val", av))
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	// Put item request
	_, err = svc.PutItem(input)
	if err != nil {
		zap.L().Fatal("Got error calling put item", zap.Any("error", err.Error()))
		return &event.Response{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	// marshal item to return
	itemMarshalled, err := json.Marshal(item)
	zap.L().Info("Returning Item", zap.Any("item", itemMarshalled))
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
