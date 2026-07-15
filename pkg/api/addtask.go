package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/exis7ence/final_project/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(
			w,
			fmt.Errorf("ошибка чтения JSON: %w", err),
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

	if err := checkDate(&task); err != nil {
		writeError(w, err)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, map[string]string{
		"id": strconv.FormatInt(id, 10),
	})
}

// checkDate проверяет дату и правило повторения задачи
//
// Если дата не указана, используется сегодняшнее число
// Если дата уже прошла
//   - обычная задача переносится на сегодня
//   - повторяющаяся задача переносится на следующую допустимую дату
func checkDate(task *db.Task) error {
	now := dateOnly(time.Now())

	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}

	taskDate, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("некорректная дата задачи: %w", err)
	}

	var nextDate string

	// Вызов NextDate нужен не только для вычисления даты
	// но и для проверки правильности правила повторения
	if task.Repeat != "" {
		nextDate, err = NextDate(
			now,
			task.Date,
			task.Repeat,
		)
		if err != nil {
			return fmt.Errorf(
				"некорректное правило повторения: %w",
				err,
			)
		}
	}

	// Сравниваем календарные даты без учёта текущего времени
	if taskDate.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format(DateFormat)
		} else {
			task.Date = nextDate
		}
	}

	return nil
}