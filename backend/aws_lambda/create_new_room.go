package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"os"
)

type CreateNewRoomParams struct {
	RoomId        string `json:"room_id"`
	RoomName      string `json:"room_name"`
	RoomType      string `json:"room_type"`
	Password      string `json:"password"`
	ThemeColorHex string `json:"theme_color_hex"`
}

type CreateNewRoomResponseStruct struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

func CreateNewRoom(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("CreateNewRoom()")
	ctx, client := InitializeHttpFuncWithFirestore()
	defer CloseFirestoreClient(client)

	var apiResp CreateNewRoomResponseStruct
	body := request.Body
	params := CreateNewRoomParams{}
	_ = json.Unmarshal([]byte(body), &params)

	roomId, roomName, roomType, themeColorHex := params.RoomId, params.RoomName, params.RoomType, params.ThemeColorHex
	password := params.Password

	if roomId == "" || roomName == "" || roomType == "" || password == "" || themeColorHex == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else if password != os.Getenv("password") {
		apiResp.Result = ERROR
		apiResp.Message = "Invalid password."
	} else {
		continueFlag := true
		// 同じ部屋がすでにないかチェック
		rooms, _ := RetrieveRooms(client, ctx)
		for _, room := range rooms {
			if room.RoomId == roomId {
				apiResp.Result = ERROR
				apiResp.Message = "The room id is already used."
				continueFlag = false
				break
			}
		}
		if continueFlag {
			_ = _CreateNewRoom(roomId, roomName, roomType, themeColorHex, client, ctx)
			apiResp.Result = OK
			apiResp.Message = "Successfully created room named " + roomId
		}
	}

	bytes, _ := json.Marshal(apiResp)
	return Response(bytes)
}

func main() {
	lambda.Start(CreateNewRoom)
}
