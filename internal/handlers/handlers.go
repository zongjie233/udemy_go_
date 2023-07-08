package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
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
	"strconv"
	"strings"
	"time"
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

func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

// NewHandlers 为处理程序设置存储库
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the Home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) { // 必须有着两个参数
	fmt.Println("home  called")

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})

}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) { // 必须有着两个参数
	// 向模板发送数据
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Reservation 渲染预定页面，展示表单
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "cant get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cant find room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res
	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservation 处理预定表单的post请求
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cant parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	layout := "2006-01-02"

	starDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cant parse startdate")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cant parse enddate")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: starDate,
		EndDate:   endDate,
		RoomID:    roomID,
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
		http.Error(w, "my own error message", http.StatusSeeOther)
		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	// 存入数据库
	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cant insert reservation into database ")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     starDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cant insert room restriction")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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
	// 获取表单上的数据,类型都为string，故需要转换为time类型
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// 没有可用房间
	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "没有可用房间")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	// 创建预订对象，并将其存储在会话中的"reservation"变量中
	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})

}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// AvailabilityJSON 处理查询请求并发送JSON响应
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	// 解析request body, 为了函数测试
	err := r.ParseForm()
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "    ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "ERROR connecting to database",
		}

		out, _ := json.MarshalIndent(resp, "", "    ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	available, err := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Error connecting to Database",
		}

		out, _ := json.MarshalIndent(resp, "", "    ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	out, _ := json.MarshalIndent(resp, "", "    ")

	w.Header().Set("Content-Type", "application/json") // 使客户端能够正确判断和解析服务器返回的数据格式
	w.Write(out)
}

// Contact 渲染查找页面，展示表单
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// ReservationSummary 展示预定总汇
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

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// ChooseRoom 展示可以预定的房间
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {

	// chi API，URLParam返回http.Request对象的URL参数。
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// 从会话中获取名为 "reservation" 的变量，并尝试将其转换为类型 "models.Reservation"
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}

	// 将更新后的预订对象存储回会话中的 "reservation" 变量
	res.RoomID = roomID
	m.App.Session.Put(r.Context(), "reservation", res)

	// 重定向用户到创建预订页面
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// BookRoom 用url的参数创建session变量，并重定向到预定界面
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {

	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var res models.Reservation

	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.Room.RoomName = room.RoomName
	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// ShowLogin shows the login screen
func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostShowLogin handles logging the user in
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout logs a user out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}

// AdminAllReservations 展示所有订单
func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, r, "admin-all-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// AdminNewReservations 展示所有新订单
func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllNewReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, r, "admin-new-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// AdminShowReservation  展示订单
func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := exploded[3]

	stringsMap := make(map[string]string)
	stringsMap["src"] = src

	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "admin-reservations-show.page.tmpl", &models.TemplateData{
		StringMap: stringsMap,
		Data:      data,
	})
}

// AdminReservationsCalendar 展示订单日期
func (m *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-reservations-calendar.page.tmpl", &models.TemplateData{})
}
