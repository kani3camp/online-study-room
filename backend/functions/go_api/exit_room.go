package go_api

import (
	"encoding/json"
	"net/http"
)

type ExitRoomResponseStruct struct {
	Result string `json:"result"`
	Message string `json:"message"`
}

func ExitRoom(w http.ResponseWriter, r *http.Request) {
	ctx, client := InitializeHttpFunc(&w)
	defer client.Close()
	
	var apiResp ExitRoomResponseStruct
	roomId, userId, idToken := r.FormValue(room_id), r.FormValue(user_id), r.FormValue(id_token)
	
	if roomId == "" || userId == "" || idToken == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else {
		if IsUserVerified(userId, idToken, ctx) {
			statement, err := LeaveRoom(roomId, userId, client, ctx)
			if err != nil {
				apiResp.Result = ERROR
			} else {
				apiResp.Result = OK
			}
			apiResp.Message = statement
		} else {
			apiResp.Result = ERROR
			apiResp.Message = UserAuthFailed
		}
	}
	
	bytes, _ := json.Marshal(apiResp)
	_, _ = w.Write(bytes)
}
