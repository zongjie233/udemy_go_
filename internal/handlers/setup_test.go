package handlers

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
	"github.com/zongjie233/udemy_lesson/internal/config"
	"github.com/zongjie233/udemy_lesson/internal/models"
	"github.com/zongjie233/udemy_lesson/internal/render"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var functions = template.FuncMap{
	"humanDate":  render.HumanDate,
	"formatDate": render.FormatDate,
	"iterate":    render.Iterate,
	"add":        render.Add,
}

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{}) // 注册 models.Reservation 类型，使其能够序列化为 gob 格式

	gob.Register(models.Room{})
	gob.Register(models.RoomRestriction{})
	gob.Register(models.User{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	app.InProduction = false // 在生产模式时请设置为true

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true                  // 让浏览器关闭后继续保留session
	session.Cookie.SameSite = http.SameSiteLaxMode // 要求 cookie 只能被同一站点请求访问
	session.Cookie.Secure = app.InProduction       // 表示 cookie 只有通过 HTTPS 才能发送，如果应用程序在生产模式下，则应启用此选项

	app.Session = session // 存储当前会话

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan
	defer close(mailChan)
	listenForMail()

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")

	}
	app.TemplateCache = tc
	app.UseCache = true

	repo := NewTestRepo(&app)
	NewHandlers(repo)
	render.NewRenderer(&app)

	os.Exit(m.Run())
}

func listenForMail() {
	go func() {
		for {
			_ = <-app.MailChan
		}
	}()
}

func getRoutes() http.Handler {

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	//mux.Use(NoSurf)      // 给所有 POST 请求添加 CSRF 保护
	mux.Use(SessionLoad) // 加载保存会话的函数

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/bigbed", Repo.Bigbed)
	mux.Get("/basicroom", Repo.Basic)
	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	mux.Get("/contact", Repo.Contact)
	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	filesServer := http.FileServer(http.Dir("./static/"))             // 建一个文件服务器，用于提供静态文件服务
	mux.Handle("/static/*", http.StripPrefix("/static", filesServer)) // 注册路由规则，将/static开头的请求映射到 filesServer上
	return mux
}

// NoSurf adds CSRF to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad 加载保存每次请求的会话
func SessionLoad(next http.Handler) http.Handler {

	//LoadAndSave提供了中间件，自动加载和保存当前请求的会话数据，并将会话令牌以cookie的形式与客户进行交流。在一个cookie中与客户端进行沟通。
	return session.LoadAndSave(next)
}

// CreateTestTemplateCache 创建模板缓存
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := make(map[string]*template.Template)

	// 获取templates中所有*.page.tmpl文件
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates)) // filepath.Glob()用于返回与指定模式匹配的所有文件或目录的名称，以切片模式返回
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
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
