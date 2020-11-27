package go_api

import (
	"encoding/json"
	"net/http"
)

type StayingAwakeResponseStruct struct {
	Result  string       `json:"result"`
	Message string       `json:"message"`
	Users   []UserStruct `json:"users"`
}

func StayingAwake(w http.ResponseWriter, r *http.Request) {
	ctx, client := InitializeHttpFunc(&w)
	defer client.Close()
	
	userId, idToken := r.PostFormValue(user_id), r.PostFormValue(id_token)
	var apiResp StayingAwakeResponseStruct
	
	if userId == "" || idToken == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else if isUserVerified, _ := IsUserVerified(userId, idToken, client, ctx); !isUserVerified {
		apiResp.Result = ERROR
		apiResp.Message = UserAuthFailed
	} else {
		_ = RecordLastAccess(userId, client, ctx)
		roomId, _ := InWhichRoom(userId, client, ctx)
		if roomId == "" {
			apiResp.Result = ERROR
			apiResp.Message = "you are not in the room."
		} else {
			users, _ := RetrieveRoomUsers(roomId, client, ctx)
			apiResp.Result = OK
			apiResp.Users = users
		}
	}
	
	bytes, _ := json.Marshal(apiResp)
	_, _ = w.Write(bytes)
}
