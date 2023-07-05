package render

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/justinas/nosurf"
	"github.com/zongjie233/udemy_lesson/internal/config"
	"github.com/zongjie233/udemy_lesson/internal/models"
	"html/template"
	"net/http"
	"path/filepath"
)

var (
	functions = template.FuncMap{}
	// "undefined app = a" is most likely due to the variable "app" not being defined in the package-level scope.
	app *config.AppConfig
	// test setup
	pathToTemplates = "./templates"
)

// NewRenderer 为模板设定配置
func NewRenderer(a *config.AppConfig) {
	app = a

}

// AddDefaultData 添加数据函数
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}
	//someone

	return td
}

// Template  渲染模板函数
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {

		tc, _ = CreateTemplateCache()
		// var err error
		// tc, err = CreateTemplateCache()
		// if err != nil {
		// 	log.Fatal(err)
		// 	return
		// }
	}

	t, ok := tc[tmpl]
	if !ok {
		return errors.New("could not get template from cache")
	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td, r)
	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
		return err
	}
	return nil
}

// CreateTemplateCache 创建模板缓存
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := make(map[string]*template.Template)

	// 获取所有 .page.tmpl 文件
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates)) // filepath.Glob()用于返回与指定模式匹配的所有文件或目录的名称，以切片模式返回
	if err != nil {
		return myCache, err
	}

	// 遍历所有 page.tmpl 文件并创建一个模板对象
	for _, page := range pages {
		name := filepath.Base(page)                                     // 返回路径中的最后一个元素,即文件名
		ts, err := template.New(name).Funcs(functions).ParseFiles(page) // 创建一个新模板对象
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates)) // 这里是layout
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates)) // 这里是layout
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts // 将模板对象存储到模板缓存中
	}
	return myCache, nil
}

/*
	第一种方法

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
