package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	Id   int       `json:"id"`
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}

type TaskStore struct {
	sync.Mutex

	tasks  map[int]Task
	nextId int
}

type taskServer struct {
	store *TaskStore
}

func NewTaskServer() *taskServer {
	store := New()
	return &taskServer{store: store}
}

func New() *TaskStore {
	ts := &TaskStore{}
	ts.tasks = make(map[int]Task)
	ts.nextId = 0
	return ts
}

// CreateTask создаёт новую задачу в хранилище.
func (ts *TaskStore) CreateTask(text string, tags []string, due time.Time) int {
	ts.Lock()
	defer ts.Unlock()
	task := Task{
		Id:   ts.nextId,
		Text: text,
		Tags: tags,
		Due:  due,
	}
	ts.tasks[ts.nextId] = task
	ts.nextId++
	return ts.nextId
}

// GetTask получает задачу из хранилища по ID. Если ID не существует -
// будет возвращена ошибка.
func (ts *TaskStore) GetTask(id int) (Task, error) {
	ts.Lock()
	defer ts.Unlock()
	t, ok := ts.tasks[id]
	if ok {
		return t, nil
	} else {
		return t, fmt.Errorf("id = %d not found", id)
	}
}

// DeleteTask удаляет задачу с заданным ID. Если ID не существует -
// будет возвращена ошибка.
func (ts *TaskStore) DeleteTask(id int) error {
	ts.Lock()
	defer ts.Unlock()
	_, ok := ts.tasks[id]
	if !ok {
		return fmt.Errorf("id = %d not found", id)
	} else {
		delete(ts.tasks, id)
		return nil
	}

}

// DeleteAllTasks удаляет из хранилища все задачи.
func (ts *TaskStore) DeleteAllTasks() error {
	ts.Lock()
	defer ts.Unlock()
	ts.tasks = make(map[int]Task)
	return nil
}

// GetAllTasks возвращает из хранилища все задачи в произвольном порядке.
func (ts *TaskStore) GetAllTasks() []Task {
	var tasks []Task

	for _, task := range ts.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// GetTasksByTag возвращает, в произвольном порядке, все задачи
// с заданным тегом.
func (ts *TaskStore) GetTasksByTag(tag string) []Task {
	ts.Lock()
	defer ts.Unlock()

	var tasks []Task

	for _, task := range ts.tasks {
		for _, taskTag := range task.Tags {
			if taskTag == tag {
				tasks = append(tasks, task)
			}
		}
	}
	return tasks
}

// GetTasksByDueDate возвращает, в произвольном порядке, все задачи, которые
// запланированы на указанную дату.
func (ts *TaskStore) GetTasksByDueDate(year int, month time.Month, day int) []Task {
	ts.Lock()
	defer ts.Unlock()

	var tasks []Task

	for _, task := range ts.tasks {
		y, m, d := task.Due.Date()
		if y == year && m == month && d == day {
			tasks = append(tasks, task)
		}
	}

	return tasks
}
