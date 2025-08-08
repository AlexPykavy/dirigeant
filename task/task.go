package task

import (
	"errors"
	"os/exec"

	"github.com/google/uuid"
)

type TaskStatus int

const (
	Running TaskStatus = iota
	Succeeded
	Stopped
	Failed
)

var (
	ErrAlreadyExists = errors.New("task already exists")
	ErrNotExists     = errors.New("task does not exist")
)

type Task struct {
	ID         uuid.UUID
	Name       string
	Executable string
	Args       []string

	Status TaskStatus

	Cmd *exec.Cmd
}
