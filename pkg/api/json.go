package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/exis7ence/final_project/pkg/db"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	body, err := json.Marshal(data)
	if err != nil {
		log.Printf("ошибка сериализации JSON: %v", err)
		http.Error(w, "Внутренняя ошибка", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	if _, err := w.Write(append(body, '\n')); err != nil {
		log.Printf("ошибка отправки JSON-ответа: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}

func writeInternalError(w http.ResponseWriter, err error) {
	log.Printf("внутренняя ошибка: %v", err)

	writeError(
		w,
		http.StatusInternalServerError,
		"Внутренняя ошибка",
	)
}

func writeDBError(w http.ResponseWriter, err error) {
	if errors.Is(err, db.ErrTaskNotFound) {
		writeError(
			w,
			http.StatusNotFound,
			db.ErrTaskNotFound.Error(),
		)
		return
	}

	writeInternalError(w, err)
}