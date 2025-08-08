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

	tasks map[uuid.UUID]*task.Task
}

func NewWorker() *Worker {
	return &Worker{
		tasks: make(map[uuid.UUID]*task.Task),
	}
}

func (w *Worker) LenTasks() int {
	w.RLock()
	defer w.RUnlock()

	return len(w.tasks)
}

func (w *Worker) ListTasks() iter.Seq[*task.Task] {
	w.RLock()
	defer w.RUnlock()

	return maps.Values(w.tasks)
}

func (w *Worker) GetTask(id uuid.UUID) *task.Task {
	w.RLock()
	defer w.RUnlock()

	return w.tasks[id]
}

func (w *Worker) StartTask(t task.Task) error {
	w.Lock()

	if _, ok := w.tasks[t.ID]; ok {
		w.Unlock()
		return task.ErrAlreadyExists
	}

	t.Cmd = exec.Command(t.Executable, t.Args...)
	t.Status = task.Running
	w.tasks[t.ID] = &t

	w.Unlock()

	stdout, err := t.Cmd.CombinedOutput()
	os.Stdout.Write(stdout)
	if err != nil {
		t.Status = task.Failed
		return err
	}

	t.Status = task.Succeeded

	return nil
}

func (w *Worker) StopTask(id uuid.UUID) error {
	w.Lock()
	defer w.Unlock()

	t := w.tasks[id]
	if t == nil {
		return task.ErrNotExists
	}

	if t.Cmd.ProcessState == nil {
		if err := t.Cmd.Process.Kill(); err != nil {
			return err
		}

		t.Status = task.Stopped
	}

	delete(w.tasks, t.ID)

	return nil
}
