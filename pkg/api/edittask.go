package api

import (
	"encoding/json"
	"final-project/pkg/db"
	"net/http"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if task.ID == "" {
		writeJSON(w, map[string]string{
			"error": "не указан идентификатор задачи",
		})
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{
			"error": "не указан заголовок задачи",
		})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, map[string]string{})
}
