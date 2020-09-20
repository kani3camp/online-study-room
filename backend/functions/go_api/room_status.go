package go_api

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type RoomStatusResponseStruct struct {
	Result string `json:"result"`
	Message string `json:"message"`
	RoomStatus RoomStruct `json:"room_status"`
	Users []UserStruct `json:"users"`
}

func GetRoomUsers(roomId string, client *firestore.Client, ctx context.Context) ([]UserStruct, error) {
	var err error
	var users []UserStruct

	app, _ := InitializeFirebaseApp(ctx)
	authClient, _ := app.Auth(ctx)

	roomInfo, err := GetRoomInfo(roomId, client, ctx)
	if err != nil {
	} else {
		for _, userId := range roomInfo.Users {
			userBody, err := GetUserInfo(userId, client, ctx)
			if err != nil {
			} else {
				user, _ := authClient.GetUser(ctx, userId)
				users = append(users, UserStruct{
					UserId: userId,
					DisplayName: user.DisplayName,
					Body: userBody,
				})
			}
		}
	}
	return users, err
}

func GetRoomInfo(roomId string, client *firestore.Client, ctx context.Context) (RoomBodyStruct, error) {
	room, err := client.Collection(ROOMS).Doc(roomId).Get(ctx)
	var roomBodyStruct RoomBodyStruct
	if err != nil {
		log.Println(err)
	} else {
		_ = room.DataTo(&roomBodyStruct)
	}
	return roomBodyStruct, err
}

func RoomStatus(w http.ResponseWriter, r *http.Request) {
	ctx, client := InitializeHttpFunc(&w)
	defer client.Close()
	
	var apiResp RoomStatusResponseStruct
	roomId := r.FormValue(room_id)
	
	if roomId == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else if IsExistRoom(roomId, client, ctx) {
		roomInfo, err := GetRoomInfo(roomId, client, ctx)
		if err != nil {
			apiResp.Result = ERROR
			apiResp.Message = err.Error()
		} else {
			apiResp.RoomStatus = RoomStruct{
				RoomId: roomId,
				Body:   roomInfo,
			}
			
			users, err := GetRoomUsers(roomId, client, ctx)
			if err != nil {
				apiResp.Result = ERROR
				apiResp.Message = err.Error()
			} else {
				apiResp.Users = users
				apiResp.Result = OK
			}
		}
	} else {
		log.Println(RoomDoesNotExist)
		apiResp.Result = ERROR
		apiResp.Message = RoomDoesNotExist
	}
	bytes, _ := json.Marshal(apiResp)
	_, _ = w.Write(bytes)
}
