package forms

import (
	"net/http"
	"net/url"
)

// Form 创建一个form结构体
type Form struct {
	url.Values
	Errors errors
}

// Valid 没有错误返回true，否则false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New 初始化form
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has 检查表单字段名在post请求中并且不为空
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field cannot be blank")
		return false
	}
	return true
}
