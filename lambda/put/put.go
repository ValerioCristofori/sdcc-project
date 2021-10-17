package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

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
// Declare a new DynamoDB instance. Note that this is safe for concurrent
// use.
var db = dynamodb.New(session.Must(session.NewSession()), aws.NewConfig().WithRegion("us-east-1"))

func putItem(args Request) Response{

	av, err := dynamodbattribute.MarshalMap(args)

	if err != nil {
		fmt.Println("Got error marshalling map:")
		fmt.Println(err.Error())
		return Response{Ack: false}
	}

	// Create item in table
	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(tableName),
	}

	_, err = db.PutItem(input)

	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		return Response{Ack: false}
	}

	fmt.Println("Successfully added" + args.Key + "sensor")
	return Response{ Ack: true }
}

func Put(req Request) (Response, error) {
	return putItem(req), nil
}

func main() {
	lambda.Start(Put)
}