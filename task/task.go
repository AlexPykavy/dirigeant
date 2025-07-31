package task

import (
	"os"

	"github.com/google/uuid"
)

type Task struct {
	ID         uuid.UUID
	Name       string
	Executable string
	Args       []string

	Process *os.Process
}
