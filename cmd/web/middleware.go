package main

import (
	"github.com/justinas/nosurf"
	"net/http"
)

// NoSurf 将 CSRF 保护添加到所有 POST 请求中
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,                 // 防止跨域脚本攻击，只允许通过 HTTP 请求发送 cookies，并防止 JavaScript 访问该 cookie
		Path:     "/",                  // cookie 的请求路径，对于所有路径都适用
		Secure:   app.InProduction,     // 如果应用程序在生产模式下，则访问 cookie 必须通过 HTTPS 发送
		SameSite: http.SameSiteLaxMode, // 防止跨站点请求伪造（CSRF）攻击
	})
	return csrfHandler // 将 CSRF 插入下一个请求处理程序
}

// SessionLoad 自动加载并保存保存每次请求的会话
func SessionLoad(next http.Handler) http.Handler {

	// 返回一个新的请求处理程序，用于加载和保存当前请求的会话数据，并将会话令牌以 cookie 的形式与客户进行交流
	return session.LoadAndSave(next)
}
