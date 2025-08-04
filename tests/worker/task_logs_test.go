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

func TestTaskLogs__PrintFile(t *testing.T) {
	tcs := []struct {
		name           string
		path           string
		responseStatus int
		responseBody   string
		stdoutRegexp   string
	}{
		{
			name:           "print-hosts-file",
			path:           helper.HostsFilePath,
			responseStatus: http.StatusCreated,
			responseBody:   "",
			stdoutRegexp:   "localhost",
		},
		{
			name:           "print-non-existing-file",
			path:           "non-existing-file.txt",
			responseStatus: http.StatusInternalServerError,
			responseBody:   "Error when executing the task: exit status 1",
			stdoutRegexp:   fmt.Sprintf(helper.NoFileErrMessage, "non-existing-file.txt"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			api := &worker.Api{
				Worker: &worker.Worker{
					Tasks: make(map[uuid.UUID]*task.Task),
				},
			}
			request := helper.NewTaskPostRequest(helper.PrintFileTask(tc.name, tc.path))
			responseRecorder := httptest.NewRecorder()

			stdout := helper.CaptureStdout(func() {
				api.HandleCreateTask(responseRecorder, request)
			})

			assert.Equal(t, tc.responseStatus, responseRecorder.Code)
			assert.Equal(t, tc.responseBody, responseRecorder.Body.String())
			assert.Regexp(t, tc.stdoutRegexp, stdout)
		})
	}
}
