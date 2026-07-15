package api

import (
	"encoding/json"
	"errors"
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
			fmt.Errorf("ошибка чтения JSON: %w", err),
		)
		return
	}

	if strings.TrimSpace(task.ID) == "" {
		writeError(
			w,
			errors.New("не указан идентификатор задачи"),
		)
		return
	}

	if strings.TrimSpace(task.Title) == "" {
		writeError(
			w,
			errors.New("не указан заголовок задачи"),
		)
		return
	}

	// Используем ту же проверку даты и правила повторения
	// которая применяется при добавлении задачи
	if err := checkDate(&task); err != nil {
		writeError(w, err)
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeError(w, err)
		return
	}

	// При успешном обновлении API должен вернуть пустой JSON
	writeJSON(w, map[string]any{})
}