package handlers

import (
	"github.com/zongjie233/udemy_lesson/models"
	"github.com/zongjie233/udemy_lesson/pkg/config"
	"github.com/zongjie233/udemy_lesson/pkg/render"
	"net/http"
)

// Repo 处理程序使用的存储库
var Repo *Repository

// Repository 是一个库的类型
type Repository struct {
	App *config.AppConfig
}

// NewRepo 创建一个新库
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers 为处理程序设置存储库
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the Home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) { // 必须有着两个参数
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP) // 将访问ip存入session中，key值为”remote_ip“

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})

}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) { // 必须有着两个参数
	// 业务逻辑
	stringMap := make(map[string]string)
	stringMap["test"] = "hello,world"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	// 向模板发送数据
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
