package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"google.golang.org/api/iterator"
	"log"
)


type RoomsResponseStruct struct {
	Result  string       `json:"result"`
	Message string       `json:"message"`
	Rooms   []RoomStruct `json:"rooms"`
}

func RetrieveRooms(client *firestore.Client, ctx context.Context) ([]RoomStruct, error) {
	var rooms []RoomStruct

	// roomsのコレクションを取得
	iter := client.Collection(ROOMS).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			err = nil
			break
		}
		if err != nil {
			log.Printf("Failed to iterate: %v", err)
			return []RoomStruct{}, err
		}
		var _room RoomBodyStruct
		_ = doc.DataTo(&_room)
		room := RoomStruct{
			RoomId: doc.Ref.ID,
			Body:   _room,
		}
		if room.Body.Users == nil {
			room.Body.Users = []string{}
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx, client := InitializeHttpFunc()
	defer client.Close()

	UpdateDatabase(client, ctx)

	var apiResp RoomsResponseStruct

	rooms, _ := RetrieveRooms(client, ctx)
	apiResp.Result = OK
	apiResp.Rooms = rooms

	jsonBytes, _ := json.Marshal(apiResp)
	return events.APIGatewayProxyResponse{
		Body: string(jsonBytes),
		StatusCode: 200,	// これないとInternal Server Errorになる
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
