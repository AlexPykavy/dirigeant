package worker

import (
	"dirigeant/task"
	"iter"
	"maps"
	"os"
	"os/exec"
	"sync"

	"github.com/google/uuid"
)

type Worker struct {
	sync.RWMutex

	Tasks map[uuid.UUID]*task.Task
}

func (w *Worker) ListTasks() iter.Seq[*task.Task] {
	w.RLock()
	defer w.RUnlock()

	return maps.Values(w.Tasks)
}

func (w *Worker) GetTask(id uuid.UUID) *task.Task {
	w.RLock()
	defer w.RUnlock()

	return w.Tasks[id]
}

func (w *Worker) StartTask(t task.Task) error {
	w.Lock()

	if _, ok := w.Tasks[t.ID]; ok {
		w.Unlock()
		return task.ErrAlreadyExists
	}

	t.Cmd = exec.Command(t.Executable, t.Args...)
	w.Tasks[t.ID] = &t

	w.Unlock()

	stdout, err := t.Cmd.CombinedOutput()
	os.Stdout.Write(stdout)
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) StopTask(id uuid.UUID) error {
	w.Lock()
	defer w.Unlock()

	t := w.Tasks[id]
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
