package go_api

import (
	"cloud.google.com/go/firestore"
	"encoding/json"
	"log"
	"net/http"
)

type StayingAwakeResponseStruct struct {
	Result string `json:"result"`
	Message string `json:"message"`
	Users []UserStruct `json:"users"`
}

func StayingAwake(w http.ResponseWriter, r *http.Request)  {
	ctx, client := InitializeHttpFunc(&w)
	defer client.Close()
	
	userId, idToken := r.FormValue(user_id), r.FormValue(id_token)
	var apiResp StayingAwakeResponseStruct
	
	if userId == "" || idToken == ""{
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else {
		if IsUserVerified(userId, idToken, ctx) {
			if IsInUsers(userId, client, ctx) {
				_, err := client.Collection(USERS).Doc(userId).Set(ctx, map[string]interface{}{
					"last-access": firestore.ServerTimestamp,
				}, firestore.MergeAll)
				if err != nil {
					log.Println(err)
					apiResp.Result = ERROR
					apiResp.Message = "failed to refresh your last-access"
				} else {
					userBody, err := GetUserInfo(userId, client, ctx)
					if err != nil {
						apiResp.Result = ERROR
						apiResp.Message = "failed to retrieve user info"
					} else {
						if userBody.In == "" {
							apiResp.Result = ERROR
							apiResp.Message = "failed to retrieve users in the room"
						} else {
							users, err := GetRoomUsers(userBody.In, client, ctx)
							if err != nil {
								apiResp.Result = ERROR
								apiResp.Message = "failed to retrieve users in the room"
							} else {
								apiResp.Result = OK
								apiResp.Users = users
							}
						}
					}
				}
			} else {
				apiResp.Result = ERROR
				apiResp.Message = InvalidUser
			}
		} else {
			apiResp.Result = ERROR
			apiResp.Message = UserAuthFailed
		}
	}
	
	bytes, _ := json.Marshal(apiResp)
	_, _ = w.Write(bytes)
}
