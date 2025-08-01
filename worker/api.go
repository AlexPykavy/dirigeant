package worker

import (
	"dirigeant/task"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Api struct {
	Worker *Worker
}

func (a *Api) HandleListTasks(w http.ResponseWriter, _ *http.Request) {
	tasks := slices.Collect(a.Worker.ListTasks())
	if tasks == nil {
		tasks = []*task.Task{}
	}

	responseBody, err := json.Marshal(tasks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshalling the response: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
	// json.NewEncoder(w).Encode(a.Worker.ListTasks())
}

func (a *Api) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	taskId, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing passed ID: %v", err)
		return
	}
	t := a.Worker.GetTask(taskId)

	if t == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "A task with %s ID not found", taskId)
		return
	}

	responseBody, err := json.Marshal(t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshalling the response: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
	// json.NewEncoder(w).Encode(a.Worker.ListTasks())
}

func (a *Api) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	t := task.Task{}
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request body: %v", err)
		return
	}

	if err := a.Worker.StartTask(t); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error when executing the task: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *Api) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	taskId, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing passed ID: %v", err)
		return
	}

	if err := a.Worker.StopTask(taskId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error when stopping the task: %v", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
