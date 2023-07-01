package config

import (
	"github.com/alexedwards/scs/v2"
	"github.com/zongjie233/udemy_lesson/internal/models"
	"html/template"
	"log"
)

// AppConfig app的全局配置
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger // 日志指针
	ErrorLog      *log.Logger
	InProduction  bool                // 是否为生产模式
	Session       *scs.SessionManager // 方便handlers函数访问
	MailChan      chan models.MailData
}
