package main

import (
	"fmt"
	"github.com/zongjie233/udemy_lesson/pkg/handlers"
	"net/http"
)

const portNumber = ":8080" // 全局变量

func main() {

	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)

	fmt.Println("starting a server")
	_ = http.ListenAndServe(portNumber, nil) // 启动监听端口
}
