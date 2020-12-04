package main

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

type LiveBroadCastsResponseStruct struct {
	Kind     string `json:"kind"`
	Etag     string `json:"etag"`
	PageInfo struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []LiveBroadCastItemStruct `json:"items"`
}

type LiveBroadCastItemStruct struct {
	Kind    string                     `json:"kind"`
	Etag    string                     `json:"etag"`
	Id      string                     `json:"id"`
	Snippet LiveBroadCastSnippetStruct `json:"snippet"`
}

type LiveBroadCastSnippetStruct struct {
	PublishedAt        time.Time   `json:"publishedAt"`
	ChannelId          string      `json:"channelId"`
	Title              string      `json:"title"`
	Description        string      `json:"description"`
	Thumbnails         interface{} `json:"thumbnails"`
	ScheduledStartTime time.Time   `json:"scheduledStartTime"`
	ActualStartTime    time.Time   `json:"actualStartTime"`
	IsDefaultBroadcast bool        `json:"isDefaultBroadcast"`
	LiveChatId         string      `json:"liveChatId"`
}

type AccessTokenResponseStruct struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type ApiConfigStruct struct {
	AccessToken          string    `firestore:"access-token"`             // 自動
	LiveChatId           string    `firestore:"live-chat-id"`             // 自動
	LiveChatMessageUrl   string    `firestore:"live-chat-message-url"`    // 要設定
	RefreshToken         string    `firestore:"refresh-token"`            // 要設定
	ExpireDate           time.Time `firestore:"expire-date"`              // 自動
	ClientId             string    `firestore:"client-id"`                // 要設定
	ClientSecret         string    `firestore:"client-secret"`            // 要設定
	OAuthRefreshTokenUrl string    `firestore:"o-auth-refresh-token-url"` // 要設定
}

type ChatMessageRequestStruct struct {
	Snippet SnippetStruct `json:"snippet"`
}
type SnippetStruct struct {
	LiveChatId         string                   `json:"liveChatId"`
	Type               string                   `json:"type"`
	TextMessageDetails TextMessageDetailsStruct `json:"textMessageDetails"`
}
type TextMessageDetailsStruct struct {
	MessageText string `json:"messageText"`
}

func RefreshLiveChatId(config *ApiConfigStruct, client *firestore.Client, ctx context.Context) error {
	const LiveBroadCastsUrl = "https://www.googleapis.com/youtube/v3/liveBroadcasts"
	req, err := http.NewRequest(
		http.MethodGet,
		LiveBroadCastsUrl,
		nil,
	)
	if err != nil {
		log.Println(err)
	}
	if req != nil {
		req.Header.Add("Authorization", "Bearer "+config.AccessToken)
		params := req.URL.Query()
		params.Add("part", "snippet")
		params.Add("broadcastStatus", "active")
		req.URL.RawQuery = params.Encode()
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	if resp != nil {
		defer func() {
			_ = resp.Body.Close()
		}()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}

		var responseBody LiveBroadCastsResponseStruct
		err = json.Unmarshal(body, &responseBody)
		if err != nil {
			log.Println(err)
		}

		log.Println(responseBody)
		if len(responseBody.Items) > 0 {
			config.LiveChatId = responseBody.Items[0].Snippet.LiveChatId
			_, err = client.Collection(CONFIG).Doc(API).Set(ctx, map[string]interface{}{
				"live-chat-id": config.LiveChatId,
			}, firestore.MergeAll)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			log.Fatalln("No live streaming now.")
		}
		return err
	} else {
		return err
	}
}

func RefreshAccessToken(config *ApiConfigStruct, client *firestore.Client, ctx context.Context) error {
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
	if err != nil {
		log.Println(err)
	}
	if req != nil {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	if resp != nil {
		defer func() {
			_ = resp.Body.Close()
		}()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}

		var responseBody AccessTokenResponseStruct
		err = json.Unmarshal(body, &responseBody)
		if err != nil {
			log.Println(err)
		}

		//log.Println(string(body))

		config.AccessToken = responseBody.AccessToken
		config.ExpireDate = time.Now().Add(time.Duration(responseBody.ExpiresIn) * time.Second)
		_, err = client.Collection(CONFIG).Doc(API).Set(ctx, map[string]interface{}{
			"access-token": config.AccessToken,
			"expire-date":  config.ExpireDate,
		}, firestore.MergeAll)
		if err != nil {
			log.Fatalln(err)
		}

		return err
	} else {
		return err
	}
}

func PostMessage(message string, config *ApiConfigStruct) (int, error) {
	// POSTなので
	requestBody := ChatMessageRequestStruct{
		Snippet: SnippetStruct{
			LiveChatId: config.LiveChatId,
			Type:       "textMessageEvent",
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
	if err != nil {
		log.Println(err)
	}
	if req != nil {
		req.Header.Add("Authorization", "Bearer "+config.AccessToken)

		values := url.Values{}        // url.Valuesオブジェクト生成
		values.Add("part", "snippet") // key-valueを追加
		req.URL.RawQuery = values.Encode()

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
		}
		if resp != nil {
			defer func() {
				_ = resp.Body.Close()
			}()

			//body, _ := ioutil.ReadAll(resp.Body)
			//var b interface{}
			//_ = json.Unmarshal(body, &b)
			//log.Println(b)

			return resp.StatusCode, err
		}
	}
	return 0, err
}

func SendLiveChatMessage(message string, client *firestore.Client, ctx context.Context) {
	var config ApiConfigStruct
	configDoc, err := client.Collection(CONFIG).Doc(API).Get(ctx)
	if err != nil {
		log.Println(err)
	}
	_ = configDoc.DataTo(&config)

	if config.ExpireDate.Before(time.Now()) {
		log.Println("access token is expired. refreshing...")
		_ = RefreshAccessToken(&config, client, ctx)
	}

	statusCode, err := PostMessage(message, &config)
	if err != nil {
		log.Println(err)
	} else if statusCode != 200 {
		log.Println("first post has failed. response status code: ", statusCode)
		err1 := RefreshLiveChatId(&config, client, ctx)
		if err1 != nil {
			log.Println(err1)
		}
		statusCode, err2 := PostMessage(message, &config)
		if err2 != nil {
			log.Println(err2)
		} else if statusCode != 200 {
			log.Println("second post has also failed. response status code: ", statusCode)
		}
	}
}
