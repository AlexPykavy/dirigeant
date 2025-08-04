package helper

import (
	"bytes"
	"dirigeant/task"
	"encoding/json"
	"io"
	"runtime"

	"github.com/google/uuid"
)

var (
	HostsFilePath      string
	NoFileErrorMessage string
)

func init() {
	switch runtime.GOOS {
	case "windows":
		HostsFilePath = "$env:windir/System32/drivers/etc/hosts"
		NoFileErrorMessage = "Cannot find path (\\r\\n)*'.+%s' because (\\r\\n)*it (\\r\\n)*does (\\r\\n)*not (\\r\\n)*exist."
	case "linux":
		HostsFilePath = "/etc/hosts"
		NoFileErrorMessage = "%s: No such file or directory"
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
