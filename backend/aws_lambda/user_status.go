package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

type UserStatusResponseStruct struct {
	Result     string     `json:"result"`
	Message    string     `json:"message"`
	UserStatus UserStruct `json:"user_status"`
}

func UserStatus(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("UserStatus()")
	ctx, client := InitializeHttpFuncWithFirestore()
	defer CloseFirestoreClient(client)

	userId := request.QueryStringParameters[UserId]
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
			log.Println("failed authClient.GetUser(ctx, userId).")
			displayName = ""
		} else {
			displayName = user.DisplayName
		}

		userInfo, _ := RetrieveUserInfo(userId, client, ctx)
		apiResp.UserStatus = UserStruct{
			UserId:      userId,
			DisplayName: displayName,
			Body:        userInfo,
		}
		apiResp.Result = OK
	}

	bytes, _ := json.Marshal(apiResp)
	return Response(bytes)
}

func main() {
	lambda.Start(UserStatus)
}
