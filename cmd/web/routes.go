package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/zongjie233/udemy_lesson/pkg/config"
	"github.com/zongjie233/udemy_lesson/pkg/handlers"
	"net/http"
)

// routes 路由管理
func routes(app *config.AppConfig) http.Handler {
	//mux := pat.New()
	//mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	//mux.Get("/about", http.HandlerFunc(handlers.Repo.About))
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/bigbed", handlers.Repo.Bigbed)
	mux.Get("/basicroom", handlers.Repo.Basic)
	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Get("/search-availability-json", handlers.Repo.AvailabilityJSON)

	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/make-reservation", handlers.Repo.Reservation)

	filesServer := http.FileServer(http.Dir("./static/"))             // 建一个文件服务器，用于提供静态文件服务
	mux.Handle("/static/*", http.StripPrefix("/static", filesServer)) // 注册路由规则，将/static开头的请求映射到 filesServer上
	return mux
}
