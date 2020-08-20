package go_api

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
)

func HelloWorld(w http.ResponseWriter, r *http.Request)  {
	var d struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		_, _ = fmt.Fprint(w, "error... Hello World!")
		_, _ = fmt.Fprint(w, err)
		return
	}
	if d.Name == "" {
		_, _ = fmt.Fprint(w, "Hello, World!")
		return
	}
	_, _ = fmt.Fprintf(w, "Hello, %s!", html.EscapeString(d.Name))
}
