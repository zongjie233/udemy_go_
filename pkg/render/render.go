package render

import (
	"fmt"
	"html/template"
	"net/http"
)

// RenderTemplate  渲染模板函数
func RenderTemplate(w http.ResponseWriter, tmpl string) {
	parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl) // 模板路径
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println("error parsing template:", err)
	}
}
