package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/exis7ence/final_project/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	writeJSON(w, task)
}