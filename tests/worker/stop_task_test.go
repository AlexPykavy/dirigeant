package worker

import (
	"dirigeant/task"
	"dirigeant/tests/helper"
	"dirigeant/worker"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

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
		Worker: worker.NewWorker(),
	}
	testTask := helper.PrintFileTask("print-task", helper.HostsFilePath)

	// 1 - Create a task
	request := helper.NewTaskPostRequest(testTask)
	responseRecorder := httptest.NewRecorder()
	api.HandleCreateTask(responseRecorder, request)

	assert.Equal(t, http.StatusCreated, responseRecorder.Code, "Response status code should be 201 Created")
	assert.Empty(t, responseRecorder.Body, "Response body should be empty")
	assert.Equal(t, 1, api.Worker.LenTasks(), "Tasks map should contain 1 task")

	// 2 - Delete a task
	request = helper.NewTaskDeleteRequest(testTask.ID)
	responseRecorder = httptest.NewRecorder()

	api.HandleDeleteTask(responseRecorder, request)

	assert.Equal(t, http.StatusNoContent, responseRecorder.Code, "Response status code should be 204 No Content")
	assert.Empty(t, responseRecorder.Body, "Response body should be empty")
	assert.Zero(t, api.Worker.LenTasks(), "Tasks map should be empty")
}

func TestStopTask__ShouldStopRunningTask(t *testing.T) {
	api := &worker.Api{
		Worker: worker.NewWorker(),
	}
	testTask := helper.PingTask("ping-task", "127.0.0.1")

	// 1 - Create a task
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		createRequest := helper.NewTaskPostRequest(testTask)
		createResponseRecorder := httptest.NewRecorder()

		api.HandleCreateTask(createResponseRecorder, createRequest)

		assert.Equal(t, http.StatusInternalServerError, createResponseRecorder.Code, "Response status code should be 500 Internal Server Error")
		assert.Equal(t, fmt.Sprintf("Error when executing the task: %s", helper.SignalKilledErrMessage), createResponseRecorder.Body.String(), "Response body should contain error message")
		assert.Zero(t, api.Worker.LenTasks(), "Tasks map should be empty")
	}()

	time.Sleep(1 * time.Second)

	assert.Equal(t, 1, api.Worker.LenTasks(), "Tasks map should contain 1 task")

	// 2 - Delete a task
	deleteRequest := helper.NewTaskDeleteRequest(testTask.ID)
	deleteResponseRecorder := httptest.NewRecorder()

	api.HandleDeleteTask(deleteResponseRecorder, deleteRequest)

	assert.Equal(t, http.StatusNoContent, deleteResponseRecorder.Code, "Response status code should be 204 No Content")
	assert.Empty(t, deleteResponseRecorder.Body, "Response body should be empty")
	assert.Zero(t, api.Worker.LenTasks(), "Tasks map should be empty")

	wg.Wait()
}
