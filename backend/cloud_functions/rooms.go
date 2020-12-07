package go_api

import (
	"encoding/json"
	"net/http"
)

// なぜか、構造体のキーを小文字から始めるとそのデータが返せないので大文字にするように。

type RoomsResponseStruct struct {
	Result  string       `json:"result"`
	Message string       `json:"message"`
	Rooms   []RoomStruct `json:"rooms"`
}

func Rooms(w http.ResponseWriter, r *http.Request) {
	ctx, client := InitializeHttpFuncWithFirestore()
	defer CloseFirestoreClient(client)

	var apiResp RoomsResponseStruct
	
	rooms, _ := RetrieveRooms(client, ctx)
	apiResp.Result = OK
	apiResp.Rooms = rooms
	
	bytes, _ := json.Marshal(apiResp)
	_, _ = w.Write(bytes)
}
