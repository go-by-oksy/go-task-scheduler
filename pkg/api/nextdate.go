package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("не указано правило повторения задачи")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", err
	}

	rules := strings.Fields(repeat)
	if len(rules) == 0 {
		return "", errors.New("не указано правило повторения задачи")
	}

	switch rules[0] {
	case "y":
		return nextYearDate(now, date, rules)

	case "d":
		return nextDayDate(now, date, rules)

	case "w":
		return nextWeekDate(now, date, rules)

	case "m":
		return nextMonthDate(now, date, rules)

	default:
		return "", errors.New("неизвестное правило повторения")
	}
}

func nextYearDate(now time.Time, date time.Time, rules []string) (string, error) {
	if len(rules) != 1 {
		return "", errors.New("для правила не нужны дополнительные параметры")
	}

	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			return date.Format(dateFormat), nil
		}
	}
}

func nextDayDate(now time.Time, date time.Time, rules []string) (string, error) {
	if len(rules) != 2 {
		return "", errors.New("необходимо указать количество дней")
	}

	days, err := strconv.Atoi(rules[1])
	if err != nil {
		return "", errors.New("количество дней должно быть числом")
	}

	if days < 1 || days > 400 {
		return "", errors.New("количество дней должно быть от 1 до 400")
	}

	for {
		date = date.AddDate(0, 0, days)
		if afterNow(date, now) {
			return date.Format(dateFormat), nil
		}
	}
}

func nextWeekDate(now time.Time, date time.Time, rules []string) (string, error) {
	if len(rules) != 2 {
		return "", errors.New("необходимо указать дни недели")
	}

	weekDays := make(map[int]bool)

	for _, value := range strings.Split(rules[1], ",") {
		day, err := strconv.Atoi(value)
		if err != nil {
			return "", errors.New("день недели должен быть числом от 1 до 7")
		}

		if day < 1 || day > 7 {
			return "", errors.New("день недели должен быть от 1 до 7")
		}

		weekDays[day] = true
	}

	for {
		date = date.AddDate(0, 0, 1)

		day := int(date.Weekday())
		if day == 0 {
			day = 7
		}

		if weekDays[day] && afterNow(date, now) {
			return date.Format(dateFormat), nil
		}
	}
}

func nextMonthDate(now time.Time, date time.Time, rules []string) (string, error) {
	if len(rules) < 2 || len(rules) > 3 {
		return "", errors.New("для правила необходимо указать дни месяца и при необходимости месяцы")
	}

	days, err := parseMonthDays(rules[1])
	if err != nil {
		return "", err
	}

	months := make(map[int]bool)
	if len(rules) == 3 {
		months, err = parseMonths(rules[2])
		if err != nil {
			return "", err
		}
	}

	for {
		date = date.AddDate(0, 0, 1)

		month := int(date.Month())
		if len(months) > 0 && !months[month] {
			continue
		}

		day := date.Day()
		lastDay := lastDayOfMonth(date)
		prevLastDay := lastDay - 1

		if days[day] || (days[-1] && day == lastDay) || (days[-2] && day == prevLastDay) {
			if afterNow(date, now) {
				return date.Format(dateFormat), nil
			}
		}
	}
}

func parseMonthDays(value string) (map[int]bool, error) {
	days := make(map[int]bool)

	for _, item := range strings.Split(value, ",") {
		day, err := strconv.Atoi(item)
		if err != nil {
			return nil, errors.New("день месяца должен быть числом от 1 до 31")
		}

		if day != -1 && day != -2 && (day < 1 || day > 31) {
			return nil, errors.New("день месяца должен быть от 1 до 31")
		}

		days[day] = true
	}

	return days, nil
}

func parseMonths(value string) (map[int]bool, error) {
	months := make(map[int]bool)

	for _, item := range strings.Split(value, ",") {
		month, err := strconv.Atoi(item)
		if err != nil {
			return nil, errors.New("месяц должен быть числом от 1 до 12")
		}

		if month < 1 || month > 12 {
			return nil, errors.New("месяц должен быть от 1 до 12")
		}

		months[month] = true
	}

	return months, nil
}

func lastDayOfMonth(date time.Time) int {
	nextMonth := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location())
	return nextMonth.Day()
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var (
		now time.Time
		err error
	)

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, "дата должна быть в формате ГГГГММДД", http.StatusBadRequest)
			return
		}
	}

	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextDate))
}
