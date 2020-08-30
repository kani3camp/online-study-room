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
	UserNames []string `json:"user_names"`
}

func GetRoomUserNames(roomId string, client *firestore.Client, ctx context.Context) ([]string, error) {
	var err error
	var userNames []string
	roomInfo, err := GetRoomInfo(roomId, client, ctx)
	if err != nil {
	} else {
		for _, id := range roomInfo.Users {
			userInfo, err := GetUserInfo(id, client, ctx)
			if err != nil {
			} else {
				userNames = append(userNames, userInfo.Name)
			}
		}
	}
	return userNames, err
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
			
			userNames, err := GetRoomUserNames(roomId, client, ctx)
			if err != nil {
				apiResp.Result = ERROR
				apiResp.Message = err.Error()
			} else {
				apiResp.UserNames = userNames
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
