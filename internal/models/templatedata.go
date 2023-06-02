package models

import "github.com/zongjie233/udemy_lesson/internal/forms"

// TemplateData 保存从处理程序发送至模板的数据
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Form
}
