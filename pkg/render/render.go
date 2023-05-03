package render

import (
	"fmt"
	"html/template"
	"net/http"
)

// RenderTemplate  渲染模板函数
func RenderTemplate(w http.ResponseWriter, tmpl string) {
	parsedTemplate, _ := template.ParseFiles("./templates/"+tmpl, "./templates/base_layout.tmpl") // 如果使用layout，则后边加上layout的文件位置
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println("error parsing template:", err)
	}
}
