package go_api

import (
	"encoding/json"
	"net/http"
)

type ExitRoomResponseStruct struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

func ExitRoom(w http.ResponseWriter, r *http.Request) {
	ctx, client := InitializeHttpFuncWithFirestore()
	defer CloseFirestoreClient(client)
	
	var apiResp ExitRoomResponseStruct
	roomId, userId, idToken := r.FormValue(room_id), r.FormValue(user_id), r.FormValue(id_token)
	
	if roomId == "" || userId == "" || idToken == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else if isUserVerified, _ := IsUserVerified(userId, idToken, client, ctx); !isUserVerified {
		apiResp.Result = ERROR
		apiResp.Message = UserAuthFailed
	} else if isInRoom, _ := IsInRoom(roomId, userId, client, ctx); !isInRoom {
		apiResp.Result = ERROR
		apiResp.Message = "you are not in the room"
	} else {
		err := LeaveRoom(roomId, userId, client, ctx)
		if err != nil {
			apiResp.Result = ERROR
			apiResp.Message = "failed to leave successfully"
		} else {
			apiResp.Result = OK
		}
	}
	
	bytes, _ := json.Marshal(apiResp)
	_, _ = w.Write(bytes)
}
