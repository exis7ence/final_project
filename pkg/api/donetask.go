package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/exis7ence/final_project/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)

		writeError(
			w,
			http.StatusMethodNotAllowed,
			"метод запроса не поддерживается",
		)
		return
	}

	id := strings.TrimSpace(r.FormValue("id"))
	if id == "" {
		writeError(
			w,
			http.StatusBadRequest,
			"не указан идентификатор задачи",
		)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeDBError(w, err)
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeDBError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{})
		return
	}

	nextDate, err := NextDate(
		time.Now(),
		task.Date,
		task.Repeat,
	)
	if err != nil {
		writeInternalError(w, err)
		return
	}

	if err := db.UpdateDate(nextDate, id); err != nil {
		writeDBError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}