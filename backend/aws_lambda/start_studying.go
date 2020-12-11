package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
)

type StartStudyingRequest struct {
	UserId     string `json:"user_id"`
	IdToken    string `json:"id_token"`
	RoomId     string `json:"room_id"`
	SeatId int `json:"seat_id"`
}

type StartStudyingResponse struct {
	IsOk    bool         `json:"is_ok"`
	Message string       `json:"message"`
	Users   []UserStruct `json:"users"`
	RoomInfo RoomBodyStruct `json:"room_info"`
	RoomLayout RoomLayoutStruct `json:"room_layout"`
}

func StartStudying(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("StartStudying()")
	_, client := InitializeEventFuncWithFirestore()
	defer CloseFirestoreClient(client)

	var requestData StartStudyingRequest
	_ = json.Unmarshal([]byte(request.Body), &requestData)
	roomId, userId, seatId, idToken := requestData.RoomId, requestData.UserId, requestData.SeatId, requestData.IdToken
	connectionId := request.RequestContext.ConnectionID

	log.Printf("request: %v", request)
	log.Println("requested user id is ", userId)
	var response StartStudyingResponse

	if roomId == "" || userId == "" || seatId == 0 || idToken == "" {
		response.IsOk = false
		response.Message = InvalidParams
	} else if isUserVerified, _ := IsUserVerified(userId, idToken, client, ctx); !isUserVerified {
		response.IsOk = false
		response.Message = UserAuthFailed
	} else if isExistRoom, _ := IsExistRoom(roomId, client, ctx); !isExistRoom {
		response.IsOk = false
		response.Message = RoomDoesNotExist
	} else if isInUsers, _ := IsInUsers(userId, client, ctx); !isInUsers {
		response.IsOk = false
		response.Message = InvalidUser
	} else if isOnline, _ := IsOnline(userId, client, ctx); isOnline {
		currentRoomId, _ := InWhichRoom(userId, client, ctx)
		currentSeatId, _ := RetrieveCurrentSeatId(userId, client, ctx)
		// すでにその部屋のその席にいた場合
		if currentRoomId == roomId && currentSeatId == seatId {
			// 他のconnection idがすでに設定されている場合、強制的に書き換える
			response.IsOk = true
			response.Message = "you are already in the " + currentRoomId
		} else {
			// 一旦退室させてから入室
			_ = LeaveRoom(currentRoomId, userId, client, ctx)
			_ = client.Close()

			client, _ = InitializeFirestoreClient(ctx)
			err := _EnterRoom(roomId, userId, seatId, client, ctx)
			if err.Body != nil {
				response.IsOk = false
				switch err.ErrorType {
				case Unknown:
					response.Message = "failed to enter room"
				default:
					response.Message = err.Body.Error()
				}
			} else {
				response.IsOk = true
				response.Message = "successfully entered " + roomId + "."
			}
		}
	} else {
		// 入室処理
		err := _EnterRoom(roomId, userId, seatId, client, ctx)
		if err.Body != nil {
			log.Println(err.ErrorType)
			response.IsOk = false
			switch err.ErrorType {
			case Unknown:
				response.Message = "failed to enter room"
			default:
				response.Message = err.Body.Error()
			}
		} else {
			response.IsOk = true
			response.Message = "successfully entered " + roomId + "."
		}
	}

	if response.IsOk {
		// connection idを設定する
		err := SetConnectionId(userId, connectionId, client, ctx)
		if err != nil {
			response.IsOk = false
			response.Message = "failed to set the connection id"
		} else {
			users, _ := RetrieveRoomUsers(roomId, client, ctx)
			response.Users = users
			roomInfo, _ := RetrieveRoomInfo(roomId, client, ctx)
			response.RoomInfo = roomInfo
			roomLayout, _ := RetrieveRoomLayout(roomId, client, ctx)
			response.RoomLayout = roomLayout.SetIsVacant(client, ctx)
		}
	} else {
		// 切断
		_ = Disconnect(connectionId, client, ctx)
	}

	log.Printf("response: %v\n", response)

	res, _ := json.Marshal(response)
	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: string(res)}, nil
}

func main() {
	lambda.Start(StartStudying)
}
