package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"log"
)

var (
		handler = "main"
		runtime = "go1.x"
)

func createFunc(svc *lambda.Lambda, zipFile []byte, funcName string) error{
	createCode := &lambda.FunctionCode{
		ZipFile: zipFile,
	}

	createArgs := &lambda.CreateFunctionInput{
		Code:         createCode,
		FunctionName: &funcName,
		Handler:      &handler,
		Role:         &configuration.roleARN,
		Runtime:      &runtime,
	}
	_, err := svc.CreateFunction(createArgs)
	return err
}

func ConfigLambdaFunctions()  {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := lambda.New(sess)

	err := createFunc( svc, []byte("put"), "Put")
	if err != nil {
		log.Println("Error on creation Put lambda func: ", err)
	}
	err = createFunc( svc, []byte("get"), "Get")
	if err != nil {
		log.Println("Error on creation Get lambda func: ", err)
	}
	err = createFunc( svc, []byte("append"), "Append")
	if err != nil {
		log.Println("Error on creation Append lambda func: ", err)
	}
	err = createFunc( svc, []byte("delete"), "Delete")
	if err != nil {
		log.Println("Error on creation Delete lambda func: ", err)
	}
}


