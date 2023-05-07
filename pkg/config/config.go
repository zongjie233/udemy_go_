package config

import (
	"html/template"
	"log"
)

// AppConfig app的全局配置
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger // 日志指针
}
