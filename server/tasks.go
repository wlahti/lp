package server

import (
	"sync"
	"time"
)

type Task struct {
	Timestamp time.Time
	Text      string
}

type Tasks []Task

type TasksCache struct {
	lock  sync.RWMutex
	tasks map[int]Tasks
}

func NewTasksCache() *TasksCache {
	return &TasksCache{
		tasks: map[int]Tasks{},
	}
}

// readTasks is run by the server in the background
// to continously read new entries from the JSON task stream
func (t *TasksCache) readTasks() error {
	// TODO implement me
	return nil
}

func (t *TasksCache) GetTasks(id int) []string {
	t.lock.RLock()
	defer t.lock.RUnlock()

	tasks := make([]string, len(t.tasks[id]))
	for i, task := range t.tasks[id] {
		tasks[i] = task.Text
	}

	return tasks
}
