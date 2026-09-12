package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/exis7ence/final_project/pkg/db"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("ошибка чтения JSON: %v", err),
		)
		return
	}

	if strings.TrimSpace(task.ID) == "" {
		writeError(
			w,
			http.StatusBadRequest,
			"не указан идентификатор задачи",
		)
		return
	}

	if strings.TrimSpace(task.Title) == "" {
		writeError(
			w,
			http.StatusBadRequest,
			"не указан заголовок задачи",
		)
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeDBError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}