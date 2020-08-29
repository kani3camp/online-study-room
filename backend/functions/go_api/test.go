package go_api

import (
	"encoding/json"
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request)  {
	ctx, client := InitializeHttpFunc(&w)
	defer client.Close()

	userId := r.FormValue("user_id")
	idToken := r.FormValue("id_token")
	token := IsUserVerified(userId, idToken, ctx)
	
	bytes, _ := json.Marshal(token)
	_, _ = w.Write(bytes)
}