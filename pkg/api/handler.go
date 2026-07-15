package api

import (
	"fmt"
	"net/http"
	"time"
)

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	nowValue := r.FormValue("now")
	if nowValue != "" {
		parsedNow, err := time.Parse(DateFormat, nowValue)
		if err != nil {
			http.Error(w, "некорректная дата now", http.StatusBadRequest)
			return
		}

		now = parsedNow
	}

	nextDate, err := NextDate(
		now,
		r.FormValue("date"),
		r.FormValue("repeat"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, nextDate)
}