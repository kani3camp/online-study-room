package go_api

import (
	"encoding/json"
	"firebase.google.com/go/auth"
	"net/http"
)

type ChangeUserInfoResponseStruct struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

func ChangeUserInfo(w http.ResponseWriter, r *http.Request) {
	ctx, client := InitializeHttpFuncWithFirestore()
	defer CloseFirestoreClient(client)
	
	var apiResp ChangeUserInfoResponseStruct
	userId, idToken := r.FormValue(user_id), r.FormValue(id_token)
	displayName, statusMessage := r.FormValue("display_name"), r.FormValue("status_message")
	
	if userId == "" || idToken == "" || displayName == "" || statusMessage == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else if _isUserVerified, _ := IsUserVerified(userId, idToken, client, ctx); !_isUserVerified {
		apiResp.Result = ERROR
		apiResp.Message = UserAuthFailed
	} else {
		authClient, _ := InitializeFirebaseAuthClient(ctx)
		params := (&auth.UserToUpdate{}).DisplayName(displayName)
		_, err := authClient.UpdateUser(ctx, userId, params)
		if err != nil {
			apiResp.Result = ERROR
			apiResp.Message = "failed to update display name of " + userId + "."
		} else {
			err := UpdateStatusMessage(userId, statusMessage, client, ctx)
			if err != nil {
				apiResp.Result = ERROR
				apiResp.Message = "failed to update user info of " + userId + "."
			} else {
				apiResp.Result = OK
			}
		}
	}
	
	bytes, _ := json.Marshal(apiResp)
	_, _ = w.Write(bytes)
}
