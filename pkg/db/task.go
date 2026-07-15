package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// Task содержит данные задачи планировщика
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет задачу в базу данных
// и возвращает идентификатор новой записи
func AddTask(task *Task) (int64, error) {
	if db == nil {
		return 0, errors.New("база данных не инициализирована")
	}

	query := `
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (?, ?, ?, ?)
	`

	result, err := db.Exec(
		query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
	)
	if err != nil {
		return 0, fmt.Errorf("не удалось добавить задачу: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf(
			"не удалось получить идентификатор задачи: %w",
			err,
		)
	}

	return id, nil
}

// Tasks возвращает ближайшие задачи
// отсортированные по дате и идентификатору
func Tasks(limit int) ([]*Task, error) {
	query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		ORDER BY date, id
		LIMIT ?
	`

	return selectTasks(query, limit)
}

// SearchTasks ищет задачи по заголовку, комментарию
// или по дате в формате 02.01.2006
func SearchTasks(limit int, search string) ([]*Task, error) {
	searchDate, err := time.Parse("02.01.2006", search)
	if err == nil {
		query := `
			SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE date = ?
			ORDER BY date, id
			LIMIT ?
		`

		return selectTasks(
			query,
			searchDate.Format("20060102"),
			limit,
		)
	}

	query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE title LIKE ? OR comment LIKE ?
		ORDER BY date, id
		LIMIT ?
	`

	searchPattern := "%" + search + "%"

	return selectTasks(
		query,
		searchPattern,
		searchPattern,
		limit,
	)
}

// selectTasks выполняет SELECT-запрос и преобразует строки БД
// в слайс структур Task
func selectTasks(query string, args ...any) ([]*Task, error) {
	if db == nil {
		return nil, errors.New("база данных не инициализирована")
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить задачи: %w", err)
	}
	defer rows.Close()

	// Создаём именно пустой слайс, а не nil
	// Тогда JSON будет {"tasks":[]}, а не {"tasks":null}
	tasks := make([]*Task, 0)

	for rows.Next() {
		var (
			task Task
			id   int64
		)

		err := rows.Scan(
			&id,
			&task.Date,
			&task.Title,
			&task.Comment,
			&task.Repeat,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"не удалось прочитать задачу: %w",
				err,
			)
		}

		task.ID = strconv.FormatInt(id, 10)
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"ошибка при чтении списка задач: %w",
			err,
		)
	}

	return tasks, nil
}

// GetTask возвращает задачу по её идентификатору
func GetTask(id string) (*Task, error) {
	if db == nil {
		return nil, errors.New("база данных не инициализирована")
	}

	query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE id = ?
	`

	var (
		task      Task
		numericID int64
	)

	err := db.QueryRow(query, id).Scan(
		&numericID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("задача не найдена")
	}
	if err != nil {
		return nil, fmt.Errorf("не удалось получить задачу: %w", err)
	}

	task.ID = strconv.FormatInt(numericID, 10)

	return &task, nil
}

// UpdateTask обновляет существующую задачу.
func UpdateTask(task *Task) error {
	if db == nil {
		return errors.New("база данных не инициализирована")
	}

	query := `
		UPDATE scheduler
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?
	`

	result, err := db.Exec(
		query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
		task.ID,
	)
	if err != nil {
		return fmt.Errorf("не удалось обновить задачу: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"не удалось определить количество изменённых задач: %w",
			err,
		)
	}

	if count == 0 {
		return errors.New("задача не найдена")
	}

	return nil
}

// DeleteTask удаляет задачу по идентификатору
func DeleteTask(id string) error {
	if db == nil {
		return errors.New("база данных не инициализирована")
	}

	result, err := db.Exec(
		`DELETE FROM scheduler WHERE id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("не удалось удалить задачу: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"не удалось определить количество удалённых задач: %w",
			err,
		)
	}

	if count == 0 {
		return errors.New("задача не найдена")
	}

	return nil
}

// UpdateDate изменяет дату выполнения задачи
func UpdateDate(nextDate string, id string) error {
	if db == nil {
		return errors.New("база данных не инициализирована")
	}

	result, err := db.Exec(
		`UPDATE scheduler SET date = ? WHERE id = ?`,
		nextDate,
		id,
	)
	if err != nil {
		return fmt.Errorf("не удалось обновить дату задачи: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"не удалось определить количество изменённых задач: %w",
			err,
		)
	}

	if count == 0 {
		return errors.New("задача не найдена")
	}

	return nil
}