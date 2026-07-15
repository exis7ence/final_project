package api

import (
	"errors"
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
		writeError(
			w,
			errors.New("метод запроса не поддерживается"),
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
		writeError(w, err)
		return
	}

	writeJSON(w, TasksResponse{
		Tasks: tasks,
	})
}