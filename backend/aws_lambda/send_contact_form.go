package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"os"
)

type SendContactFormParams struct {
	UserId      string `json:"UserId"`
	IdToken     string `json:"IdToken"`
	MailAddress string `json:"mail_address"`
	Message     string `json:"message"`
	ContactType string `json:"contact_type"`
}
type SendContactFormResponseStruct struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

// 環境変数はコンソールの関数の編集から設定してる
func SendContactForm(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("SendContactForm()")
	ctx, client := InitializeHttpFuncWithFirestore()
	defer CloseFirestoreClient(client)

	var apiResp SendContactFormResponseStruct
	body := request.Body
	params := SendContactFormParams{}
	_ = json.Unmarshal([]byte(body), &params)

	userId, idToken, mailAddress := params.UserId, params.IdToken, params.MailAddress
	message, contactType := params.Message, params.ContactType

	message =
		"Contact type : " + contactType + "\n\n" +
			"From : " + mailAddress + "\n\n" +
			"Message : \n" + message

	if userId == "" || idToken == "" || mailAddress == "" || contactType == "" || message == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else if isUserVerified, _ := IsUserVerified(userId, idToken, client, ctx); !isUserVerified {
		apiResp.Result = ERROR
		apiResp.Message = UserAuthFailed
	} else {
		messageDestinationId := os.Getenv("DESTINATION_LINE_ID")
		bot, err := linebot.New(
			os.Getenv("CHANNEL_SECRET"),
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
				apiResp.Message = "successfully sent your message."
			}
		}
	}

	bytes, _ := json.Marshal(apiResp)
	return Response(bytes)
}

func main() {
	lambda.Start(SendContactForm)
}
