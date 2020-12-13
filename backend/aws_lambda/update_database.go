package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

type UpdateDatabaseResponseStruct struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

func UpdateDatabase(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("UpdateDatabase()")
	ctx, client := InitializeHttpFuncWithFirestore()
	defer CloseFirestoreClient(client)

	var apiResp UpdateDatabaseResponseStruct

	_ = _UpdateDatabase(client, ctx)

	apiResp.Result = OK

	jsonBytes, _ := json.Marshal(apiResp)
	return Response(jsonBytes)
}

func main() {
	lambda.Start(UpdateDatabase)
}
