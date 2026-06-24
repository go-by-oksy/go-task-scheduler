package api

import (
	"final-project/pkg/db"
	"net/http"
	"strconv"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "некорректный идентификатор задачи",
		})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Внутренняя ошибка сервера",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{})
}
