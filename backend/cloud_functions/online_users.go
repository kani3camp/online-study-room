package go_api

import (
	"encoding/json"
	"net/http"
)

type OnlineUsersResponse struct {
	Result      string       `json:"result"`
	Message     string       `json:"message"`
	OnlineUsers []UserStruct `json:"online_users"`
}

func OnlineUsers(w http.ResponseWriter, r *http.Request) {
	ctx, client := InitializeHttpFunc(&w)
	defer client.Close()
	
	var apiResp OnlineUsersResponse
	onlineUsers, _ := RetrieveOnlineUsers(client, ctx)
	apiResp.Result = OK
	apiResp.OnlineUsers = onlineUsers
	
	bytes, _ := json.Marshal(apiResp)
	_, _ = w.Write(bytes)
}
