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
	//create dynamodb client
	svc := dynamodb.New(sess)

	// build the query input parameter
	params := &dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE")),
	}

	//scan table
	result, err := svc.Scan(params)
	if err != nil {
		zap.L().Fatal("Query API called failed", zap.Any("error", err.Error()))
		return &event.Response{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	itemArray := []models.Item{}

	for _, i := range result.Items {
		item := models.Item{}
		//unmarshalmap result.item to item
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			zap.L().Fatal("Got error unmarshalling", zap.Any("error", err.Error()))
			return &event.Response{
				StatusCode: http.StatusInternalServerError,
			}, err
		}
		itemArray = append(itemArray, item)
	}
	zap.L().Info("Succesfully unmarshal item", zap.Any("itemArray", itemArray))
	itemArrayString, err := json.Marshal(itemArray)
	if err != nil {
		zap.L().Fatal("Got error marshalling result", zap.Any("error", err.Error()))
		return &event.Response{
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	return &event.Response{
		StatusCode: http.StatusOK,
		Body:       string(itemArrayString),
	}, nil

}

func init() {
	logger.InitLogger()
}

func main() {
	lambda.Start(Handler)
}
