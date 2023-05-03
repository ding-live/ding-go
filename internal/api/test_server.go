package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

const (
	testApiKey        = "valid_api_key"
	testInvalidApiKey = "invalid_api_key"
)

func emulateResponse(w http.ResponseWriter, r *http.Request, rawResponse string) {
	if r.Header.Get(APIKeyHeader) != testApiKey {
		msg := GatewayErrorMessage{
			Message: "Forbidden",
		}

		b, _ := json.Marshal(msg)

		w.WriteHeader(http.StatusForbidden)
		w.Write(b)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write([]byte(rawResponse))
}

func testServer(res string, status ...int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(status) > 0 {
			w.WriteHeader(status[0])
		}

		emulateResponse(w, r, res)
	}))
}
