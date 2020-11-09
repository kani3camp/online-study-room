package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)


type RoomsResponseStruct struct {
	Result  string       `json:"result"`
	Message string       `json:"message"`
	Rooms   []RoomStruct `json:"rooms"`
}


func Rooms(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx, client := InitializeHttpFunc()
	defer client.Close()

	UpdateDatabase(client, ctx)

	var apiResp RoomsResponseStruct

	rooms, _ := RetrieveRooms(client, ctx)
	apiResp.Result = OK
	apiResp.Rooms = rooms

	jsonBytes, _ := json.Marshal(apiResp)
	return Response(jsonBytes)
}

func main() {
	lambda.Start(Rooms)
}
