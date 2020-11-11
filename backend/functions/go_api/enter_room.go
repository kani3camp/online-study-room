package go_api

import (
	"encoding/json"
	"log"
	"net/http"
)

type ApiResponseStruct struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

func EnterRoom(w http.ResponseWriter, r *http.Request) {
	ctx, client := InitializeHttpFunc(&w)
	defer client.Close()
	
	var apiResp ApiResponseStruct
	roomId, userId, idToken := r.PostFormValue(room_id), r.PostFormValue(user_id), r.PostFormValue(id_token)
	
	if roomId == "" || userId == "" || idToken == "" {
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
				client.Close()
				
				client, _ = InitializeFirestoreClient(ctx)
				_ = _EnterRoom(roomId, userId, client, ctx)
				apiResp.Result = OK
				apiResp.Message = "successfully entered " + roomId + "."
			}
		} else {
			// 入室処理
			_ = _EnterRoom(roomId, userId, client, ctx)
			apiResp.Result = OK
			apiResp.Message = "successfully entered " + roomId + "."
		}
	}
	log.Println(apiResp)
	bytes, _ := json.Marshal(apiResp)
	_, _ = w.Write(bytes)
}
