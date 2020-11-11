package go_api

import (
	"encoding/json"
	"net/http"
)

type RoomStatusResponseStruct struct {
	Result     string       `json:"result"`
	Message    string       `json:"message"`
	RoomStatus RoomStruct   `json:"room_status"`
	Users      []UserStruct `json:"users"`
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
