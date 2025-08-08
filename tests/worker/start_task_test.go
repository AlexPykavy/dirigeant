package worker

import (
	"context"
	"dirigeant/task"
	"dirigeant/tests/helper"
	"dirigeant/worker"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartTask__ShouldPersistTask(t *testing.T) {
	api := &worker.Api{
		Worker: worker.NewWorker(),
	}
	testTask := helper.PrintFileTask("print-task", helper.HostsFilePath)
	request := helper.NewTaskPostRequest(testTask)
	responseRecorder := httptest.NewRecorder()

	api.HandleCreateTask(responseRecorder, request)

	assert.Equal(t, http.StatusCreated, responseRecorder.Code, "Response status code should be 201 Created")
	assert.Empty(t, responseRecorder.Body, "Response body should be empty")
	assert.Equal(t, 1, api.Worker.LenTasks(), "Worker should have 1 task")

	persistedTask := api.Worker.GetTask(testTask.ID)

	assert.NotNil(t, persistedTask, "Persisted task ID should match the one from request")
	assert.Equal(t, task.Succeeded, persistedTask.Status, "Persisted task Status should be Succeeded")
}

func TestStartTask__ShouldReturnAnErrorIfCreatingTheSameTaskTwice(t *testing.T) {
	api := &worker.Api{
		Worker: worker.NewWorker(),
	}
	testTask := helper.PrintFileTask("print-task", helper.HostsFilePath)

	// 1 - Create a task for the first time
	firstRequest := helper.NewTaskPostRequest(testTask)
	firstResponseRecorder := httptest.NewRecorder()

	api.HandleCreateTask(firstResponseRecorder, firstRequest)

	assert.Equal(t, http.StatusCreated, firstResponseRecorder.Code, "Response status code should be 201 Created")
	assert.Empty(t, firstResponseRecorder.Body, "Response body should be empty")
	assert.Equal(t, 1, api.Worker.LenTasks(), "Worker should have 1 task")

	persistedTask := api.Worker.GetTask(testTask.ID)

	assert.NotNil(t, persistedTask, "Persisted task ID should match the one from request")
	assert.Equal(t, task.Succeeded, persistedTask.Status, "Persisted task Status should be Succeeded")

	// 2 - Create the same task for the second time
	secondRequest := helper.NewTaskPostRequest(testTask)
	secondResponseRecorder := httptest.NewRecorder()

	api.HandleCreateTask(secondResponseRecorder, secondRequest)

	assert.Equal(t, http.StatusConflict, secondResponseRecorder.Code, "Response status code should be 409 Conflict")
	assert.Equal(t, fmt.Sprintf("Error when executing the task: %s", task.ErrAlreadyExists), secondResponseRecorder.Body.String(), "Response body should contain error message")

	persistedTask = api.Worker.GetTask(testTask.ID)

	assert.NotNil(t, persistedTask, "Persisted task ID should match the one from request")
	assert.Equal(t, task.Succeeded, persistedTask.Status, "Persisted task Status should be Succeeded")
}

func TestStartTask__AllButOneRequestsShouldFailIfCreatingTheSameTaskSimultaneously(t *testing.T) {
	api := &worker.Api{
		Worker: worker.NewWorker(),
	}
	testTask := helper.PrintFileTask("print-task", helper.HostsFilePath)
	numOfRequests := 10
	requests := make([]*http.Request, numOfRequests)
	responseRecorders := make([]*httptest.ResponseRecorder, numOfRequests)

	var wg sync.WaitGroup
	for i := range numOfRequests {
		wg.Add(1)

		requests[i] = helper.NewTaskPostRequest(testTask)
		responseRecorders[i] = httptest.NewRecorder()

		go func() {
			defer wg.Done()

			api.HandleCreateTask(responseRecorders[i], requests[i])
		}()
	}

	wg.Wait()

	succeededRequests, conflictedRequests := 0, 0
	for i := range numOfRequests {
		switch responseRecorders[i].Code {
		case http.StatusCreated:
			succeededRequests++
		case http.StatusConflict:
			conflictedRequests++
		}
	}

	assert.Equal(t, 1, succeededRequests, "There should be only 1 succeeded request")
	assert.Equal(t, numOfRequests-1, conflictedRequests, "There should be only N-1 conflicted requests")
	assert.Equal(t, 1, api.Worker.LenTasks(), "Worker should have 1 task")

	persistedTask := api.Worker.GetTask(testTask.ID)

	assert.NotNil(t, persistedTask, "Persisted task ID should match the one from request")
	assert.Equal(t, task.Succeeded, persistedTask.Status, "Persisted task Status should be Succeeded")
}

func TestStartTask__ShouldHandleClientClosedRequest(t *testing.T) {
	api := &worker.Api{
		Worker: worker.NewWorker(),
	}
	testTask := helper.PingTask("ping-task", "127.0.0.1")
	ctx, cancel := context.WithCancel(context.TODO())

	// 1 - Create a task
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		createRequest := helper.NewTaskPostRequest(testTask).WithContext(ctx)
		createResponseRecorder := httptest.NewRecorder()

		stdout := helper.CaptureStdout(func() {
			api.HandleCreateTask(createResponseRecorder, createRequest)
		})

		assert.Equal(t, 499, createResponseRecorder.Code, "Response status code should be 499 Client Closed Request")
		assert.Equal(t, "Error when executing the task: client closed request", createResponseRecorder.Body.String(), "Response body should contain error message")
		assert.NotEmpty(t, stdout, "Task logs shouldn't be empty")
		assert.Zero(t, api.Worker.LenTasks(), "Worker should have no tasks")
	}()

	time.Sleep(1 * time.Second)

	assert.Equal(t, 1, api.Worker.LenTasks(), "Worker should have 1 task")

	persistedTask := api.Worker.GetTask(testTask.ID)

	assert.NotNil(t, persistedTask, "Persisted task ID should match the one from request")
	assert.Equal(t, task.Running, persistedTask.Status, "Persisted task Status should be Running")

	// 2 - Cancel a request
	cancel()

	wg.Wait()
}
