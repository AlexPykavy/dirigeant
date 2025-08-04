package task

import (
	"errors"
	"os/exec"

	"github.com/google/uuid"
)

var (
	ErrNotExists = errors.New("task does not exist")
)

type Task struct {
	ID         uuid.UUID
	Name       string
	Executable string
	Args       []string

	Cmd *exec.Cmd
}
