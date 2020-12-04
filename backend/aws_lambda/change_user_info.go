package main

import (
	"encoding/json"
	"firebase.google.com/go/auth"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type ChangeUserInfoParams struct {
	UserId        string `json:"user_id"`
	IdToken       string `json:"id_token"`
	DisplayName   string `json:"display_name"`
	StatusMessage string `json:"status_message"`
}

type ChangeUserInfoResponseStruct struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

func ChangeUserInfo(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx, client := InitializeHttpFunc()
	defer CloseFirestoreClient(client)

	var apiResp ChangeUserInfoResponseStruct
	body := request.Body
	params := ChangeUserInfoParams{}
	_ = json.Unmarshal([]byte(body), &params)

	userId, idToken := params.UserId, params.IdToken
	displayName, statusMessage := params.DisplayName, params.StatusMessage

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
	return Response(bytes)
}

func main() {
	lambda.Start(ChangeUserInfo)
}
