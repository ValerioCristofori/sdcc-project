package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func createLambdaServiceClient() *lambda.Lambda {
	// Create Lambda service client
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return lambda.New(sess, &aws.Config{Region: aws.String(configuration.AwsRegion)})
}

func GetLambda(request Args)  Response{

	client := createLambdaServiceClient()

	payload, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshalling Get request")
		return Response{}
	}
	result, err := client.Invoke(&lambda.InvokeInput{FunctionName: aws.String("get"), Payload: payload})
	if err != nil {
		fmt.Println("Error calling Get ", err)
		return Response{}
	}

	var resp Response

	err = json.Unmarshal(result.Payload, &resp)
	if err != nil {
		fmt.Println("Error unmarshalling Get response")
		os.Exit(0)
	}

	return resp
}

func PutLambda(request Args) {

	client := createLambdaServiceClient()

	payload, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshalling Put request")
		return
	}
	result, err := client.Invoke(&lambda.InvokeInput{FunctionName: aws.String("put"), Payload: payload})
	if err != nil {
		fmt.Println("Error calling Put ", err)
		return
	}

	var resp Response

	err = json.Unmarshal(result.Payload, &resp)
	if err != nil {
		fmt.Println("Error unmarshalling Put response")
		return
	}

	if !resp.Ack {
		log.Println("Error: PUT LAMBDA FUNC")
	}
}

func DeleteLambda(request Args) {

	client := createLambdaServiceClient()

	payload, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshalling Delete request")
		return
	}
	result, err := client.Invoke(&lambda.InvokeInput{FunctionName: aws.String("delete"), Payload: payload})
	if err != nil {
		fmt.Println("Error calling Delete ", err)
		return
	}

	var resp Response

	err = json.Unmarshal(result.Payload, &resp)
	if err != nil {
		fmt.Println("Error unmarshalling Delete response")
		return
	}

	if !resp.Ack {
		log.Println("Error: DELETE LAMBDA FUNC")
	}
}

func AppendLambda(request Args) {

	client := createLambdaServiceClient()

	payload, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshalling Append request")
		return
	}
	result, err := client.Invoke(&lambda.InvokeInput{FunctionName: aws.String("append"), Payload: payload})
	if err != nil {
		fmt.Println("Error calling Append ", err)
		return
	}

	var resp Response

	err = json.Unmarshal(result.Payload, &resp)
	if err != nil {
		fmt.Println("Error unmarshalling Append response")
		return
	}

	if !resp.Ack {
		log.Println("Error: APPEND LAMBDA FUNC")
	}
}