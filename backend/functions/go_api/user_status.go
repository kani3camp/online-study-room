package go_api

import (
	"encoding/json"
	"log"
	"net/http"
)

type UserStatusResponseStruct struct {
	Result     string     `json:"result"`
	Message    string     `json:"message"`
	UserStatus UserStruct `json:"user_status"`
}

func UserStatus(w http.ResponseWriter, r *http.Request) {
	ctx, client := InitializeHttpFunc(&w)
	defer client.Close()
	
	userId := r.FormValue(user_id)
	var apiResp UserStatusResponseStruct
	
	if userId == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else if isInUsers, _ := IsInUsers(userId, client, ctx); !isInUsers {
		apiResp.Result = ERROR
		apiResp.Message = InvalidUser
	} else {
		authClient, _ := InitializeFirebaseAuthClient(ctx)
		var displayName string
		user, err := authClient.GetUser(ctx, userId)
		if err != nil {
			// テストユーザーだとauthに登録されてないかもなので起こりうる
			log.Println("faield authClient.GetUser(ctx, userId).")
			displayName = ""
		} else {
			displayName = user.DisplayName
		}

		log.Println("p4")
		userInfo, _ := RetrieveUserInfo(userId, client, ctx)
		apiResp.UserStatus = UserStruct{
			UserId:      userId,
			DisplayName: displayName,
			Body:        userInfo,
		}
		apiResp.Result = OK
	}
	
	bytes, _ := json.Marshal(apiResp)
	_, _ = w.Write(bytes)
}
