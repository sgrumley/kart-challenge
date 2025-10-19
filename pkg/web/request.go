package web

import (
	"encoding/json"
	"io"
	"net/http"
)

func DecodeBody(r *http.Request, req interface{}) error {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}

	return json.Unmarshal(body, req)
}
