package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"os"
	"time"
)

type SendContactFormParams struct {
	UserId      string `json:"user_id"`
	IdToken     string `json:"id_token"`
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

	if userId == "" || idToken == "" || mailAddress == "" || contactType == "" || message == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else if isUserVerified, _ := IsUserVerified(userId, idToken, client, ctx); !isUserVerified {
		apiResp.Result = ERROR
		apiResp.Message = UserAuthFailed
	} else {
		// Firestoreに記録
		contactBody := ContactBodyStruct{
			Contacted:   time.Now(),
			UserId:      userId,
			MailAddress: mailAddress,
			Message:     message,
			ContactType: contactType,
		}
		err := RecordContact(contactBody, client, ctx)
		if err != nil {
			log.Println()
			apiResp.Result = ERROR
			apiResp.Message = "failed to record the contact data"
		} else {
			// LINEに送信
			message =
				"Contact type : " + contactType + "\n" +
					"From : " + mailAddress + "\n" +
					"Message : \n" + message
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
					apiResp.Message = "successfully sent your message"
				}
			}
		}

	}

	bytes, _ := json.Marshal(apiResp)
	return Response(bytes)
}

func main() {
	lambda.Start(SendContactForm)
}
