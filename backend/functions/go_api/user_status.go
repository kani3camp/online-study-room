package go_api

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type UserStatusResponseStruct struct {
	Result string `json:"result"`
	Message string `json:"message"`
	UserStatus UserStruct `json:"user_status"`
}

func GetUserInfo(userId string, client *firestore.Client, ctx context.Context) (UserBodyStruct, error) {
	var userBodyStruct UserBodyStruct
	userDoc, err := client.Collection(USERS).Doc(userId).Get(ctx)
	if err != nil {
		log.Println(err)
	} else {
		_ = userDoc.DataTo(&userBodyStruct)
	}
	return userBodyStruct, err
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
		// todo userが存在するか？
		userInfo, err := GetUserInfo(userId, client, ctx)
		if err != nil {
			log.Println(err)
			apiResp.Result = ERROR
			apiResp.Message = Failed
		} else {
			apiResp.UserStatus = UserStruct{
				UserId: userId,
				Body:   userInfo,
			}
			apiResp.Result = OK
		}
	}
	
	bytes, _ := json.Marshal(apiResp)
	_, _ = w.Write(bytes)
}
