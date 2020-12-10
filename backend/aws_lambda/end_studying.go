package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)


func EndStudying(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) CustomError {
	log.Println("EndStudying")
	_, client := InitializeEventFuncWithFirestore()
	defer CloseFirestoreClient(client)

	connectionId := request.RequestContext.ConnectionID

	userInfo, err := FindUserWithConnectionId(connectionId, client, ctx)
	if err.Body != nil {
		log.Println(err)
		return err
	} else if isInRoom, _ := IsInRoom(userInfo.Body.In, userInfo.UserId, client, ctx); !isInRoom {
		errString := "this user(" + userInfo.UserId + ") is not in any room"
		return UserNotInAnyRoom.New(errString)
	}
	err = LeaveRoom(userInfo.Body.In, userInfo.UserId, client, ctx)
	return err
}

func main() {
	lambda.Start(EndStudying)
}
