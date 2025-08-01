package worker

import (
	"dirigeant/task"
	"fmt"
	"iter"
	"maps"
	"os"
	"os/exec"

	"github.com/google/uuid"
)

type Worker struct {
	Tasks map[uuid.UUID]*task.Task
}

func (w *Worker) ListTasks() iter.Seq[*task.Task] {
	return maps.Values(w.Tasks)
}

func (w *Worker) GetTask(id uuid.UUID) *task.Task {
	return w.Tasks[id]
}

func (w *Worker) StartTask(t task.Task) error {
	cmd := exec.Command(t.Executable, t.Args...)
	t.Process = cmd.Process

	stdout, err := cmd.CombinedOutput()
	os.Stdout.Write(stdout)
	if err != nil {
		return err
	}

	w.Tasks[t.ID] = &t

	return nil
}

func (w *Worker) StopTask(id uuid.UUID) error {
	t := w.GetTask(id)
	if t == nil {
		return fmt.Errorf("%s not found", id)
	}

	return t.Process.Kill()
}
