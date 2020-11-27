package go_api

import (
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
	"os"
)


// 環境変数はコンソールの関数の編集から設定してる
func LineWebhookCallBack(w http.ResponseWriter, r *http.Request)  {
	client, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
		)
	if err != nil {
		http.Error(w, "Error init client", http.StatusBadRequest)
		log.Fatal(err)
		return
	}
	events, err := client.ParseRequest(r)
	if err != nil {
		http.Error(w, "Error parse request", http.StatusBadRequest)
		log.Fatal(err)
		return
	}
	for _, e := range events {
		switch e.Type {
		case linebot.EventTypeMessage:
			message := linebot.NewTextMessage("Test")
			_, err := client.ReplyMessage(e.ReplyToken, message).Do()
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
	_, _ = fmt.Fprint(w, "ok")
}
