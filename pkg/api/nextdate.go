package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

// Григорианский календарь полностью повторяет свой цикл за 400 лет
// Ограничение защищает поиск даты от бесконечного цикла
const maxSearchDays = 366 * 400

// monthDays хранит правила для дней месяца
//
// exact - обычные дни от 1 до 31
// last - последний день месяца, то есть -1
// penultimate - предпоследний день месяца, то есть -2
type monthDays struct {
	exact       [32]bool
	last        bool
	penultimate bool
}

// NextDate вычисляет следующую дату выполнения задачи
//
// Возвращаемая дата всегда должна быть позже now
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	start, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("некорректная начальная дата: %w", err)
	}

	now = dateOnly(now)

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("правило повторения не указано")
	}

	var next time.Time

	switch parts[0] {
	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf(
				"некорректное правило повторения %q",
				repeat,
			)
		}

		next = nextYearDate(start, now)

	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf(
				"некорректное правило повторения %q",
				repeat,
			)
		}

		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("интервал дней должен быть числом")
		}

		if interval < 1 || interval > 400 {
			return "", errors.New(
				"интервал дней должен быть от 1 до 400",
			)
		}

		next = nextDayDate(start, now, interval)

	case "w":
		if len(parts) != 2 {
			return "", fmt.Errorf(
				"некорректное правило повторения %q",
				repeat,
			)
		}

		weekdays, err := parseWeekdays(parts[1])
		if err != nil {
			return "", err
		}

		next, err = nextWeekDate(start, now, weekdays)
		if err != nil {
			return "", err
		}

	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", fmt.Errorf(
				"некорректное правило повторения %q",
				repeat,
			)
		}

		days, err := parseMonthDays(parts[1])
		if err != nil {
			return "", err
		}

		var months [13]bool

		if len(parts) == 3 {
			months, err = parseMonths(parts[2])
			if err != nil {
				return "", err
			}
		} else {
			// Если месяцы не указаны, разрешаем все
			for month := 1; month <= 12; month++ {
				months[month] = true
			}
		}

		next, err = nextMonthDate(start, now, days, months)
		if err != nil {
			return "", err
		}

	default:
		return "", fmt.Errorf(
			"неподдерживаемое правило повторения %q",
			repeat,
		)
	}

	return next.Format(DateFormat), nil
}

// dateOnly убирает из времени часы, минуты и секунды
// Для вычислений нас интересует только календарная дата
func dateOnly(value time.Time) time.Time {
	return time.Date(
		value.Year(),
		value.Month(),
		value.Day(),
		0,
		0,
		0,
		0,
		time.UTC,
	)
}

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

// nextYearDate обрабатывает правило y
func nextYearDate(start, now time.Time) time.Time {
	date := start

	for {
		date = date.AddDate(1, 0, 0)

		if afterNow(date, now) {
			return date
		}
	}
}

// nextDayDate обрабатывает правило d N
func nextDayDate(start, now time.Time, interval int) time.Time {
	date := start

	for {
		date = date.AddDate(0, 0, interval)

		if afterNow(date, now) {
			return date
		}
	}
}

// parseWeekdays разбирает строку вида 1,4,5
func parseWeekdays(value string) ([8]bool, error) {
	var weekdays [8]bool

	values, err := parseNumberList(value)
	if err != nil {
		return weekdays, fmt.Errorf(
			"некорректные дни недели: %w",
			err,
		)
	}

	for _, day := range values {
		if day < 1 || day > 7 {
			return weekdays, errors.New(
				"день недели должен быть от 1 до 7",
			)
		}

		weekdays[day] = true
	}

	return weekdays, nil
}

// nextWeekDate обрабатывает правило w
//
// В Go воскресенье имеет значение 0, поэтому мы преобразуем его в 7
func nextWeekDate(
	start time.Time,
	now time.Time,
	weekdays [8]bool,
) (time.Time, error) {
	date := start

	// Для недельного правила достаточно продолжать поиск от более поздней
	// из двух дат: start или now
	if now.After(date) {
		date = now
	}

	// В течение следующих семи дней подходящий день обязательно найдётся
	for i := 0; i < 7; i++ {
		date = date.AddDate(0, 0, 1)

		weekday := int(date.Weekday())
		if weekday == 0 {
			weekday = 7
		}

		if weekdays[weekday] {
			return date, nil
		}
	}

	return time.Time{}, errors.New(
		"не удалось вычислить следующую дату",
	)
}

// parseMonthDays разбирает дни месяца
//
// Допустимы
// 1-31 - конкретный день
// -1 - последний день
// -2 - предпоследний день
func parseMonthDays(value string) (monthDays, error) {
	var days monthDays

	values, err := parseNumberList(value)
	if err != nil {
		return days, fmt.Errorf(
			"некорректные дни месяца: %w",
			err,
		)
	}

	for _, day := range values {
		if day < -2 || day == 0 || day > 31 {
			return days, errors.New(
				"день месяца должен быть от 1 до 31, -1 или -2",
			)
		}

		switch day {
		case -1:
			days.last = true

		case -2:
			days.penultimate = true

		default:
			days.exact[day] = true
		}
	}

	return days, nil
}

// parseMonths разбирает строку месяцев вида 1,3,6
func parseMonths(value string) ([13]bool, error) {
	var months [13]bool

	values, err := parseNumberList(value)
	if err != nil {
		return months, fmt.Errorf(
			"некорректные месяцы: %w",
			err,
		)
	}

	for _, month := range values {
		if month < 1 || month > 12 {
			return months, errors.New(
				"месяц должен быть от 1 до 12",
			)
		}

		months[month] = true
	}

	return months, nil
}

// nextMonthDate обрабатывает правило m
func nextMonthDate(
	start time.Time,
	now time.Time,
	days monthDays,
	months [13]bool,
) (time.Time, error) {
	date := start

	if now.After(date) {
		date = now
	}

	for i := 0; i < maxSearchDays; i++ {
		date = date.AddDate(0, 0, 1)

		if !afterNow(date, now) {
			continue
		}

		if !months[int(date.Month())] {
			continue
		}

		currentDay := date.Day()
		lastDay := lastDayOfMonth(date)

		if days.exact[currentDay] {
			return date, nil
		}

		if days.last && currentDay == lastDay {
			return date, nil
		}

		if days.penultimate && currentDay == lastDay-1 {
			return date, nil
		}
	}

	return time.Time{}, errors.New(
		"для правила не существует следующей даты",
	)
}

// lastDayOfMonth возвращает номер последнего дня текущего месяца
//
// Нулевой день следующего месяца - это последний день текущего месяца
func lastDayOfMonth(date time.Time) int {
	lastDate := time.Date(
		date.Year(),
		date.Month()+1,
		0,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	return lastDate.Day()
}

// parseNumberList преобразует строку 1,4,5 в []int{1, 4, 5}
func parseNumberList(value string) ([]int, error) {
	parts := strings.Split(value, ",")

	numbers := make([]int, 0, len(parts))

	for _, part := range parts {
		if part == "" {
			return nil, errors.New(
				"в списке обнаружено пустое значение",
			)
		}

		number, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf(
				"%q не является числом",
				part,
			)
		}

		numbers = append(numbers, number)
	}

	return numbers, nil
}