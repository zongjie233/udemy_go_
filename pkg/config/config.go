package config

import (
	"github.com/alexedwards/scs/v2"
	"html/template"
	"log"
)

// AppConfig app的全局配置
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger         // 日志指针
	InProduction  bool                // 是否为生产模式
	Session       *scs.SessionManager // 方便handlers函数访问
}
