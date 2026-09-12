package server

import (
	"fmt"
	"net/http"

	"github.com/exis7ence/final_project/pkg/api"
)

const webDir = "web"

func Run(port string) error {
	mux := http.NewServeMux()

	// Сначала регистрируем API
	api.Init(mux)

	// Затем регистрируем раздачу файлов фронтенда
	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	address := ":" + port

	fmt.Printf("Сервер запущен: http://localhost:%s\n", port)

	return http.ListenAndServe(address, mux)
}