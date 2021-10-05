package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"os"
	"strings"
)
var table_name string

func initDynamoDB(tableName string) error{
	// Initialize a session in us-east-1 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	table_name = tableName
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		fmt.Println("Got error calling NewSession:")
		fmt.Println(err.Error())
		return err
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// Create table Movies
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Key"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Key"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(table_name),
	}

	_, err = svc.CreateTable(input)

	if err != nil {

		fmt.Println("Got error calling CreateTable:")
		fmt.Println(err.Error())
		return err
	}
	return nil
}

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
		os.Exit(1)
	}

	// Create item in table
	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(table_name),
	}

	_, err = svc.PutItem(input)

	if strings.Contains(err.Error(), "ResourceInUseException"){
		fmt.Println(err.Error())
	}  else if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		//os.Exit(1)
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

	result, err := svc.DeleteItem(&dynamodb.DeleteItemInput{TableName: aws.String(table_name),
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

func getItem (key string){
	// Initialize a session in us-east-1 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("Movies"),
		Key: map[string]*dynamodb.AttributeValue{
			"Key": {
				S: aws.String(key),
			},
		},
	})

	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}

	item := Args{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)

	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	if item.Value == "" {
		fmt.Println("Could not find ", key, " sensor")
		return
	}

	fmt.Println("Found item:")
	fmt.Println("Key:  ", item.Key)
	fmt.Println("Value: ", item.Value)

}


func callTable() int {
	// Initialize a session in us-east-1 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// create the input configuration instance
	input := &dynamodb.ListTablesInput{}

	fmt.Printf("Tables:\n")

	for {
		// Get the list of tables
		result, err := svc.ListTables(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeInternalServerError:
					fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return 0
		}

		return len(result.TableNames)

	}
}