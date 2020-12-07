package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type EnterRoomParams struct {
	RoomId  string `json:"room_id"`
	UserId  string `json:"user_id"`
	SeatId int `json:"seat_id"`
	IdToken string `json:"id_token"`
}

type EnterRoomResponseStruct struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

// todo EnterRoom関数いる？いるなら席も指定できるよう修正
func EnterRoom(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx, client := InitializeHttpFuncWithFirestore()
	defer CloseFirestoreClient(client)

	var apiResp EnterRoomResponseStruct
	body := request.Body
	params := EnterRoomParams{}
	_ = json.Unmarshal([]byte(body), &params)

	roomId, userId, seatId, idToken := params.RoomId, params.UserId, params.SeatId, params.IdToken

	if roomId == "" || userId == "" || seatId == 0 || idToken == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else {
		// auth
		if isUserVerified, _ := IsUserVerified(userId, idToken, client, ctx); !isUserVerified {
			apiResp.Result = ERROR
			apiResp.Message = UserAuthFailed
		} else if isExistRoom, _ := IsExistRoom(roomId, client, ctx); !isExistRoom {
			apiResp.Result = ERROR
			apiResp.Message = RoomDoesNotExist
		} else if isInUsers, _ := IsInUsers(userId, client, ctx); !isInUsers {
			apiResp.Result = ERROR
			apiResp.Message = InvalidUser
		} else if isOnline, _ := IsOnline(userId, client, ctx); isOnline {
			// 一旦退室させてから入室
			currentRoomId, _ := InWhichRoom(userId, client, ctx)
			if currentRoomId == roomId {
				apiResp.Result = OK
				apiResp.Message = "you are already in the " + currentRoomId
			} else {
				_ = LeaveRoom(currentRoomId, userId, client, ctx)
				_ = client.Close()

				client, _ = InitializeFirestoreClient(ctx)
				err := _EnterRoom(roomId, userId, seatId, client, ctx)
				if err != nil {
					apiResp.Result = ERROR
					apiResp.Message = "failed to enter room"
				} else {
					apiResp.Result = OK
					apiResp.Message = "successfully entered " + roomId + "."
				}
			}
		} else {
			// 入室処理
			err := _EnterRoom(roomId, userId, seatId, client, ctx)
			if err != nil {
				apiResp.Result = ERROR
				apiResp.Message = "failed to enter room"
			} else {
				apiResp.Result = OK
				apiResp.Message = "successfully entered " + roomId + "."
			}
		}
	}

	bytes, _ := json.Marshal(apiResp)
	return Response(bytes)
}

func main() {
	lambda.Start(EnterRoom)
}
