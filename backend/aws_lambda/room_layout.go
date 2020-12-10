package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

type RoomLayoutResponseStruct struct {
	Result         string `json:"result"`
	Message        string `json:"message"`
	RoomLayoutData RoomLayoutStruct `json:"room_layout_data"`
}

func RoomLayout(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("RoomLayout()")
	ctx, client := InitializeEventFuncWithFirestore()
	defer CloseFirestoreClient(client)

	var apiResp RoomLayoutResponseStruct

	roomId := request.QueryStringParameters[RoomId]

	if roomId == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else {
		roomLayout, err := RetrieveRoomLayout(roomId, client, ctx)
		if err != nil {
			apiResp.Result = ERROR
			apiResp.Message = "failed."
		} else {
			apiResp.Result = OK
			apiResp.RoomLayoutData = roomLayout
		}
	}

	bytes, _ := json.Marshal(apiResp)
	return Response(bytes)
}

func main() {
	lambda.Start(RoomLayout)
}
