package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/exis7ence/final_project/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(
			w,
			errors.New("метод запроса не поддерживается"),
		)
		return
	}

	id := strings.TrimSpace(r.FormValue("id"))
	if id == "" {
		writeError(
			w,
			errors.New("не указан идентификатор задачи"),
		)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err)
		return
	}

	// Одноразовая задача после выполнения удаляется
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, map[string]any{})
		return
	}

	// Для периодической задачи вычисляем следующую дату
	nextDate, err := NextDate(
		time.Now(),
		task.Date,
		task.Repeat,
	)
	if err != nil {
		writeError(w, err)
		return
	}

	if err := db.UpdateDate(nextDate, id); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, map[string]any{})
}