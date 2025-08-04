package worker

import (
	"dirigeant/task"
	"dirigeant/tests/helper"
	"dirigeant/worker"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestListTasks__ShouldReturnAnEmptySlice(t *testing.T) {
	api := &worker.Api{
		Worker: &worker.Worker{},
	}
	request := helper.NewTaskGetRequest(uuid.Nil)
	responseRecorder := httptest.NewRecorder()

	api.HandleListTasks(responseRecorder, request)

	tasks := []task.Task{}
	json.NewDecoder(responseRecorder.Body).Decode(&tasks)
	assert.Equal(t, http.StatusOK, responseRecorder.Code, "Response status code should be 200 OK")
	assert.Empty(t, tasks, "Response body should be an empty slice")
}
