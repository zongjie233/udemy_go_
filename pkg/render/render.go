package render

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// RenderTemplate  渲染模板函数
func RenderTemplateTEst(w http.ResponseWriter, tmpl string) {
	//parsedTemplate, _ := template.ParseFiles("D:\\udemy_lessons\\templates"+"\\"+tmpl, "D:\\udemy_lessons\\templates\\base.layout.tmpl") // 如果使用layout，则后边加上layout的文件位置
	//parsedTemplate, e := template.ParseFiles("D:\\udemy_lessons\\templates" + "\\" + tmpl) // 使用绝对路径能够保证不出错
	parsedTemplate, _ := template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.tmpl") // 相对路径弊端：相对于当前启动程序的目录。
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println("error parsing template:", err)
		return
	}
}

var tc = make(map[string]*template.Template)

// 缓存机制
func RenderTemplate(w http.ResponseWriter, t string) {
	var tmpl *template.Template
	var err error

	// 检查是否有模板已经在缓存里
	_, inMap := tc[t]
	if !inMap {
		// 需要创建template
		log.Println("创建模板并添加到缓存")
		err = createTemplateCache(t)
		if err != nil {
			log.Println(err)
		}
	} else {
		// 缓存中已经存在
		fmt.Println("缓存中已经有了")
	}
	tmpl = tc[t]

	err = tmpl.Execute(w, nil) // 生成HTML页面
	if err != nil {
		log.Println(err)
	}
}

// 创建模板缓存
func createTemplateCache(t string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", t),
		"./templates/base.layout.tmpl",
	}

	// 解析模板
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

	// 向缓存中加入模板
	tc[t] = tmpl

	return nil
}
