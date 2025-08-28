package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TaskHeap struct {
	tasks     []*Task
	positions map[int]int
}

func NewHeap() TaskHeap {
	return TaskHeap{
		tasks:     make([]*Task, 0),
		positions: make(map[int]int),
	}
}

func (h *TaskHeap) Add(task *Task) {
	h.tasks = append(h.tasks, task)

	lastIndex := len(h.tasks) - 1
	h.positions[task.Identifier] = lastIndex
	h.siftUp(lastIndex)
}

func (h *TaskHeap) siftUp(position int) int {
	parent := (position - 1) / 2
	for position > 0 && h.tasks[parent].Priority < h.tasks[position].Priority {
		currentTask := h.tasks[position]
		parentTask := h.tasks[parent]

		h.tasks[position], h.tasks[parent] = h.tasks[parent], h.tasks[position]

		h.positions[currentTask.Identifier],
			h.positions[parentTask.Identifier] = h.positions[parentTask.Identifier],
			h.positions[currentTask.Identifier]

		position = parent
		parent = (position - 1) / 2
	}
	return position
}

func (h *TaskHeap) IsEmpty() bool {
	return len(h.tasks) == 0
}

func (h *TaskHeap) PopMax() (*Task, error) {
	if h.IsEmpty() {
		return nil, errors.New("Heap is empty!")
	}

	lastIndex := len(h.tasks) - 1

	max := h.tasks[0]
	h.tasks[0] = h.tasks[lastIndex]
	h.tasks = h.tasks[:lastIndex]

	if !h.IsEmpty() {
		id := h.tasks[0].Identifier
		h.positions[id] = h.siftDown(0)
	}

	return max, nil
}

func (h *TaskHeap) siftDown(position int) int {
	lastIndex := len(h.tasks) - 1

	for {
		leftChild := 2*position + 1
		rightChild := 2*position + 2
		largestChild := position

		if leftChild <= lastIndex && h.tasks[leftChild].Priority > h.tasks[largestChild].Priority {
			largestChild = leftChild
		}

		if rightChild <= lastIndex && h.tasks[rightChild].Priority > h.tasks[largestChild].Priority {
			largestChild = rightChild
		}

		if largestChild == position {
			break
		}

		h.tasks[position], h.tasks[largestChild] = h.tasks[largestChild], h.tasks[position]
		position = largestChild
	}

	return position
}

func (h *TaskHeap) ChangeTaskPriority(id int, priority int) {
	position := h.positions[id]
	task := h.tasks[position]

	if task.Priority > priority {
		task.Priority = priority
		h.positions[id] = h.siftDown(position)
	} else {
		task.Priority = priority
		h.positions[id] = h.siftUp(position)
	}
}

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	taskHeap TaskHeap
}

func NewScheduler() Scheduler {
	return Scheduler{
		taskHeap: NewHeap(),
	}
}

func (s *Scheduler) AddTask(task *Task) {
	s.taskHeap.Add(task)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	s.taskHeap.ChangeTaskPriority(taskID, newPriority)
}

func (s *Scheduler) GetTask() *Task {
	if !s.taskHeap.IsEmpty() {
		task, _ := s.taskHeap.PopMax()
		return task
	} else {
		return nil
	}
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(&task1)
	scheduler.AddTask(&task2)
	scheduler.AddTask(&task3)
	scheduler.AddTask(&task4)
	scheduler.AddTask(&task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, *task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, *task)

	scheduler.ChangeTaskPriority(1, 100)

	task = scheduler.GetTask()
	assert.Equal(t, task1, *task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, *task)
}
