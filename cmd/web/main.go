package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/zongjie233/udemy_lesson/internal/config"
	"github.com/zongjie233/udemy_lesson/internal/driver"
	"github.com/zongjie233/udemy_lesson/internal/handlers"
	"github.com/zongjie233/udemy_lesson/internal/helpers"
	"github.com/zongjie233/udemy_lesson/internal/models"
	"github.com/zongjie233/udemy_lesson/internal/render"
	"log"
	"net/http"
	"os"
	"time"
)

const portNumber = ":8080"

var app config.AppConfig        // 声明应该在main函数外，这样同为main包下的middleware也能使用声明的配置文件
var session *scs.SessionManager // 便于管理session
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err) // Fatal 会停止应用
	}
	defer db.SQL.Close()

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
func run() (*driver.DB, error) {
	gob.Register(models.Reservation{}) // 将数据类型“models.Reservation”注册到gob包中，允许以二进制格式进行编码和解码。
	gob.Register(models.Room{})
	gob.Register(models.RoomRestriction{})
	gob.Register(models.User{})
	gob.Register(models.Restriction{})

	app.InProduction = false // 在生产模式时请设置为true

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true // 关闭浏览器之后继续保留
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	// 链接到数据库
	log.Println("连接数据库中...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=990321")
	if err != nil {
		log.Fatal("连接不到数据库")
	}

	log.Println("成功连接数据库！")

	app.Session = session
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}
	app.TemplateCache = tc
	app.UseCache = true

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)
	return db, nil
}
