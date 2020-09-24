package go_api

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"google.golang.org/api/iterator"
	"log"
	"net/http"
	"strconv"
)

// なぜか、構造体のキーを小文字から始めるとそのデータが返せないので大文字にするように。

type NewsResponseStruct struct {
	Result   string       `json:"result"`
	Message  string       `json:"message"`
	NewsList []NewsStruct `json:"news_list"`
}

func RetrieveNews(numNews int, client *firestore.Client, ctx context.Context) ([]NewsStruct, error) {
	var newsList []NewsStruct
	var err error
	
	// roomsのコレクションを取得
	iter := client.Collection(NEWS).OrderBy("updated", firestore.Desc).Limit(numNews).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			err = nil
			break
		}
		if err != nil {
			log.Printf("failed to iterate: %v\n", err)
			return []NewsStruct{}, err
		}
		var _news NewsBodyStruct
		_ = doc.DataTo(&_news)
		news := NewsStruct{
			NewsId:   doc.Ref.ID,
			NewsBody: _news,
		}
		newsList = append(newsList, news)
	}
	return newsList, err
}

func News(w http.ResponseWriter, r *http.Request) {
	ctx, client := InitializeHttpFunc(&w)
	defer client.Close()
	
	var apiResp NewsResponseStruct
	
	_numNews := r.FormValue("num_news")
	if _numNews == "" {
		apiResp.Result = ERROR
		apiResp.Message = InvalidParams
	} else {
		numNews, err := strconv.Atoi(_numNews)
		if err != nil {
			apiResp.Result = ERROR
			apiResp.Message = InvalidValue
		} else if numNews < 1 {
			apiResp.Result = ERROR
			apiResp.Message = InvalidValue
		} else {
			newsList, _ := RetrieveNews(numNews, client, ctx)
			apiResp.Result = OK
			apiResp.NewsList = newsList
		}
	}
	
	bytes, _ := json.Marshal(apiResp)
	_, _ = w.Write(bytes)
}
