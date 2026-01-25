package todo

import (
	"fmt"
	"time"
)

type TodoState struct {
	state int
}

func (t *TodoState) SetNew() {
	t.state = 0
}

func (t *TodoState) SetWorking() {
	t.state = 1
}

func (t *TodoState) SetHold() {
	t.state = 2
}

func (t *TodoState) SetDone() {
	t.state = 3
}

func (t *TodoState) SetCancel() {
	t.state = 4
}

type TodoTask interface {
	Run()
}

type TodoFunc func()

type Todo struct {
	title string
	start time.Time
	end   time.Time
	state TodoState
	task  TodoFunc
}

func NewTodo(title string, task TodoFunc) *Todo {
	state := TodoState{}
	state.SetNew()
	return &Todo{
		title: title,
		state: state,
		task:  task,
	}
}

func (t *Todo) SetState(state int) {
	switch state {
	case 0:
		t.state.SetNew()
	case 1:
		t.state.SetWorking()
	case 2:
		t.state.SetHold()
	case 3:
		t.state.SetDone()
	case 4:
		t.state.SetCancel()
	default:
		panic(fmt.Sprintf("Unsupported state %d", state))
	}
}

func (t *Todo) RunTask() {
	t.task()
	fmt.Println("Run task title", t.title)
}
