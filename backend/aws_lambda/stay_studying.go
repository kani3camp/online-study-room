package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
)

type StayStudyingRequest struct {
	UserId string `json:"user_id"`
	IdToken string `json:"id_token"`
	RoomId string `json:"room_id"`
	DeviceType string `json:"device_type"`
}

type StayStudyingResponse struct {
	IsOk bool `json:"is_ok"`
	Message string `json:"message"`
	Users   []UserStruct `json:"users"`
}

func StayStudying(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	_, client := InitializeEventFunc()
	defer client.Close()
	
	var requestData StayStudyingRequest
	_ = json.Unmarshal([]byte(request.Body), &requestData)
	
	var response StayStudyingResponse
	
	if isInRoom, _ := IsInRoom(requestData.RoomId, requestData.UserId, client, ctx); !isInRoom {
		response.IsOk = false
		response.Message = "you are not in the room."
	} else {
		_ = RecordLastAccess(requestData.UserId, client, ctx)
		
		users, _ := RetrieveRoomUsers(requestData.RoomId, client, ctx)
		
		response.IsOk = true
		response.Users = users
	}
	
	res , _ := json.Marshal(response)
	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: string(res)}, nil
}


func main()  {
	lambda.Start(StayStudying)
}