package utils_test

import (
	"jokes-bapak2-api/utils"
	"net/http"
	"testing"
)

func TestRequest_Get(t *testing.T) {
	res, err := utils.Request(utils.RequestConfig{
		URL:    "https://jsonplaceholder.typicode.com/todos/1",
		Method: http.MethodGet,
		Headers: map[string]interface{}{
			"User-Agent": "Jokesbapak2 Test API",
			"Accept":     "application/json",
		},
	})
	if err != nil {
		t.Error(err.Error())
	}
	if res.StatusCode != 200 {
		t.Error("response does not have 200 status", res.Status)
	}
}
