package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

type OnlineUsersResponse struct {
	Result      string       `json:"result"`
	Message     string       `json:"message"`
	OnlineUsers []UserStruct `json:"online_users"`
}

func OnlineUsers(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("OnlineUsers()")
	ctx, client := InitializeHttpFuncWithFirestore()
	defer CloseFirestoreClient(client)

	var apiResp OnlineUsersResponse
	onlineUsers, _ := RetrieveOnlineUsersAsStatus(client, ctx)
	apiResp.Result = OK
	apiResp.OnlineUsers = onlineUsers

	bytes, _ := json.Marshal(apiResp)
	return Response(bytes)
}

func main() {
	lambda.Start(OnlineUsers)
}
