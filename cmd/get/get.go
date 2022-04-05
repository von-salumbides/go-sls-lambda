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

	// Getting id from path parameter
	pathParamId := request.PathParameters["id"]
	zap.L().Info("Derived pathParamId from path params: ", zap.Any("id", pathParamId))

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE")),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(pathParamId),
			},
		},
	})
	if err != nil {
		zap.L().Fatal("Internal server error", zap.Any("msg", err.Error()))
		return &event.Response{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	// checking type
	if len(result.Item) == 0 {
		zap.L().Fatal("0 Item")
		return &event.Response{
			StatusCode: http.StatusNoContent,
		}, err
	}

	// created of item of type Item
	item := models.Item{}
	// UnmarshallMap result.item into item
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		zap.L().Panic("Failed to UnmarshalMap result.Item", zap.Any("err", err.Error()))
	}
	// marshal to type bytes
	marshalledItem, err := json.Marshal(item)
	return &event.Response{
		StatusCode: http.StatusOK,
		Body:       string(marshalledItem),
	}, nil
}

func init() {
	logger.InitLogger()
}

func main() {
	lambda.Start(Handler)
}
