package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/zongjie233/udemy_lesson/internal/config"
	"github.com/zongjie233/udemy_lesson/internal/driver"
	"github.com/zongjie233/udemy_lesson/internal/forms"
	"github.com/zongjie233/udemy_lesson/internal/helpers"
	"github.com/zongjie233/udemy_lesson/internal/models"
	"github.com/zongjie233/udemy_lesson/internal/render"
	"github.com/zongjie233/udemy_lesson/internal/repository"
	"github.com/zongjie233/udemy_lesson/internal/repository/dbrepo"
	"log"
	"net/http"
)

// Repo 处理程序使用的存储库
var Repo *Repository

// Repository 是一个库的类型
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo 创建一个新库
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers 为处理程序设置存储库
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the Home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) { // 必须有着两个参数
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})

}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) { // 必须有着两个参数
	// 向模板发送数据
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Reservation 渲染预定页面，展示表单
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation
	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{Form: forms.New(nil), Data: data})
}

// PostReservation 处理预定表单的post请求
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
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
	//如果form校验未通过,则渲染make-reservation模板并返回
	form.Required("first_name", "last_name", "email")
	//form.MinLength("first_name", 3, r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{Form: form, Data: data})
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation) //将一个名为“reservation”的变量存储在应用程序的会话中，以便在请求之间进行访问。

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther) // 重定向请求到“/reservation-summary”网址，并告诉浏览器以HTTP状态码“http.StatusSeeOther”的形式进行请求。
}

// Bigbed 渲染大床房页面，展示表单
func (m *Repository) Bigbed(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "bigbed.page.tmpl", &models.TemplateData{})
}

// Basic 渲染标准间页面，展示表单
func (m *Repository) Basic(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "basicroom.page.tmpl", &models.TemplateData{})
}

// Availability 渲染查找页面，展示表单
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
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

// AvailabilityJSON 处理查询请求并发送JSON响应
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{OK: true, Message: "Availability"}
	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		helpers.ServerError(w, err)
	}
	w.Header().Set("Content-Type", "application/json") // 使客户端能够正确判断和解析服务器返回的数据格式
	w.Write(out)
}

// Contact 渲染查找页面，展示表单
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation) // 从应用程序的会话中获取名为“reservation”的变量。如果变量存在，将尝试将其转换为类型“models.Reservation”
	if !ok {
		m.App.ErrorLog.Println("cannot get item from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "reservation") // 清除session
	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
