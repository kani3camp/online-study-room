package main

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)


func EndStudying(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) error {
	_, client := InitializeEventFuncWithFirestore()
	defer CloseFirestoreClient(client)

	connectionId := request.RequestContext.ConnectionID

	userInfo, err := FindUserWithConnectionId(connectionId, client, ctx)
	if err != nil {
		log.Println(err)
		return err
	} else if isInRoom, _ := IsInRoom(userInfo.Body.In, userInfo.UserId, client, ctx); !isInRoom {
		errString := "this user(" + userInfo.UserId + ") is not in any room"
		return errors.New(errString)
	}
	err = LeaveRoom(userInfo.Body.In, userInfo.UserId, client, ctx)
	return err
}

func main() {
	lambda.Start(EndStudying)
}
