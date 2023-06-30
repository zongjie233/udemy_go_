package handlers

import (
	"context"
	"fmt"
	"github.com/zongjie233/udemy_lesson/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
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
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"bigbed", "/bigbed", "GET", http.StatusOK},
	{"basic", "/basicroom", "GET", http.StatusOK},
	{"search", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	//{"post-avl", "/search-availability", "POST", []postData{
	//	{key: "start", value: "2020-01-01"},
	//	{key: "end", value: "2020-02-01"},
	//}, http.StatusOK},
	//{"post-search-avl-json", "/search-availability-json", "POST", []postData{
	//	{key: "start", value: "2020-01-01"},
	//	{key: "end", value: "2020-02-01"},
	//}, http.StatusOK},
	//{"make-reservation", "/make-reservation", "POST", []postData{
	//	{key: "first_name", value: "2020-01-01"},
	//	{key: "last_name", value: "2020-02-01"},
	//	{key: "email", value: "hs@hs.com"},
	//	{key: "phone", value: "123456"},
	//}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes) // 创建一个带有自签名证书的 HTTPS 服务器，并将路由传入。 使用httptest.NewTLSServer创建一个测试服务器；
	defer ts.Close()

	for _, e := range theTests {

		resp, err := ts.Client().Get(ts.URL + e.url) // 使用ts.URL获取测试服务器的URL；
		if err != nil {
			t.Log(err)
			t.Fatal(err)

		}
		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}

	}

}

/*
创建一个模拟的预订对象 reservation，其中包含房间 ID 和房间名称。
创建一个新的 HTTP 请求，并使用 /make-reservation 作为请求的路径。
使用 getCtx 函数获取一个上下文对象 ctx，并将其设置为请求的上下文。
使用 httptest.NewRecorder 创建一个响应记录器，用于记录处理程序返回的响应。
将预订对象存储到会话中，以便在处理程序中可以访问。
创建一个处理程序 handler，并将存储库中的 Reservation 方法设置为处理程序。
使用 handler.ServeHTTP 方法调用处理程序来处理请求，并将响应记录到响应记录器中。
检查响应的状态码是否与预期的状态码相符，如果不符，则报告测试失败。
*/
func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "豪华大床房",
		},
	}

	// 创建一个HTTP请求，方法为GET，路径为"/make-reservation"，无请求体
	req, _ := http.NewRequest("GET", "/make-reservation", nil)

	// 调用getCtx函数获取上下文
	ctx := getCtx(req)

	// 将获取到的上下文设置给请求
	req = req.WithContext(ctx)

	// 创建一个ResponseRecorder
	rr := httptest.NewRecorder()

	// 将reservation存储在会话中
	session.Put(ctx, "reservation", reservation)

	// 创建一个处理Reservation的HTTP处理函数
	handler := http.HandlerFunc(Repo.Reservation)

	// 处理HTTP请求并记录响应
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d,want %d", rr.Code, http.StatusOK)
	}

	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong reponse code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test with no-existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong reponse code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostReservation(t *testing.T) {

	// 构建reqBody
	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=zdq")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=sss")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=zdq@zdq.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123123123")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	//创建一个httptest.NewRecorder()对象，用于捕获PostReservation函数生成的HTTP响应。
	rr := httptest.NewRecorder()

	//使用http.HandlerFunc将PostReservation函数转换为http.Handler接口的处理程序函数，以便能够使用handler.ServeHTTP处理请求并调用PostReservation函数
	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for missing post body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	//使用http.HandlerFunc将PostReservation函数转换为http.Handler接口的处理程序函数，以便能够使用handler.ServeHTTP处理请求并调用PostReservation函数
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code for missing post body: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid start date
	reqBody = "start_date=invalid date"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=zdq")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=sss")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=zdq@zdq.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123123123")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	//使用http.HandlerFunc将PostReservation函数转换为http.Handler接口的处理程序函数，以便能够使用handler.ServeHTTP处理请求并调用PostReservation函数
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code for invalid start date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid end date
	reqBody = "start_date=2050-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=invalid end date")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=zdq")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=sss")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=zdq@zdq.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123123123")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	//使用http.HandlerFunc将PostReservation函数转换为http.Handler接口的处理程序函数，以便能够使用handler.ServeHTTP处理请求并调用PostReservation函数
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code for invalid end date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid room id
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=zdq")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=sss")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=zdq@zdq.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123123123")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=invalid")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	//使用http.HandlerFunc将PostReservation函数转换为http.Handler接口的处理程序函数，以便能够使用handler.ServeHTTP处理请求并调用PostReservation函数
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code for invalid id: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

//func TestRepository_AvailabilityJSON(t *testing.T) {
//	reqBody := "start=2050-01-01"
//	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-01")
//	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
//
//	req, _ := http.NewRequest("POST", "search-availability-json", strings.NewReader(reqBody))
//
//	ctx := getCtx(req)
//	req = req.WithContext(ctx)
//
//	req.Header.Set("Content-Type", "x-www-form-urlencoded")
//
//	handler := http.HandlerFunc(Repo.AvailabilityJSON)
//
//	rr := httptest.NewRecorder()
//
//	handler.ServeHTTP(rr, req)
//
//	var j jsonResponse
//	err := json.Unmarshal([]byte(rr.Body.String()), &j)
//	if err != nil {
//		t.Errorf("failed to parse json")
//	}
//
//	// test case where reservation is not in session
//	req, _ = http.NewRequest("GET", "/make-reservation", nil)
//	ctx = getCtx(req)
//	req = req.WithContext(ctx)
//	rr = httptest.NewRecorder()
//
//	handler.ServeHTTP(rr, req)
//	if rr.Code != http.StatusTemporaryRedirect {
//		t.Errorf("Reservation handler returned wrong reponse code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
//	}
//
//}

// getCtx 接受一个http.Request作为参数，返回一个context.Context
func getCtx(req *http.Request) context.Context {

	// 调用session.Load方法加载会话数据
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}

	return ctx
}
