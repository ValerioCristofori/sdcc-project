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

	return Response{Ack: true, Key: item.Key, Value: item.Value}

}

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

func appendItem(args Request) Response{
	// Initialize a session in us-east-1 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.

	item := getItem(args.Key)

	if item.Value == "" {
		fmt.Println(args.Key + " sensor isn't in dynamoDB table")
		return putItem(args)
	}
	//append value to existing item
	return putItem(Request{Key: args.Key,Value: item.Value + "\n" + args.Value})
}



func Append(req Request) (Response, error) {
	return appendItem(req), nil
}

func main() {
	lambda.Start(Append)
}