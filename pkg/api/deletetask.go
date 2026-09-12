package api

import (
	"net/http"
	"strings"

	"github.com/exis7ence/final_project/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.FormValue("id"))
	if id == "" {
		writeError(
			w,
			http.StatusBadRequest,
			"не указан идентификатор задачи",
		)
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeDBError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}