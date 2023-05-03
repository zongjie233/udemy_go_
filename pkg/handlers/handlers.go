package handlers

import (
	"github.com/zongjie233/udemy_lesson/pkg/render"
	"net/http"
)

// Home is the Home page handler
func Home(w http.ResponseWriter, r *http.Request) { // 必须有着两个参数
	render.RenderTemplate(w, "home.page.tmpl")
}

// About is the about page handler
func About(w http.ResponseWriter, r *http.Request) { // 必须有着两个参数
	render.RenderTemplate(w, "about.page.tmpl")
}
