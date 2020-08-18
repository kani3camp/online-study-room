package go_api

import (
	"encoding/json"
	"net/http"
)

func OnlineUsers(w http.ResponseWriter, r *http.Request)  {
	ctx, client := InitializeHttpFunc(&w)
	defer client.Close()
	
	onlineUsers := _OnlineUsers(client, ctx)
	res, _ := json.Marshal(onlineUsers)
	_, _ = w.Write(res)
}
