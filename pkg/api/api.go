package api

import "net/http"

// Init регистрирует обработчики API
func Init(mux *http.ServeMux) {
	// Вход не должен требовать предварительной аутентификации
	mux.HandleFunc("/api/signin", signInHandler)

	// Вычисление следующей даты оставляем открытым
	mux.HandleFunc("/api/nextdate", nextDateHandler)

	// Основные операции защищаем middleware
	mux.HandleFunc("/api/task", auth(taskHandler))
	mux.HandleFunc("/api/tasks", auth(tasksHandler))
	mux.HandleFunc("/api/task/done", auth(doneTaskHandler))
}