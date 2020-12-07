package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
)

type StayStudyingRequest struct {
	UserId     string `json:"user_id"`
	IdToken    string `json:"id_token"`
	RoomId     string `json:"room_id"`
	//SeatId int `json:"seat_id"` 部屋にいることが確認できればいいや
}

type StayStudyingResponse struct {
	IsOk    bool         `json:"is_ok"`
	Message string       `json:"message"`
	Users   []UserStruct `json:"users"`
	RoomInfo RoomBodyStruct `json:"room_info"`
}

func StayStudying(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	_, client := InitializeEventFuncWithFirestore()
	defer CloseFirestoreClient(client)

	var requestData StayStudyingRequest
	_ = json.Unmarshal([]byte(request.Body), &requestData)
	userId, idToken, roomId := requestData.UserId, requestData.IdToken, requestData.RoomId

	connectionId := request.RequestContext.ConnectionID
	log.Printf("request: %v", request)

	var response StayStudyingResponse

	if roomId == "" || userId == "" || idToken == "" {
		response.IsOk = false
		response.Message = InvalidParams
	} else if isUserVerified, _ := IsUserVerified(userId, idToken, client, ctx); !isUserVerified {
		response.IsOk = false
		response.Message = UserAuthFailed
	} else {
		// firestoreで他のconnection idになっていたり、connection idが指定されていなかったら切断
		currentConnectionId, _ := RetrieveConnectionId(userId, client, ctx)
		if connectionId != currentConnectionId {
			response.IsOk = false
			response.Message = "this connection id is invalid"
		} else if isInRoom, _ := IsInRoom(requestData.RoomId, requestData.UserId, client, ctx); !isInRoom {
			response.IsOk = false
			response.Message = "you are not in the room."
		} else {
			response.IsOk = true
		}
	}

	if response.IsOk {
		_ = RecordLastAccess(userId, client, ctx)
		users, _ := RetrieveRoomUsers(requestData.RoomId, client, ctx)
		response.Users = users
		roomInfo, _ := RetrieveRoomInfo(roomId, client, ctx)
		response.RoomInfo = roomInfo
	} else {
		// 切断
		_ = Disconnect(connectionId, client, ctx)
	}

	res, _ := json.Marshal(response)
	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: string(res)}, nil
}

func main() {
	lambda.Start(StayStudying)
}
