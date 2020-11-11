package go_api

import (
	"encoding/json"
	"net/http"
	"os"
)

type CreateNewRoomResponseStruct struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

func CreateNewRoom(w http.ResponseWriter, r *http.Request) {
	ctx, client := InitializeHttpFunc(&w)
	defer client.Close()
	
	var apiResp ChangeUserInfoResponseStruct
	roomId, roomName := r.PostFormValue(room_id), r.PostFormValue("room_name")
	roomType, themeColorHex := r.PostFormValue("room_type"), r.PostFormValue("theme_color_hex")
	password := r.PostFormValue("password")
	
	if roomId == "" || roomName == "" || roomType == "" || password == "" || themeColorHex == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else
	if password != os.Getenv("password") {
		apiResp.Result = ERROR
		apiResp.Message = "Invalid password."
	} else {
		continueFlag := true
		// 同じ部屋がすでにないかチェック
		rooms, _ := RetrieveRooms(client, ctx)
		for _, room := range rooms {
			if room.RoomId == roomId {
				apiResp.Result = ERROR
				apiResp.Message = "The room id is already used."
				continueFlag = false
				break
			}
		}
		if continueFlag {
			_ = _CreateNewRoom(roomId, roomName, roomType, themeColorHex, client, ctx)
			apiResp.Result = OK
			apiResp.Message = "Successfully created room named " + roomId
		}
	}
	
	bytes, _ := json.Marshal(apiResp)
	_, _ = w.Write(bytes)
}
