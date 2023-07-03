package main

import (
	"github.com/justinas/nosurf"
	"github.com/zongjie233/udemy_lesson/internal/helpers"
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

// Auth 给定的代码定义了一个名为 Auth 的中间件函数，用于对传入的 HTTP 请求进行身份验证。它接受一个 http.Handler 作为输入，并返回一个新的 http.Handler，用于执行身份验证逻辑。
// 在 Auth 函数内部，它使用 helpers.IsAuthenticated 函数检查用户是否已经通过身份验证。如果用户未经过身份验证，它会使用 session.Put 方法在会话上下文中设置错误消息，并使用 http.Redirect 将用户重定向到登录页面。http.StatusSeeOther 状态码表示临时重定向。
// 如果用户已经通过身份验证，它会调用 next 处理程序的 ServeHTTP 方法继续处理请求。
// 这个中间件可以用于保护需要身份验证的路由。它确保只有经过身份验证的用户才能访问这些路由，否则将被重定向到登录页面。
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			session.Put(r.Context(), "error", "log in first!")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
