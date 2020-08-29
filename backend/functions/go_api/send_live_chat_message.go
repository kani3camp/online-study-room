package go_api

import (
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AccessTokenResponseStruct struct {
	AccessToken string `json:"access_token"`
	ExpiresIn int `json:"expires_in"`
	Scope string `json:"scope"`
	TokenType string `json:"token_type"`
}

type ApiConfigStruct struct {
	AccessToken string `firestore:"access-token"`
	LiveChatId string `firestore:"live-chat-id"`
	LiveChatMessageUrl string `firestore:"live-chat-message-url"`
	RefreshToken string `firestore:"refresh-token"`
	ExpireDate time.Time `firestore:"expire-date"`
	ClientId string `firestore:"client-id"`
	ClientSecret string `firestore:"client-secret"`
	OAuthRefreshTokenUrl string `firestore:"o-auth-refresh-token-url"`
}


func RefreshToken(config *ApiConfigStruct, client *firestore.Client, ctx context.Context) error {
	data := url.Values{}
	data.Set("client_id", config.ClientId)
	data.Add("client_secret", config.ClientSecret)
	data.Add("refresh_token", config.RefreshToken)
	data.Add("grant_type", "refresh_token")
	
	req, err := http.NewRequest(
		http.MethodPost,
		config.OAuthRefreshTokenUrl,
		strings.NewReader(data.Encode()),
		)
	if err != nil {log.Println(err)}
	if req != nil {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {log.Println(err)}
	if resp != nil {
		defer resp.Body.Close()
		
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {log.Println(err)}
		
		var responseBody AccessTokenResponseStruct
		err = json.Unmarshal(body, &responseBody)
		if err != nil {log.Println(err)}
		
		log.Println(string(body))
		
		config.AccessToken = responseBody.AccessToken
		config.ExpireDate = time.Now().Add(time.Duration(responseBody.ExpiresIn))
		_, err = client.Collection(CONFIG).Doc(API).Set(ctx, map[string]interface{}{
			"access-token" : config.AccessToken,
			"expire-date" : config.ExpireDate,
		}, firestore.MergeAll)
		if err != nil {log.Fatalln(err)}
		
		return err
	} else {
		return err
	}
}

type ChatMessageRequestStruct struct {
	Snippet SnippetStruct `json:"snippet"`
}
type SnippetStruct struct {
	LiveChatId string `json:"liveChatId"`
	Type string `json:"type"`
	TextMessageDetails TextMessageDetailsStruct `json:"textMessageDetails"`
}
type TextMessageDetailsStruct struct {
	MessageText string `json:"messageText"`
}


func PostMessage(message string, config *ApiConfigStruct) error {
	
	requestBody := ChatMessageRequestStruct{
		Snippet: SnippetStruct{
			LiveChatId: config.LiveChatId,
			Type: "textMessageEvent",
			TextMessageDetails: TextMessageDetailsStruct{
				MessageText: message,
			},
		},
	}
	jsonStr, _ := json.Marshal(requestBody)
	
	req, err := http.NewRequest(
		http.MethodPost,
		config.LiveChatMessageUrl,
		bytes.NewBuffer(jsonStr))
	if err != nil {log.Println(err)}
	if req != nil {
		req.Header.Add("Authorization", "Bearer " + config.AccessToken)
		
		values := url.Values{} // url.Valuesオブジェクト生成
		values.Add("part", "snippet") // key-valueを追加
		req.URL.RawQuery = values.Encode()
		
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {log.Println(err)}
		if resp != nil {
			defer resp.Body.Close()
			
			//body, _ := ioutil.ReadAll(resp.Body)
			//var b interface{}
			//_ = json.Unmarshal(body, &b)
			//log.Println(b)
		}
	}
	return err
}


func SendLiveChatMessage(message string, client *firestore.Client, ctx context.Context) {
	var config ApiConfigStruct
	configDoc, err := client.Collection(CONFIG).Doc(API).Get(ctx)
	if err != nil {log.Fatalln(err)}
	_ = configDoc.DataTo(&config)
	
	if config.LiveChatId == "" {	// live-chat-idが空欄だったら配信はないと判断して終了
		return
	}
	
	if config.ExpireDate.Before(time.Now()) {
		log.Println("Access token is expired. Refreshing...")
		_ = RefreshToken(&config, client, ctx)
		//log.Println("New access_token is : " + config.AccessToken)
	}
	
	err = PostMessage(message, &config)
	if err != nil {log.Fatalln(err)}
}
