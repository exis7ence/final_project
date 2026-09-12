package api

import (
	"net/http"
	"strings"

	"github.com/exis7ence/final_project/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	writeJSON(w, http.StatusOK, task)
}