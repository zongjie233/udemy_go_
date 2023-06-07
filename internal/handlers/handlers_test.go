package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes) // 创建一个带有自签名证书的 HTTPS 服务器，并将路由传入。 使用httptest.NewTLSServer创建一个测试服务器；
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url) // 使用ts.URL获取测试服务器的URL；
			if err != nil {
				t.Log(err)
				t.Fatal(err)

			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}

		} else {

		}
	}
}
