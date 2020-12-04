package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type StayingAwakeParams struct {
	UserId  string `json:"user_id"`
	IdToken string `json:"id_token"`
}
type StayingAwakeResponseStruct struct {
	Result  string       `json:"result"`
	Message string       `json:"message"`
	Users   []UserStruct `json:"users"`
}

func StayingAwake(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx, client := InitializeHttpFunc()
	defer CloseFirestoreClient(client)

	var apiResp StayingAwakeResponseStruct
	body := request.Body
	params := StayingAwakeParams{}
	_ = json.Unmarshal([]byte(body), &params)

	userId, idToken := params.UserId, params.IdToken

	if userId == "" || idToken == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else if isUserVerified, _ := IsUserVerified(userId, idToken, client, ctx); !isUserVerified {
		apiResp.Result = ERROR
		apiResp.Message = UserAuthFailed
	} else {
		_ = RecordLastAccess(userId, client, ctx)
		roomId, _ := InWhichRoom(userId, client, ctx)
		if roomId == "" {
			apiResp.Result = ERROR
			apiResp.Message = "you are not in the room."
		} else {
			users, _ := RetrieveRoomUsers(roomId, client, ctx)
			apiResp.Result = OK
			apiResp.Users = users
		}
	}

	bytes, _ := json.Marshal(apiResp)
	return Response(bytes)
}

func main() {
	lambda.Start(StayingAwake)
}
