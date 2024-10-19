package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"test/repository"
	taskservice "test/task-service"
)

type Handler struct {
	RP repository.RepositoryProcesser
}

type HandleProcesser interface {
	HandleTask(w http.ResponseWriter, r *http.Request)
	HandleDate(w http.ResponseWriter, r *http.Request)
}

func NewHandler() Handler {
	rp, err := repository.NewRepo()
	if err != nil {
		panic(err)
	}
	return Handler{RP: rp}
}

func (h Handler) HandleDate(w http.ResponseWriter, r *http.Request) {
	task := new(taskservice.Task)
	task.Date = r.FormValue("date")
	task.Repeat = r.FormValue("repeat")
	now := r.FormValue("now")
	nextDt, err := task.GetNextRepeatDateTest(now)
	if err != nil {
		JsonErr(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Write([]byte(nextDt))
}

func (h Handler) HandleTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var newTask taskservice.Task
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)

		if err != nil {
			JsonErr(w, http.StatusBadRequest, "1234455N")
			return
		}

		w.Header().Set("Content-type", "application/json")

		if err := json.Unmarshal(buf.Bytes(), &newTask); err != nil {
			JsonErr(w, http.StatusBadRequest, "Ошибка десериализации JSON")
			return
		}

		id, err := h.RP.AddTask(newTask)
		if err != nil {
			JsonErr(w, http.StatusBadRequest, err.Error())
			return
		}

		JsonResponse(w, http.StatusOK, id)
	default:
		return
	}
}

func JsonErr(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)

	errorResponse := map[string]string{
		"error": message,
	}

	// Сериализуем карту в JSON
	response, err := json.Marshal(errorResponse)
	if err != nil {
		// В случае ошибки сериализации возвращаем простую текстовую ошибку
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Записываем результат в http.ResponseWriter
	w.Write(response)
}

func JsonResponse(w http.ResponseWriter, statusCode int, id string) {
	w.WriteHeader(statusCode)

	errorResponse := map[string]string{
		"id": id,
	}

	// Сериализуем карту в JSON
	response, err := json.Marshal(errorResponse)
	if err != nil {
		// В случае ошибки сериализации возвращаем простую текстовую ошибку
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Записываем результат в http.ResponseWriter
	w.Write(response)
}
