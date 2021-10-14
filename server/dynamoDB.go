package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

var tableName = "Sensors"


func putItem(args Args){
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	av, err := dynamodbattribute.MarshalMap(args)

	if err != nil {
		fmt.Println("Got error marshalling map:")
		fmt.Println(err.Error())
		return
	}

	// Create item in table
	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)

	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Successfully added" + args.Key + "sensor")
}

func deleteItem(args Args){
	// Initialize a session in us-east-1 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	result, err := svc.DeleteItem(&dynamodb.DeleteItemInput{TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Key": {
				S: aws.String(args.Key),
			},

		}})
	fmt.Printf(result.String())
	if err != nil {
		fmt.Println("Got error calling DeleteItem")
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Deleted " + args.Key + "sensor")

}

func getItem (key string) Args{
	// Initialize a session in us-east-1 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Key": {
				S: aws.String(key),
			},
		},
	})

	if err != nil{
		log.Printf("Got error calling GetItem: %s \n", err)
	}

	item := Args{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)

	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	if item.Value == "" {
		fmt.Println("Could not find ", key, " sensor")
	} else {
		item.Counter++
		fmt.Println("Found item:")
		fmt.Println("Key:  ", item.Key)
		fmt.Println("Value: ", item.Value)
	}


	return item

}


func appendItem(args Args) {
	// Initialize a session in us-east-1 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.

	item := getItem(args.Key)

	if item.Value == "" {
		fmt.Println(args.Key + " sensor isn't in dynamoDB table")
		putItem(args)
		return
	}
	//append value to existing item
	item.Value = item.Value + "\n" + args.Value
	putItem(item)

	fmt.Println("Correctly updated " + args.Key + " sensor value")
}


