package main

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"sdcc-project/rpc-logic/dataformat"
)

// aws-dynamoDB conf
var svc *dynamodb.DynamoDB
var tableName 	= "SensorData"


func GetData() (*dataformat.Data, error) {

	// Call GetItem to get the item from the table.
	// If we encounter an error, print the error message.
	// Otherwise, display information about the item.
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Year": {
				N: aws.String(movieYear),
			},
			"Title": {
				S: aws.String(movieName),
			},
		},
	})
	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}

	// If an item was returned, unmarshall the return value and display values.

	if result.Item == nil {
		msg := "Could not find '" + *title + "'"
		return nil, errors.New(msg)
	}

	data := dataformat.Data{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &data)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return &data, nil
	//fmt.Println("Found item:")
	//fmt.Println("Year:  ", item.Year)
	//fmt.Println("Title: ", item.Title)
	//fmt.Println("Plot:  ", item.Plot)
	//fmt.Println("Rating:", item.Rating)
}

func PutData(data *dataformat.Data) {

	// Marshall that data into a map of AttributeValue objects.
	av, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		log.Fatalf("Got error marshalling new movie item: %s", err)
	}

	// Create the input for PutItem and call it.
	// If an error occurs, print the error and exit.
	// If no error occurs, print an message that the data was added to the table.
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	//year := strconv.Itoa(item.Year)

	//fmt.Println("Successfully added '" + item.Title + "' (" + year + ") to table " + tableName)

}


func createTable() {
	// Call CreateTable.
	// If an error occurs, print the error and exit.
	// If no error occurs, print an message that the table was created.
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Year"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("Title"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Key"),
				KeyType:       aws.String("HASH"), // partition key ( better sort key by timestamp ???)
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err := svc.CreateTable(input)
	if err != nil {
		log.Fatalf("Got error calling CreateTable: %s", err)
	}

	fmt.Println("Created the table", tableName)
}


func InitDynamo()  {

	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc = dynamodb.New(sess)

	// Create Sensor Data table
	createTable()

}
