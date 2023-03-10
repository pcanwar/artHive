package handlers

import (
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/gorilla/mux"
)

func List(req events.APIGatewayProxyRequest, tableName string, dyna dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	// Create the key for the DynamoDB GetItem operation.
	ownerAddress := mux.Vars(req)["ownerAddress"]
	artworkID := mux.Vars(req)["artworkID"]
	key := map[string]*dynamodb.AttributeValue{
		"address": {
			S: aws.String(ownerAddress),
		},
		"_id": {
			S: aws.String(artworkID),
		},
	}

	input := &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(tableName),
	}

	output, err := dyna.GetItem(input)

	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecords)
	}

	if len(output.Item) == 0 || output.Item == nil {
		return nil, nil
	}

	artwork := new(Artwork)
	err = dynamodbattribute.UnmarshalMap(output.Item, artwork)

	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return artwork, nil
}
