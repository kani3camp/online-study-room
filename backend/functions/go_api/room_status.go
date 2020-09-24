package go_api

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type RoomStatusResponseStruct struct {
	Result     string       `json:"result"`
	Message    string       `json:"message"`
	RoomStatus RoomStruct   `json:"room_status"`
	Users      []UserStruct `json:"users"`
}

func RetrieveRoomUsers(roomId string, client *firestore.Client, ctx context.Context) ([]UserStruct, error) {
	var err error
	var users []UserStruct
	
	authClient, _ := InitializeFirebaseAuthClient(ctx)
	
	roomInfo, err := RetrieveRoomInfo(roomId, client, ctx)
	if err != nil {
	} else {
		for _, userId := range roomInfo.Users {
			userBody, err := RetrieveUserInfo(userId, client, ctx)
			if err != nil {
			} else {
				user, _ := authClient.GetUser(ctx, userId)
				users = append(users, UserStruct{
					UserId:      userId,
					DisplayName: user.DisplayName,
					Body:        userBody,
				})
			}
		}
	}
	return users, err
}

func RetrieveRoomInfo(roomId string, client *firestore.Client, ctx context.Context) (RoomBodyStruct, error) {
	var roomBodyStruct RoomBodyStruct
	
	room, err := client.Collection(ROOMS).Doc(roomId).Get(ctx)
	if err != nil {
		log.Println(err)
		return RoomBodyStruct{}, err
	} else {
		_ = room.DataTo(&roomBodyStruct)
		return roomBodyStruct, nil
	}
}

func RoomStatus(w http.ResponseWriter, r *http.Request) {
	ctx, client := InitializeHttpFunc(&w)
	defer client.Close()
	
	var apiResp RoomStatusResponseStruct
	roomId := r.FormValue(room_id)
	
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
	_, _ = w.Write(bytes)
}
