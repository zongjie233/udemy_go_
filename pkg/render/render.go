package render

import (
	"bytes"
	"github.com/zongjie233/udemy_lesson/models"
	"github.com/zongjie233/udemy_lesson/pkg/config"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates 为模板设定配置
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData 添加数据函数
func AddDefaultData(td *models.TemplateData) *models.TemplateData {

	return td
}

// RenderTemplate  渲染模板函数
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	// 从缓存中取得模板
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("不能从缓存中拿到模板")
	}

	td = AddDefaultData(td)
	buf := new(bytes.Buffer) // bytes.buffer实现字节缓冲区
	_ = t.Execute(buf, td)   // 将数据写入

	// 渲染模板
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

}

// CreateTemplateCache 创建模板缓存
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := make(map[string]*template.Template)

	// 获取templates中所有*.page.tmpl文件
	pages, err := filepath.Glob("./templates/*.page.tmpl") // filepath.Glob()用于返回与指定模式匹配的所有文件或目录的名称，以切片模式返回
	if err != nil {
		return myCache, err
	}

	// 遍历所有page.tmpl文件
	for _, page := range pages {
		name := filepath.Base(page)                    // 返回路径中的最后一个元素,即文件名
		ts, err := template.New(name).ParseFiles(page) // 创建一个模板对象
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}

/* 第一种方法
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

*/
