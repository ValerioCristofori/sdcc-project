package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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

	exp := 1
	client := createLambdaServiceClient()

	payload, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshalling Get request")
		return Response{}
	}
exponentialBackOffLabelGet:
	result, err := client.Invoke(&lambda.InvokeInput{FunctionName: aws.String("get"), Payload: payload})
	if err != nil {
		if strings.Contains(err.Error(), "TooManyRequestsException"){
			fmt.Println("Exponential back-off")
			time.Sleep( time.Duration(exp)*time.Millisecond)
			exp = exp*2
			goto exponentialBackOffLabelGet
		}
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

	exp := 1
	client := createLambdaServiceClient()

	payload, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshalling Put request")
		return
	}

exponentialBackOffLabelPut:
	result, err := client.Invoke(&lambda.InvokeInput{FunctionName: aws.String("put"), Payload: payload})
	if err != nil {
		if strings.Contains(err.Error(), "TooManyRequestsException"){
			fmt.Println("Exponential back-off")
			time.Sleep( time.Duration(exp)*time.Millisecond)
			exp = exp*2
			goto exponentialBackOffLabelPut
		}
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

	exp := 1
	client := createLambdaServiceClient()

	payload, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshalling Delete request")
		return
	}

exponentialBackOffLabelDel:
	result, err := client.Invoke(&lambda.InvokeInput{FunctionName: aws.String("delete"), Payload: payload})
	if err != nil {
		if strings.Contains(err.Error(), "TooManyRequestsException"){
			fmt.Println("Exponential back-off")
			time.Sleep( time.Duration(exp)*time.Millisecond)
			exp = exp*2
			goto exponentialBackOffLabelDel
		}
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

	exp := 1
	client := createLambdaServiceClient()

	payload, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshalling Append request")
		return
	}

exponentialBackOffLabelApp:
	result, err := client.Invoke(&lambda.InvokeInput{FunctionName: aws.String("append"), Payload: payload})
	if err != nil {
		if strings.Contains(err.Error(), "TooManyRequestsException"){
			fmt.Println("Exponential back-off")
			time.Sleep( time.Duration(exp)*time.Millisecond)
			exp = exp*2
			goto exponentialBackOffLabelApp
		}
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