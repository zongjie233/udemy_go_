package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"bigbed", "/bigbed", "GET", []postData{}, http.StatusOK},
	{"basic", "/basicroom", "GET", []postData{}, http.StatusOK},
	{"search", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"make", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"post-avl", "/search-availability", "POST", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-02-01"},
	}, http.StatusOK},
	{"post-search-avl-json", "/search-availability-json", "POST", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-02-01"},
	}, http.StatusOK},
	{"make-reservation", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "2020-01-01"},
		{key: "last_name", value: "2020-02-01"},
		{key: "email", value: "hs@hs.com"},
		{key: "phone", value: "123456"},
	}, http.StatusOK},
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
			values := url.Values{} //url.Values 被用于创建 POST 请求的表单数据
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+e.url, values) //PostForm() 函数会自动将表单数据进行编码，并将其作为消息体发送到目标 URL 地址
			if err != nil {
				t.Log(err)
				t.Fatal(err)

			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
