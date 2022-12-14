package todo

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/chandanacharya1/sda-aws-todo/pkg/validators"
)

var (
	ErrorFailedToUnmarshalRecord = "failed to unmarshal record"
	ErrorFailedToFetchRecord     = "There was an error loading the To-Do objects."
	ErrorInvalidTodoData         = "invalid id data"
	ErrorInvalidID               = "invalid ID"
	ErrorCouldNotMarshalItem     = "could not marshal item"
	ErrorCouldNotDeleteItem      = "There has been an error while deleting the To-Do object."
	ErrorCouldNotDynamoPutItem   = "There was an error while creating a new To-Do object."
	ErrorTodoAlreadyExists       = "todo already exists"
	ErrorTodoDoesNotExist        = "todo does not exist"
)

type Todo struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Success struct {
	Message string `json:"message"`
}

func FetchTodo(id, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*Todo, error) {

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.GetItem(input)
	if result == nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	} else if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	} else {
		item := new(Todo)
		err = dynamodbattribute.UnmarshalMap(result.Item, item)
		if err != nil {
			return nil, errors.New(ErrorFailedToUnmarshalRecord)
		}
		if item.Id == "" {
			return nil, errors.New(ErrorFailedToFetchRecord)
		}
		return item, nil
	}
}

func FetchTodos(tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*[]Todo, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}
	item := new([]Todo)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	return item, nil
}

func CreateTodo(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*Success,
	error,
) {
	var t Todo
	message := Success{
		"To-Do object created successfully.",
	}
	if err := json.Unmarshal([]byte(req.Body), &t); err != nil {
		return nil, errors.New(ErrorInvalidTodoData)
	}
	if !validators.IsValidUUID(t.Id) {
		return nil, errors.New(ErrorInvalidID)
	}
	currentTodo, _ := FetchTodo(t.Id, tableName, dynaClient)
	if currentTodo != nil && len(currentTodo.Id) != 0 {
		return nil, errors.New(ErrorTodoAlreadyExists)
	}

	av, err := dynamodbattribute.MarshalMap(t)

	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &message, nil
}

func UpdateTodo(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*Success,
	error,
) {
	var t Todo
	message := Success{
		"To-Do object updated successfully.",
	}
	if err := json.Unmarshal([]byte(req.Body), &t); err != nil {
		return nil, errors.New(ErrorInvalidID)
	}

	currentTodo, _ := FetchTodo(t.Id, tableName, dynaClient)
	if currentTodo != nil && len(currentTodo.Id) == 0 {
		return nil, errors.New(ErrorTodoDoesNotExist)
	}

	av, err := dynamodbattribute.MarshalMap(t)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &message, nil
}

func DeleteTodo(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*Success, error) {
	id := req.QueryStringParameters["id"]
	if !validators.IsValidUUID(id) {
		return nil, errors.New(ErrorCouldNotDeleteItem + ErrorInvalidID)
	}
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(tableName),
	}
	message := Success{
		"To-Do object deleted successfully.",
	}
	currentTodo, _ := FetchTodo(id, tableName, dynaClient)
	if currentTodo != nil && len(currentTodo.Id) == 0 {
		return nil, errors.New(ErrorCouldNotDeleteItem + ErrorTodoDoesNotExist)
	} else {
		_, err := dynaClient.DeleteItem(input)
		if err != nil {
			return nil, errors.New(ErrorCouldNotDeleteItem)
		}

		return &message, nil
	}
}
