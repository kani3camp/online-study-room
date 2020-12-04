package go_api

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// なぜか、構造体のキーを小文字から始めるとそのデータが返せないので大文字にするように。

type NewsResponseStruct struct {
	Result   string       `json:"result"`
	Message  string       `json:"message"`
	NewsList []NewsStruct `json:"news_list"`
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
