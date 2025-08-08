package worker

import (
	"dirigeant/task"
	"dirigeant/tests/helper"
	"dirigeant/worker"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskLogs__PrintFile(t *testing.T) {
	tcs := []struct {
		name           string
		path           string
		responseStatus int
		responseBody   string
		stdoutRegexp   string
		taskStatus     task.TaskStatus
	}{
		{
			name:           "print-hosts-file",
			path:           helper.HostsFilePath,
			responseStatus: http.StatusCreated,
			responseBody:   "",
			stdoutRegexp:   "localhost",
			taskStatus:     task.Succeeded,
		},
		{
			name:           "print-non-existing-file",
			path:           "non-existing-file.txt",
			responseStatus: http.StatusInternalServerError,
			responseBody:   "Error when executing the task: exit status 1",
			stdoutRegexp:   fmt.Sprintf(helper.NoFileErrMessage, "non-existing-file.txt"),
			taskStatus:     task.Failed,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			api := &worker.Api{
				Worker: worker.NewWorker(),
			}
			testTask := helper.PrintFileTask(tc.name, tc.path)
			request := helper.NewTaskPostRequest(testTask)
			responseRecorder := httptest.NewRecorder()

			stdout := helper.CaptureStdout(func() {
				api.HandleCreateTask(responseRecorder, request)
			})

			assert.Equal(t, tc.responseStatus, responseRecorder.Code)
			assert.Equal(t, tc.responseBody, responseRecorder.Body.String())
			assert.Equal(t, 1, api.Worker.LenTasks(), "Worker should have 1 task")
			assert.Regexp(t, tc.stdoutRegexp, stdout)

			persistedTask := api.Worker.GetTask(testTask.ID)

			assert.NotNil(t, persistedTask, "Persisted task ID should match the one from request")
			assert.Equal(t, tc.taskStatus, persistedTask.Status)
		})
	}
}
