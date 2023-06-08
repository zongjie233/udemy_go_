package render

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/zongjie233/udemy_lesson/internal/config"
	"github.com/zongjie233/udemy_lesson/internal/models"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{}) // 将数据类型“models.Reservation”注册到gob包中，允许以二进制格式进行编码和解码。

	testApp.InProduction = false // 在生产模式时请设置为true

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true // 关闭浏览器之后继续保留
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = testApp.InProduction

	app = &testApp

	testApp.Session = session
	os.Exit(m.Run()) // 这个函数实在任何测试之前调用的
}
