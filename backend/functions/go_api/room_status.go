package go_api

import (
	"encoding/json"
	"log"
	"net/http"
)

type RoomStatusResponseStruct struct {
	Result string `json:"result"`
	Message string `json:"message"`
	RoomStatus RoomStruct `json:"room_status"`
}

func RoomStatus(w http.ResponseWriter, r *http.Request) {
	ctx, client := InitializeHttpFunc(&w)
	defer client.Close()
	
	var apiResp RoomStatusResponseStruct
	roomId := r.FormValue(room_id)
	
	if roomId == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else {
		if IsExistRoom(roomId, client, ctx) {
			room, err := client.Collection(ROOMS).Doc(roomId).Get(ctx)
			if err != nil {
				log.Fatalln(err)
			} else {
				apiResp.Result = OK
				var roomBody RoomBodyStruct
				_ = room.DataTo(&roomBody)
				apiResp.RoomStatus = RoomStruct{
					RoomId: roomId,
					Body:   roomBody,
				}
			}
		} else {
			log.Println(RoomDoesNotExist)
			apiResp.Result = ERROR
			apiResp.Message = RoomDoesNotExist
		}
	}
	bytes, _ := json.Marshal(apiResp)
	_, _ = w.Write(bytes)
}
