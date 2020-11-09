package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type ExitRoomParams struct {
	RoomId string `json:"room_id"`
	UserId string `json:"user_id"`
	IdToken string `json:"id_token"`
}

type ExitRoomResponseStruct struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

func ExitRoom(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx, client := InitializeHttpFunc()
	defer client.Close()

	var apiResp ExitRoomResponseStruct
	body := request.Body
	params := ExitRoomParams{}
	_ = json.Unmarshal([]byte(body), &params)

	roomId, userId, idToken := params.RoomId, params.UserId, params.IdToken

	if roomId == "" || userId == "" || idToken == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else if isUserVerified, _ := IsUserVerified(userId, idToken, client, ctx); !isUserVerified {
		apiResp.Result = ERROR
		apiResp.Message = UserAuthFailed
	} else if isInRoom, _ := IsInRoom(roomId, userId, client, ctx); !isInRoom {
		apiResp.Result = ERROR
		apiResp.Message = "you are not in the room."
	} else {
		_ = LeaveRoom(roomId, userId, client, ctx)
		apiResp.Result = OK
	}

	bytes, _ := json.Marshal(apiResp)
	return Response(bytes)
}

func main()  {
	lambda.Start(ExitRoom)
}
