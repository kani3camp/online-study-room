package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type UpdateDatabaseResponseStruct struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

func UpdateDatabase(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx, client := InitializeHttpFunc()
	defer client.Close()

	var apiResp UpdateDatabaseResponseStruct

	_UpdateDatabase(client, ctx)

	apiResp.Result = OK

	jsonBytes, _ := json.Marshal(apiResp)
	return Response(jsonBytes)
}

func main() {
	lambda.Start(UpdateDatabase)
}
