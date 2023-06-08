package render

import (
	"github.com/zongjie233/udemy_lesson/internal/models"
	"net/http"
	"testing"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession() // 调用 getSession 函数，获取一个带有 session 的 *http.Request
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123") // 往 session 中添加一个键值对，键为 "flash"，值为 "123"
	result := AddDefaultData(&td, r)         // 调用 AddDefaultData 函数，传入 td 和 r，并接收返回的 TemplateData

	if result.Flash != "123" { // 检查返回的 TemplateData 是否包含正确的 flash 值
		t.Error("flash value of 123 not found in session")
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	r, err := getSession() // 创建request
	if err != nil {
		t.Error(err)
	}
	var ww myWriter

	err = RenderTemplate(&ww, r, "home.page.tmpl", &models.TemplateData{})

	if err != nil {
		t.Error("error writing template to browser")
	}

	err = RenderTemplate(&ww, r, "non-existent.page.tmpl", &models.TemplateData{})

	if err == nil {
		t.Error("rendered template that does not exist")
	}
}

// getSession 是一个辅助函数，用于创建一个带有 session 的 *http.Request
func getSession() (*http.Request, error) {

	r, err := http.NewRequest("GET", "some-url", nil) // 创建一个 GET 请求，URL 为 "some-url"
	if err != nil {
		return nil, err
	}

	ctx := r.Context()                                    // 从 r 中获取 Context
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session")) // 从请求头中获取 "X-Session" 字段的值，并将其加载到 ctx
	r = r.WithContext(ctx)                                // 使用带有新 ctx 的 r 替换原来的 r

	return r, nil
}

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}

// 这个类型和方法的主要目的是实现 http.ResponseWriter 接口，在自定义的 HTTP 处理程序中使用自定义的响应方式。
type myWriter struct {
}

func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *myWriter) WriteHeader(i int) {

}

func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}

/*
单元测试：这段代码展示了如何编写一个Go语言的单元测试用例。测试用例函数以Test开头，并接收一个*testing.T类型的参数。在测试用例中，可使用t.Error或t.Errorf等方法记录错误信息。

net/http包：这段代码使用了http.NewRequest函数创建一个HTTP请求。http.Request结构体包含了一个HTTP请求的所有信息，如请求方法、URL、头部等。

context：http.Request结构体包含一个context.Context类型的字段。Context可以在函数之间传递数据和取消信号，是Go语言中实现跨函数数据传输和协作的一种方式。

session：这段代码使用了github.com/alexedwards/scs/v2包实现session管理。session.Load函数用于加载session数据到context.Context，session.Put函数用于往session中添加键值对。

错误处理：在Go语言中，函数通常会返回一个错误值（error类型）。我们需要检查这个错误值，判断是否发生了错误。同时，函数也可以返回一个带有错误信息的error类型值。
*/
