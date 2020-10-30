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
	ctx, client := InitializeHttpFunc()
	defer client.Close()

	userId := request.QueryStringParameters[user_id]
	var apiResp UserStatusResponseStruct

	if userId == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else if isInUsers, _ := IsInUsers(userId, client, ctx); !isInUsers {
		log.Println("p1")
		apiResp.Result = ERROR
		apiResp.Message = InvalidUser
	} else {
		log.Println("p2")
		authClient, _ := InitializeFirebaseAuthClient(ctx)
		log.Println("p3")
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
		log.Println("p5")
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
