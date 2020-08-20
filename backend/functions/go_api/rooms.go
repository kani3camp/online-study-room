package go_api

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"google.golang.org/api/iterator"
	"log"
	"net/http"
)

// なぜか、構造体のキーを小文字から始めるとそのデータが返せないので大文字にするように。

type RoomsResponseStruct struct {
	Result string `json:"result"`
	Message string `json:"message"`
	Rooms []RoomStruct `json:"rooms"`
}

func retrieveRooms(client *firestore.Client, ctx context.Context) ([]RoomStruct, error) {
	var rooms []RoomStruct
	var err error
	
	// roomsのコレクションを取得
	iter := client.Collection(ROOMS).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			err = nil
			break
		}
		if err != nil {log.Fatalf("Failed to iterate: %v", err)}
		var _room RoomBodyStruct
		_ = doc.DataTo(&_room)
		room := RoomStruct {
			RoomId: doc.Ref.ID,
			Body:   _room,
		}
		rooms = append(rooms, room)
	}
	return rooms, err
}

func Rooms(w http.ResponseWriter, r *http.Request) {
	ctx, client := InitializeHttpFunc(&w)
	defer client.Close()
	
	UpdateDatabase(client, ctx)
	
	var apiResp RoomsResponseStruct
	
	rooms, err := retrieveRooms(client, ctx)
	if err != nil {
		log.Println(err)
		apiResp.Result = ERROR
		apiResp.Message = err.Error()
	} else {
		apiResp.Result = OK
		apiResp.Rooms = rooms
	}
	
	bytes, _ := json.Marshal(apiResp)
	_, _ = w.Write(bytes)
}
