package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/what", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/what", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("something wrong")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	// 创建第二个request请求，便于测试Required函数
	r, _ = http.NewRequest("POST", "/what", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("show does not have required fields when it does")
	}

}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/what", nil)
	form := New(r.PostForm)

	if form.Has("a") {
		t.Error("出错咯")
	}

	postedData := url.Values{}
	postedData.Set("a", "a")
	form = New(postedData)

	if !form.Has("a") {
		t.Error("出错啦")
	}

}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/what", nil)
	form := New(r.PostForm)

	form.MinLength("a", 10)
	if form.Valid() {
		t.Error("minlength在没有字段的情况下调用了")
	}

	isError := form.Errors.Get("a")
	if isError == "" {
		t.Error("应该有错误，但是没有报错")
	}

	postedValues := url.Values{}
	postedValues.Add("some", "some")

	form = New(postedValues)

	form.MinLength("some", 100)

	if form.Valid() {
		t.Error("minlength函数没有起作用")
	}

	postedValues = url.Values{}
	postedValues.Add("another", "some")
	form = New(postedValues)

	form.MinLength("another", 1)
	isError = form.Errors.Get("another")
	if isError != "" {
		t.Error("不应该有错误，但是报错了")
	}

}

func TestForm_IsEmail(t *testing.T) {
	postedValue := url.Values{}

	form := New(postedValue)

	form.IsEmail("x") // 会添加错误

	if form.Valid() {
		t.Error("form shows valid email for non field")
	}

	postedValue.Add("email", "hs@zdq.com")
	form = New(postedValue)

	form.IsEmail("email")

	if !form.Valid() {
		t.Error("我们得到了一个无效的邮箱,而此时不应该得到无效邮箱")
	}

	postedValue = url.Values{}
	postedValue.Add("email", "s")
	form = New(postedValue)

	form.IsEmail("email")

	if form.Valid() {
		t.Error("这个错误信息表示:我们得到了一个针对无效邮箱地址的有效回应,而此时不应该得到这样的回应。")
	}

}
