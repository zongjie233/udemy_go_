package main

import (
	"github.com/justinas/nosurf"
	"net/http"
)

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
