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

// todo connection idも考慮に入れる
func UpdateDatabase(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx, client := InitializeHttpFuncWithFirestore()
	defer CloseFirestoreClient(client)

	var apiResp UpdateDatabaseResponseStruct

	_UpdateDatabase(client, ctx)

	apiResp.Result = OK

	jsonBytes, _ := json.Marshal(apiResp)
	return Response(jsonBytes)
}

func main() {
	lambda.Start(UpdateDatabase)
}
