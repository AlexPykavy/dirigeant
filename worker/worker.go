package worker

import (
	"dirigeant/task"
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
	t.Cmd = exec.Command(t.Executable, t.Args...)
	w.Tasks[t.ID] = &t

	stdout, err := t.Cmd.CombinedOutput()
	os.Stdout.Write(stdout)
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) StopTask(id uuid.UUID) error {
	t := w.GetTask(id)
	if t == nil {
		return task.ErrNotExists
	}

	if t.Cmd.ProcessState == nil {
		if err := t.Cmd.Process.Kill(); err != nil {
			return err
		}
	}

	delete(w.Tasks, t.ID)

	return nil
}
