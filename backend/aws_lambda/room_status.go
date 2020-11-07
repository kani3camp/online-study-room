package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RoomStatusResponseStruct struct {
	Result     string       `json:"result"`
	Message    string       `json:"message"`
	RoomStatus RoomStruct   `json:"room_status"`
	Users      []UserStruct `json:"users"`
}

func RoomStatus(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx, client := InitializeHttpFunc()
	defer client.Close()

	var apiResp RoomStatusResponseStruct
	roomId := request.QueryStringParameters[room_id]

	if roomId == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else if isExistRoom, _ := IsExistRoom(roomId, client, ctx); !isExistRoom {
		apiResp.Result = ERROR
		apiResp.Message = RoomDoesNotExist
	} else {
		roomInfo, _ := RetrieveRoomInfo(roomId, client, ctx)
		apiResp.RoomStatus = RoomStruct{
			RoomId: roomId,
			Body:   roomInfo,
		}

		users, _ := RetrieveRoomUsers(roomId, client, ctx)
		apiResp.Users = users
		apiResp.Result = OK
	}
	bytes, _ := json.Marshal(apiResp)
	return Response(bytes)
}

func main() {
	lambda.Start(RoomStatus)
}
