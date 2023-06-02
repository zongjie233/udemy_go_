package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/zongjie233/udemy_lesson/internal/config"
	"github.com/zongjie233/udemy_lesson/internal/forms"
	"github.com/zongjie233/udemy_lesson/internal/models"
	"github.com/zongjie233/udemy_lesson/internal/render"
	"log"
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

	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})

}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) { // 必须有着两个参数
	// 业务逻辑
	stringMap := make(map[string]string)
	stringMap["test"] = "hello,world"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	// 向模板发送数据
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation 渲染预定页面，展示表单
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{Form: forms.New(nil)})
}

// PostReservation 处理预定表单的post请求
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	//从请求表单中获取预订信息,保存到reservation结构体
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	//- 使用forms包解析请求表单
	form := forms.New(r.PostForm)
	//- 调用Has方法校验first_name字段是否存在
	form.Has("first_name", r)
	//如果form校验未通过,则渲染make-reservation模板并返回
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{Form: form, Data: data})
		return
	}
}

// Bigbed 渲染大床房页面，展示表单
func (m *Repository) Bigbed(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "bigbed.page.tmpl", &models.TemplateData{})
}

// Basic 渲染标准间页面，展示表单
func (m *Repository) Basic(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "basicroom.page.tmpl", &models.TemplateData{})
}

// Availability 渲染查找页面，展示表单
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability 渲染查找页面，展示表单
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	// 获取表单上的数据
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("开始日期是%s,结束日期是%s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJson 处理查询请求并发送JSON响应
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{OK: true, Message: "Availability"}
	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json") // 使客户端能够正确判断和解析服务器返回的数据格式
	w.Write(out)
}

// Contact 渲染查找页面，展示表单
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}
