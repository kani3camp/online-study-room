package go_api

import (
	"encoding/json"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
	"os"
)


type SendContactFormResponseStruct struct {
	Result string `json:"result"`
	Message string `json:"message"`
}


// 環境変数はコンソールの関数の編集から設定してる
func SendContactForm(w http.ResponseWriter, r *http.Request)  {
	ctx, client := InitializeHttpFunc(&w)
	defer client.Close()
	
	
	var apiResp SendContactFormResponseStruct
	userId, idToken, mailAddress := r.FormValue(user_id), r.FormValue(id_token), r.FormValue("mail_address")
	message, contactType := r.FormValue("message"), r.FormValue("contact_type")
	
	message =
		"Contact type : " + contactType + "\n\n" +
		"From : " + mailAddress + "\n\n" +
		"Message : \n" + message
	
	if userId == "" || idToken == "" || mailAddress == "" || contactType == "" || message == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else {
		if IsUserVerified(userId, idToken, ctx) {
			messageDestinationId := os.Getenv("DESTINATION_LINE_ID")
			bot, err := linebot.New(
				os.Getenv("CHANNEL_SECRET"), // todo
				os.Getenv("CHANNEL_TOKEN"),
			)
			if err != nil {
				log.Println(err)
				apiResp.Result = ERROR
			} else {
				if _, err := bot.PushMessage(messageDestinationId, linebot.NewTextMessage(message)).Do(); err != nil {
					log.Println(err)
					apiResp.Result = ERROR
				} else {
					apiResp.Result = OK
					apiResp.Message = "Successfully sent your message."
				}
			}
		} else {
			apiResp.Result = ERROR
			apiResp.Message = UserAuthFailed
		}
	}
	
	bytes, _ := json.Marshal(apiResp)
	_, _ = w.Write(bytes)
}
