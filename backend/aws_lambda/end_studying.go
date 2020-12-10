package main

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)


// AWSの決まりとして、必ずerrorを返す必要あり
func EndStudying(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) error {
	log.Println("EndStudying()")
	_, client := InitializeEventFuncWithFirestore()
	defer CloseFirestoreClient(client)

	connectionId := request.RequestContext.ConnectionID
	log.Printf("connectionId: %s\n", connectionId)

	userInfo, err := FindUserWithConnectionId(connectionId, client, ctx)
	if err.Body != nil {
		log.Println(err)
		return err.Body
	} else if isInRoom, _, _ := IsInRoom(userInfo.Body.In, userInfo.UserId, client, ctx); !isInRoom {
		errString := "this user(" + userInfo.UserId + ") is not in any room"
		return errors.New(errString)
	}
	err = LeaveRoom(userInfo.Body.In, userInfo.UserId, client, ctx)
	return err.Body
}

func main() {
	lambda.Start(EndStudying)
}
