package task

import (
	"context"
	"sync"
	"time"

	"github.com/go-cmd/cmd"
)

type TaskManager struct {
	mu      sync.Mutex
	tasks   map[string]*cmd.Cmd
	timeout time.Duration
}

func NewManager(timeout time.Duration) *TaskManager {
	return &TaskManager{
		tasks:   make(map[string]*cmd.Cmd),
		timeout: timeout,
	}
}

func (tm *TaskManager) Execute(taskID string, command string, args ...string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Create new command
	c := cmd.NewCmd(command, args...)
	tm.tasks[taskID] = c

	// Start execution with timeout
	ctx, cancel := context.WithTimeout(context.Background(), tm.timeout)
	defer cancel()

	statusChan := c.Start()

	select {
	case <-ctx.Done():
		c.Stop()
		return ctx.Err()
	case finalStatus := <-statusChan:
		if finalStatus.Error != nil {
			return finalStatus.Error
		}
		return nil
	}
}

func (tm *TaskManager) Status(taskID string) (cmd.Status, bool) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	task, exists := tm.tasks[taskID]
	if !exists {
		return cmd.Status{}, false
	}
	return task.Status(), true
}

func (tm *TaskManager) Stop(taskID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	task, exists := tm.tasks[taskID]
	if !exists {
		return nil
	}

	err := task.Stop()
	delete(tm.tasks, taskID)
	return err
}

func (tm *TaskManager) Purge() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for _, task := range tm.tasks {
		task.Stop()
	}
	tm.tasks = make(map[string]*cmd.Cmd)
}
