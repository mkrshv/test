package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	GetTasksHandle(w http.ResponseWriter, r *http.Request)
	GetTaskHandle(w http.ResponseWriter, r *http.Request)
	PutTaskHandle(w http.ResponseWriter, r *http.Request)
	DoneTaskeHandle(w http.ResponseWriter, r *http.Request)
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
		h.PostHandle(w, r)
	case "GET":
		h.GetTaskHandle(w, r)
	case "PUT":
		h.PutTaskHandle(w, r)
	case "DELETE":
		h.DeleteTaskeHandle(w, r)
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

func (h Handler) PostHandle(w http.ResponseWriter, r *http.Request) {
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
}

func (h Handler) GetTasksHandle(w http.ResponseWriter, r *http.Request) {
	taskSLice, err := h.RP.GetTaskList()
	if err != nil {
		fmt.Println("qwe")
		JsonErr(w, http.StatusBadRequest, err.Error())
	}

	respMap := make(map[string][]taskservice.Task)
	respMap["tasks"] = taskSLice

	resp, err := json.Marshal(respMap)
	if err != nil {
		fmt.Println("123")
		JsonErr(w, http.StatusInternalServerError, err.Error())
	}
	k := make(map[string][]taskservice.Task)
	err = json.Unmarshal(resp, &k)
	if err != nil {
		fmt.Println("123")
	}
	fmt.Println(k)
	w.Write(resp)
}

func (h Handler) GetTaskHandle(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	task, err := h.RP.GetTask(id)
	if err != nil {
		JsonErr(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := json.Marshal(task)
	if err != nil {
		JsonErr(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Write(resp)
}

func (h Handler) PutTaskHandle(w http.ResponseWriter, r *http.Request) {

	var buf bytes.Buffer
	task := taskservice.Task{}
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		JsonErr(w, http.StatusInternalServerError, err.Error())
		fmt.Println(task)
		return
	}
	err = h.RP.UpdateTask(task)
	if err != nil {
		JsonErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	response := struct{}{}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (h Handler) DoneTaskeHandle(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	err := h.RP.DoneTask(id)
	if err != nil {
		JsonErr(w, http.StatusBadRequest, "wrong id")
		return
	}

	response := struct{}{}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h Handler) DeleteTaskeHandle(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	err := h.RP.DeleteTask(id)
	if err != nil {
		JsonErr(w, http.StatusBadRequest, "wrong id")
		return
	}
	response := struct{}{}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
