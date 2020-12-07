package go_api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type ApiResponseStruct struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

// todo websocketのConnectionIdを返す
func EnterRoom(w http.ResponseWriter, r *http.Request) {
	ctx, client := InitializeHttpFuncWithFirestore()
	defer CloseFirestoreClient(client)
	
	var apiResp ApiResponseStruct
	roomId, userId, seatIdStr, idToken := r.PostFormValue(room_id), r.PostFormValue(user_id), r.PostFormValue("seat_id"), r.PostFormValue(id_token)
	seatId, _ := strconv.Atoi(seatIdStr)
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
	log.Println(apiResp)
	bytes, _ := json.Marshal(apiResp)
	_, _ = w.Write(bytes)
}
