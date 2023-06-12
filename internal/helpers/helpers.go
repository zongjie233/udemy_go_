package helpers

import "github.com/zongjie233/udemy_lesson/internal/config"

var app *config.AppConfig

func NewHelpers(a *config.AppConfig) {
	app = a
}
