package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

type RoomStatusResponseStruct struct {
	Result     string       `json:"result"`
	Message    string       `json:"message"`
	RoomStatus RoomStruct   `json:"room_status"`
	RoomLayout RoomLayoutStruct `json:"room_layout"`
	Users      []UserStruct `json:"users"`
}

func RoomStatus(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("RoomStatus()")
	ctx, client := InitializeHttpFuncWithFirestore()
	defer CloseFirestoreClient(client)

	var apiResp RoomStatusResponseStruct
	roomId := request.QueryStringParameters[RoomId]

	if roomId == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else if isExistRoom, _ := IsExistRoom(roomId, client, ctx); !isExistRoom {
		apiResp.Result = ERROR
		apiResp.Message = RoomDoesNotExist
	} else {
		roomLayout, err := RetrieveRoomLayout(roomId, client, ctx)
		if err != nil {
			apiResp.Result = ERROR
			apiResp.Message = "failed to retrieve room layout"
		} else {
			apiResp.RoomLayout = roomLayout

			roomInfo, _ := RetrieveRoomInfo(roomId, client, ctx)
			apiResp.RoomStatus = RoomStruct{
				RoomId: roomId,
				Body:   roomInfo,
			}

			users, _ := RetrieveRoomUsers(roomId, client, ctx)
			apiResp.Users = users
			apiResp.Result = OK
		}
	}
	bytes, _ := json.Marshal(apiResp)
	return Response(bytes)
}

func main() {
	lambda.Start(RoomStatus)
}
