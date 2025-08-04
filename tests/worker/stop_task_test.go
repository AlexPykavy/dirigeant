package worker

import (
	"dirigeant/task"
	"dirigeant/tests/helper"
	"dirigeant/worker"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStopTask__ShouldReturnAnErrorIfNotFound(t *testing.T) {
	api := &worker.Api{
		Worker: &worker.Worker{},
	}
	request := helper.NewTaskDeleteRequest(uuid.New())
	responseRecorder := httptest.NewRecorder()

	api.HandleDeleteTask(responseRecorder, request)

	assert.Equal(t, http.StatusNotFound, responseRecorder.Code, "Response status code should be 404 Not Found")
	assert.Equal(t, fmt.Sprintf("Error when stopping the task: %v", task.ErrNotExists), responseRecorder.Body.String(), "Response body should contain error message")
}

func TestStopTask__ShouldStopCompletedTask(t *testing.T) {
	api := &worker.Api{
		Worker: &worker.Worker{
			Tasks: make(map[uuid.UUID]*task.Task),
		},
	}
	testTask := helper.PrintFileTask("print-task", helper.HostsFilePath)

	// 1 - Create a task
	request := helper.NewTaskPostRequest(testTask)
	responseRecorder := httptest.NewRecorder()
	api.HandleCreateTask(responseRecorder, request)

	assert.Equal(t, http.StatusCreated, responseRecorder.Code, "Response status code should be 201 Created")
	assert.Empty(t, responseRecorder.Body, "Response body should be empty")
	assert.Equal(t, 1, len(api.Worker.Tasks), "Tasks map should contain 1 task")

	// 2 - Delete a task
	request = helper.NewTaskDeleteRequest(testTask.ID)
	responseRecorder = httptest.NewRecorder()

	api.HandleDeleteTask(responseRecorder, request)

	assert.Equal(t, http.StatusNoContent, responseRecorder.Code, "Response status code should be 204 No Content")
	assert.Empty(t, responseRecorder.Body, "Response body should be empty")
	assert.Empty(t, api.Worker.Tasks, "Tasks map should be empty")
}
