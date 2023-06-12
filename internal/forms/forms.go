package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/url"
	"strings"
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

// Required 检查必填字段合法性
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Has 检查表单字段名在post请求中并且不为空
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	if x == "" {

		return false
	}
	return true
}

// MinLength 检查长度是否符合要求
func (f *Form) MinLength(field string, length int) bool {
	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}
