package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/zongjie233/udemy_lesson/internal/config"
	"github.com/zongjie233/udemy_lesson/internal/handlers"
	"github.com/zongjie233/udemy_lesson/internal/models"
	"github.com/zongjie233/udemy_lesson/internal/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

var app config.AppConfig        // 声明应该在main函数外，这样同为main包下的middleware也能使用声明的配置文件
var session *scs.SessionManager // 便于管理session

// main is the main function
func main() {

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
	app.UseCache = true

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
