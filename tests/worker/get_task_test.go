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

func TestGetTask__ShouldReturnAnErrorIfNotFound(t *testing.T) {
	api := &worker.Api{
		Worker: &worker.Worker{},
	}
	id := uuid.New()
	request := helper.NewTaskGetRequest(id)
	responseRecorder := httptest.NewRecorder()

	api.HandleGetTask(responseRecorder, request)

	assert.Equal(t, http.StatusNotFound, responseRecorder.Code, "Response status code should be 404 Not Found")
	assert.Equal(t, fmt.Sprintf("A task with %s ID not found", id), responseRecorder.Body.String(), "Response body should contain error message")
}

func TestGetTask__ShouldReturnFinishedTask(t *testing.T) {
	tcs := []struct {
		name           string
		path           string
		responseStatus int
		responseBody   string
		taskStatus     task.TaskStatus
	}{
		{
			name:           "print-hosts-file",
			path:           helper.HostsFilePath,
			responseStatus: http.StatusCreated,
			responseBody:   "",
			taskStatus:     task.Succeeded,
		},
		{
			name:           "print-non-existing-file",
			path:           "non-existing-file.txt",
			responseStatus: http.StatusInternalServerError,
			responseBody:   "Error when executing the task: exit status 1",
			taskStatus:     task.Failed,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			api := &worker.Api{
				Worker: worker.NewWorker(),
			}
			testTask := helper.PrintFileTask(tc.name, tc.path)

			// 1 - create a task
			createRequest := helper.NewTaskPostRequest(testTask)
			createResponseRecorder := httptest.NewRecorder()

			api.HandleCreateTask(createResponseRecorder, createRequest)

			assert.Equal(t, tc.responseStatus, createResponseRecorder.Code)
			assert.Equal(t, tc.responseBody, createResponseRecorder.Body.String())
			assert.Equal(t, 1, api.Worker.LenTasks(), "Worker should have 1 task")

			persistedTask := api.Worker.GetTask(testTask.ID)

			assert.NotNil(t, persistedTask, "Persisted task ID should match the one from request")
			assert.Equal(t, tc.taskStatus, persistedTask.Status)

			// 2 - get a task
			getRequest := helper.NewTaskGetRequest(testTask.ID)
			getResponseRecorder := httptest.NewRecorder()

			api.HandleGetTask(getResponseRecorder, getRequest)

			assert.Equal(t, http.StatusOK, getResponseRecorder.Code, "Response status code should be 200 OK")

			gottenTask := helper.JsonDecodeTask(getResponseRecorder.Body)

			assert.Equal(t, testTask.ID, gottenTask.ID, "Gotten task ID should match the created one")
			assert.Equal(t, tc.taskStatus, gottenTask.Status)
		})
	}
}

func TestGetTask__ShouldReturnRunningTask(t *testing.T) {
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

		assert.Equal(t, http.StatusCreated, createResponseRecorder.Code, "Response status code should be 201 Created")
		assert.Empty(t, createResponseRecorder.Body, "Response body should be empty")
		assert.Equal(t, 1, api.Worker.LenTasks(), "Worker should have 1 task")

		persistedTask := api.Worker.GetTask(testTask.ID)

		assert.NotNil(t, persistedTask, "Persisted task ID should match the one from request")
		assert.Equal(t, task.Succeeded, persistedTask.Status, "Persisted task status should be Succeeded")
	}()

	time.Sleep(1 * time.Second)

	assert.Equal(t, 1, api.Worker.LenTasks(), "Worker should have 1 task")

	persistedTask := api.Worker.GetTask(testTask.ID)

	assert.NotNil(t, persistedTask, "Persisted task ID should match the one from request")
	assert.Equal(t, task.Running, persistedTask.Status, "Persisted task Status should be Running")

	// 2 - Get a task
	getRequest := helper.NewTaskGetRequest(testTask.ID)
	getResponseRecorder := httptest.NewRecorder()

	api.HandleGetTask(getResponseRecorder, getRequest)

	assert.Equal(t, http.StatusOK, getResponseRecorder.Code, "Response status code should be 200 OK")

	gottenTask := helper.JsonDecodeTask(getResponseRecorder.Body)

	assert.Equal(t, testTask.ID, gottenTask.ID, "Gotten task ID should match the created one")
	assert.Equal(t, task.Running, gottenTask.Status)

	wg.Wait()
}
