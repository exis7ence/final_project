package api

import "net/http"

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)

	case http.MethodPost:
		addTaskHandler(w, r)

	case http.MethodPut:
		updateTaskHandler(w, r)

	case http.MethodDelete:
		deleteTaskHandler(w, r)

	default:
		w.Header().Set("Allow", "GET, POST, PUT, DELETE")

		writeError(
			w,
			http.StatusMethodNotAllowed,
			"метод запроса не поддерживается",
		)
	}
}