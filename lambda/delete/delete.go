package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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
var db = dynamodb.New(session.Must(session.NewSession()), aws.NewConfig().WithRegion("us-east-1"))



func deleteItem(args Request) Response{

	result, err := db.DeleteItem(&dynamodb.DeleteItemInput{TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Key": {
				S: aws.String(args.Key),
			},

		}})
	fmt.Printf(result.String())
	if err != nil {
		fmt.Println("Got error calling DeleteItem")
		fmt.Println(err.Error())
		return Response{Ack: false}
	}

	fmt.Println("Deleted " + args.Key + "sensor")
	return Response{Ack: true}
}


func Delete(req Request) (Response, error) {
	return deleteItem(req), nil
}

func main() {
	lambda.Start(Delete)
}