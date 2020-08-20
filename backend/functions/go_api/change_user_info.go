package go_api

import (
	"cloud.google.com/go/firestore"
	"encoding/json"
	"log"
	"net/http"
)

type ChangeUserInfoResponseStruct struct {
	Result string `json:"result"`
	Message string `json:"message"`
}

func ChangeUserInfo(w http.ResponseWriter, r *http.Request)  {
	ctx, client := InitializeHttpFunc(&w)
	defer client.Close()
	
	var apiResponse ChangeUserInfoResponseStruct
	userId, idToken := r.FormValue(user_id), r.FormValue(id_token)
	displayName, statusMessage := r.FormValue("display_name"), r.FormValue("status_message")
	
	if userId == "" || idToken == "" || displayName == "" || statusMessage == "" {
		apiResponse.Result = ERROR
		apiResponse.Message = InvalidParams
	} else {
		if IsUserVerified(userId, idToken, ctx) {
			_, err := client.Collection(USERS).Doc(userId).Set(ctx, map[string]interface{}{
				"last-access": firestore.ServerTimestamp,
				"name": displayName,
				"status": statusMessage,
			}, firestore.MergeAll)
			if err != nil {
				log.Fatalln("Failed to update user info of " + userId)
			}
			apiResponse.Result = OK
		} else {
			apiResponse.Result = ERROR
			apiResponse.Message = UserAuthFailed
		}
	}
	
	bytes, _ := json.Marshal(apiResponse)
	_, _ = w.Write(bytes)
}
