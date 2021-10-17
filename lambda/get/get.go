package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

type Args struct {
	Key 	string
	Value 	string
}

type Request struct {
	Key 	string
	Value 	string
}

type Response struct {
	Ack 	bool
	Key 	string
	Value 	string
}

var tableName = "Sensors"
var db = dynamodb.New(session.Must(session.NewSession()), aws.NewConfig().WithRegion("us-east-1"))

func getItem (key string) Response{

	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Key": {
				S: aws.String(key),
			},
		},
	})

	if err != nil{
		log.Printf("Got error calling GetItem: %s \n", err)
		return Response{Ack: false}
	}

	item := Args{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)

	if err != nil {
		log.Printf(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		return Response{Ack: false}
	}

	if item.Value == "" {
		// item not found
		return Response{Ack: false, Key: item.Key}
	} else {
		//found item
		return Response{Ack: true, Key: item.Key, Value: item.Value}
	}




}


func Get(req Request) (Response, error) {
	resp := getItem(req.Key)
	return resp, nil
}

func main() {
	lambda.Start(Get)
}