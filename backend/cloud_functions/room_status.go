package go_api

import (
	"encoding/json"
	"net/http"
)

type RoomStatusResponseStruct struct {
	Result     string       `json:"result"`
	Message    string       `json:"message"`
	RoomStatus RoomStruct   `json:"room_status"`
	RoomLayout RoomLayoutStruct `json:"room_layout"`
	Users      []UserStruct `json:"users"`
}

func RoomStatus(w http.ResponseWriter, r *http.Request) {
	ctx, client := InitializeHttpFuncWithFirestore()
	defer CloseFirestoreClient(client)
	
	var apiResp RoomStatusResponseStruct
	roomId := r.FormValue(room_id)
	
	if roomId == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else if isExistRoom, _ := IsExistRoom(roomId, client, ctx); !isExistRoom {
		apiResp.Result = ERROR
		apiResp.Message = RoomDoesNotExist
	} else {
		roomLayout, err := RetrieveRoomLayout(roomId, client, ctx)
		if err != nil {
			apiResp.Result = ERROR
			apiResp.Message = "failed to retrieve room layout"
		} else {
			apiResp.RoomLayout = roomLayout

			roomInfo, _ := RetrieveRoomInfo(roomId, client, ctx)
			apiResp.RoomStatus = RoomStruct{
				RoomId: roomId,
				Body:   roomInfo,
			}

			users, _ := RetrieveRoomUsers(roomId, client, ctx)
			apiResp.Users = users
			apiResp.Result = OK
		}
	}
	bytes, _ := json.Marshal(apiResp)
	_, _ = w.Write(bytes)
}
