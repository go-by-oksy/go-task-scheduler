package api

import (
	"database/sql"
	"errors"
	"github.com/go-by-oksy/go-task-scheduler/pkg/db"
	"net/http"
	"strconv"
	"time"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "некорректный идентификатор задачи",
		})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, map[string]string{
				"error": "Задача не найдена",
			})
			return
		}

		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Внутренняя ошибка сервера",
		})
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writeJSON(w, http.StatusNotFound, map[string]string{
					"error": "Задача не найдена",
				})
				return
			}

			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "Внутренняя ошибка сервера",
			})
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{})
		return
	}

	next, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{})
}
