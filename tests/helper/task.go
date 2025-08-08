package helper

import (
	"bytes"
	"context"
	"dirigeant/task"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var (
	HostsFilePath          string
	NoFileErrMessage       string
	SignalKilledErrMessage string
)

func init() {
	switch runtime.GOOS {
	case "windows":
		HostsFilePath = "$env:windir/System32/drivers/etc/hosts"
		NoFileErrMessage = "Cannot find path (\\r\\n)*'.+%s' because (\\r\\n)*it (\\r\\n)*does (\\r\\n)*not (\\r\\n)*exist."
		SignalKilledErrMessage = "exit status 1"
	case "linux":
		HostsFilePath = "/etc/hosts"
		NoFileErrMessage = "%s: No such file or directory"
		SignalKilledErrMessage = "signal: killed"
	}
}

func PingTask(name, host string) task.Task {
	switch runtime.GOOS {
	case "windows":
		return task.Task{
			ID:         uuid.New(),
			Name:       name,
			Executable: "ping",
			Args:       []string{"-n", "20", host},
		}
	case "linux":
		fallthrough
	default:
		return task.Task{
			ID:         uuid.New(),
			Name:       name,
			Executable: "ping",
			Args:       []string{"-c", "20", host},
		}
	}
}

func PrintFileTask(name, path string) task.Task {
	switch runtime.GOOS {
	case "windows":
		return task.Task{
			ID:         uuid.New(),
			Name:       name,
			Executable: "powershell",
			Args:       []string{"Get-Content", path},
		}
	case "linux":
		fallthrough
	default:
		return task.Task{
			ID:         uuid.New(),
			Name:       name,
			Executable: "cat",
			Args:       []string{path},
		}
	}
}

func JsonEncodeTask(t task.Task) io.Reader {
	w := &bytes.Buffer{}
	json.NewEncoder(w).Encode(t)
	return w
}

func NewTaskGetRequest(id uuid.UUID) *http.Request {
	if id == uuid.Nil {
		return httptest.NewRequest("GET", "/tasks", nil)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id.String())

	r := httptest.NewRequest("GET", fmt.Sprintf("/tasks/%s", id), nil)

	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func NewTaskPostRequest(t task.Task) *http.Request {
	return httptest.NewRequest("POST", "/tasks", JsonEncodeTask(t))
}

func NewTaskDeleteRequest(id uuid.UUID) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id.String())

	r := httptest.NewRequest("DELETE", fmt.Sprintf("/tasks/%s", id), nil)

	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}
