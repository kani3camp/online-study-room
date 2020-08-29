package go_api

import (
	"cloud.google.com/go/firestore"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type CreateNewRoomResponseStruct struct {
	Result string `json:"result"`
	Message string `json:"message"`
}

func CreateNewRoom(w http.ResponseWriter, r *http.Request)  {
	ctx, client := InitializeHttpFunc(&w)
	defer client.Close()
	
	var apiResp ChangeUserInfoResponseStruct
	roomId, roomName, roomType := r.FormValue(room_id), r.FormValue("room_name"), r.FormValue("room_type")
	password := r.FormValue("password")
	
	if roomId == "" || roomName == "" || roomType == "" || password == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else {
		if password == os.Getenv("password") {
			continueFlag := true
			// 同じ部屋がすでにないかチェック
			rooms, _ := retrieveRooms(client, ctx)
			for _, room := range rooms {
				if room.RoomId == roomId {
					apiResp.Result = ERROR
					apiResp.Message = "The room id is already used."
					continueFlag = false
					break
				}
			}
			if continueFlag {
				_, err := client.Collection(ROOMS).Doc(roomId).Set(ctx, map[string]interface{}{
					"name":    roomName,
					"type":    roomType,
					"users":   []string{},
					"created": firestore.ServerTimestamp,
				}, firestore.MergeAll)
				if err != nil {
					log.Println(err)
					apiResp.Result = ERROR
					apiResp.Message = "Failed to create room."
				} else {
					apiResp.Result = OK
					apiResp.Message = "Successfully created room named " + roomId
				}
			}
		} else {
			apiResp.Result = ERROR
			apiResp.Message = "Invalid password."
		}
	}
	bytes, _ := json.Marshal(apiResp)
	_, _ = w.Write(bytes)
}
