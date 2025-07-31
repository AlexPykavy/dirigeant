package worker

import (
	"dirigeant/task"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListTasks__ShouldReturnAnEmptyList(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/tasks")
	if err != nil {
		t.Error(err)
	}

	tasks := []task.Task{}
	json.NewDecoder(resp.Body).Decode(&tasks)

	assert.Equal(t, tasks, []task.Task{}, "Returned tasks list should an empty slice")
}
