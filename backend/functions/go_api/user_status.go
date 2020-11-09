package go_api

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type UserStatusResponseStruct struct {
	Result     string     `json:"result"`
	Message    string     `json:"message"`
	UserStatus UserStruct `json:"user_status"`
}

func RetrieveUserInfo(userId string, client *firestore.Client, ctx context.Context) (UserBodyStruct, error) {
	var userBodyStruct UserBodyStruct
	userDoc, err := client.Collection(USERS).Doc(userId).Get(ctx)
	if err != nil {
		log.Println(err)
	} else {
		_ = userDoc.DataTo(&userBodyStruct)
	}
	return userBodyStruct, err
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
		user, _ := authClient.GetUser(ctx, userId)
		
		userInfo, _ := RetrieveUserInfo(userId, client, ctx)
		apiResp.UserStatus = UserStruct{
			UserId:      userId,
			DisplayName: user.DisplayName,
			Body:        userInfo,
		}
		apiResp.Result = OK
	}
	
	bytes, _ := json.Marshal(apiResp)
	_, _ = w.Write(bytes)
}
