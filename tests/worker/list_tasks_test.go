package worker

import (
	"dirigeant/task"
	"dirigeant/worker"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListTasks__ShouldReturnAnEmptySlice(t *testing.T) {
	request := httptest.NewRequest("GET", "/tasks", nil)
	responseRecorder := httptest.NewRecorder()

	api := &worker.Api{
		Worker: &worker.Worker{},
	}
	api.HandleListTasks(responseRecorder, request)

	tasks := []task.Task{}
	json.NewDecoder(responseRecorder.Body).Decode(&tasks)
	assert.Equal(t, responseRecorder.Code, http.StatusOK, "Response status code should be 200")
	assert.Equal(t, tasks, []task.Task{}, "Response body should be an empty slice")
}
