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
	"path/filepath"
	"time"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./templates" // 便于测试用例找到模板位置
var functions = template.FuncMap{}

func getRoutes() http.Handler {
	gob.Register(models.Reservation{}) // 将数据类型“models.Reservation”注册到gob包中，允许以二进制格式进行编码和解码。

	app.InProduction = false // 在生产模式时请设置为true

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true // 关闭浏览器之后继续保留
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")

	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := NewRepo(&app)
	NewHandlers(repo)

	render.NewTemplates(&app)

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/bigbed", Repo.Bigbed)
	mux.Get("/basicroom", Repo.Basic)
	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)

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

// CreateTemplateCache 创建模板缓存
func CreateTemplateCache() (map[string]*template.Template, error) {
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
