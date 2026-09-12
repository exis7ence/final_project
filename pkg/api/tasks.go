package api

import (
	"net/http"
	"strings"

	"github.com/exis7ence/final_project/pkg/db"
)

const tasksLimit = 50

type TasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)

		writeError(
			w,
			http.StatusMethodNotAllowed,
			"метод запроса не поддерживается",
		)
		return
	}

	search := strings.TrimSpace(r.FormValue("search"))

	var (
		tasks []*db.Task
		err   error
	)

	if search == "" {
		tasks, err = db.Tasks(tasksLimit)
	} else {
		tasks, err = db.SearchTasks(tasksLimit, search)
	}

	if err != nil {
		writeInternalError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, TasksResponse{
		Tasks: tasks,
	})
}