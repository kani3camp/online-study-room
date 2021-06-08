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
	RoomLayout RoomLayoutStruct `json:"room_layout"`
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
	roomLayout, password := params.RoomLayout, params.Password

	if roomId == "" || roomName == "" || roomType == "" || password == "" || themeColorHex == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else if password != os.Getenv("password") {
		apiResp.Result = ERROR
		apiResp.Message = "invalid password"
	} else {
		continueFlag := true
		// 同じ部屋がすでにないかチェック
		rooms, _ := RetrieveRooms(client, ctx)
		for _, room := range rooms {
			if room.RoomId == roomId {
				apiResp.Result = ERROR
				apiResp.Message = "the room id is already used"
				continueFlag = false
				break
			}
		}
		if continueFlag {
			// roomLayoutが有効かチェック
			if roomId != roomLayout.RoomId {
				apiResp.Result = ERROR
				apiResp.Message = "room layout's room id is not the same as what you specified"
			} else {
				customErr := CheckRoomLayoutData(roomLayout, client, ctx)
				if customErr.Body != nil {
					apiResp.Result = ERROR
					apiResp.Message = customErr.Body.Error()
				} else {
					err := _CreateNewRoom(roomId, roomName, roomType, themeColorHex, client, ctx)
					if err != nil {
						apiResp.Result = ERROR
						apiResp.Message = err.Error()
					} else {
						err := SaveRoomLayout(roomLayout, client, ctx)
						if err != nil {
							apiResp.Result = ERROR
							apiResp.Message = "failed. please review the log in CloudWatch."
						} else {
							apiResp.Result = OK
							apiResp.Message = "successfully created room named " + roomId
						}
					}
				}
			}
		}
	}

	bytes, _ := json.Marshal(apiResp)
	return Response(bytes)
}

func main() {
	lambda.Start(CreateNewRoom)
}
