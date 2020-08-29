package go_api

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type ApiResponseStruct struct {
	Result string `json:"result"`
	Message string `json:"message"`
}

func enterRoom(roomId string, userId string, client *firestore.Client, ctx context.Context) (*firestore.WriteResult, error) {
	return client.Collection(ROOMS).Doc(roomId).Set(ctx, map[string]interface{}{
		"users": firestore.ArrayUnion(userId),
	}, firestore.MergeAll)
}

func EnterRoom(w http.ResponseWriter, r *http.Request)  {
	ctx, client := InitializeHttpFunc(&w)
	defer client.Close()
	
	var apiResponse ApiResponseStruct
	roomId, userId, idToken := r.FormValue(room_id), r.FormValue(user_id), r.FormValue(id_token)
	
	if roomId == "" || userId == "" || idToken == "" {
		apiResponse.Result = ERROR
		apiResponse.Message = InvalidParams
	} else {
		// auth
		if IsUserVerified(userId, idToken, ctx) {
			if IsExistRoom(roomId, client, ctx) {
				if IsInUsers(userId, client, ctx) {
					if IsOnline(userId, client, ctx) {
						currentRoomId, err := InWhichRoom(userId, client, ctx)
						if err != nil {
							log.Println(err)
							apiResponse.Result = ERROR
							apiResponse.Message = "Failed InWhichRoom()"
						} else {
							if currentRoomId == roomId {
								apiResponse.Result = OK
								apiResponse.Message = "You are already in the " + currentRoomId
							} else {
								_, _ = LeaveRoom(currentRoomId, userId, client, ctx)
								_, err := enterRoom(roomId, userId, client, ctx)
								if err != nil {
									log.Println(err)
									apiResponse.Result = ERROR
									apiResponse.Message = "Failed enterRoom()"
								} else {
									apiResponse.Result = OK
									apiResponse.Message = "Successfully entered " + roomId + "."
								}
							}
						}
					} else {
						// 入室処理
						_, err := enterRoom(roomId, userId, client, ctx)
						if err != nil {
							log.Println(err)
							apiResponse.Result = ERROR
							apiResponse.Message = "Failed enterRoom()"
						} else {
							apiResponse.Result = OK
							apiResponse.Message = "Successfully entered " + roomId + "."
						}
					}
				} else {
					apiResponse.Result = ERROR
					apiResponse.Message = InvalidUser
				}
			} else {
				apiResponse.Result = ERROR
				apiResponse.Message = RoomDoesNotExist
			}
		} else {
			apiResponse.Result = ERROR
			apiResponse.Message = UserAuthFailed
		}
	}
	log.Println(apiResponse)
	bytes, _ := json.Marshal(apiResponse)
	_, _ = w.Write(bytes)
}
