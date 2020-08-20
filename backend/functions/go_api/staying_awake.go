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
					log.Fatalln(err)
				} else {
					apiResp.Result = OK
					apiResp.Message = "Successfully updated your last-access."
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
