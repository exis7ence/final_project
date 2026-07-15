package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/exis7ence/final_project/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.FormValue("id"))
	if id == "" {
		writeError(
			w,
			errors.New("не указан идентификатор задачи"),
		)
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, map[string]any{})
}