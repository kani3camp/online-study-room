package go_api

import (
	"encoding/json"
	"log"
	"net/http"
)

type UserStatusResponseStruct struct {
	Result string `json:"result"`
	Message string `json:"message"`
	UserStatus UserStruct `json:"user_status"`
}


func UserStatus(w http.ResponseWriter, r *http.Request)  {
	ctx, client := InitializeHttpFunc(&w)
	defer client.Close()
	
	userId := r.FormValue(user_id)
	var apiResp UserStatusResponseStruct
	
	if userId == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else {
		userDoc, err := client.Collection(USERS).Doc(userId).Get(ctx)
		
		if err != nil {
			log.Fatalln(err)
		}
		if userDoc.Data() != nil  {
			var body UserBodyStruct
			_ = userDoc.DataTo(&body)
			apiResp.UserStatus = UserStruct {
				UserId: userId,
				Body: body,
			}
			apiResp.Result = OK
		} else {
			apiResp.Result = ERROR
			apiResp.Message = UserDoesNotExist
		}
	}
	
	bytes, _ := json.Marshal(apiResp)
	_, _ = w.Write(bytes)
}
