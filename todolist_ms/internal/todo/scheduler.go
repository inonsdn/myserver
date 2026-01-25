package todo

import "sync"

type Scheduler struct {
	todoCh chan *Todo
	ch     chan int
}

func NewScheduler(buff int) *Scheduler {
	return &Scheduler{
		todoCh: make(chan *Todo, buff),
		ch:     make(chan int, 1),
	}
}

func (s *Scheduler) RunScheduler(wg *sync.WaitGroup) {
	defer close(s.ch)
	defer close(s.todoCh)
	defer wg.Done()
	for {
		select {
		case todo := <-s.todoCh:
			todo.RunTask()
		case <-s.ch:
			return
		}
	}
}

func (s *Scheduler) AddTodo(t *Todo) {
	s.todoCh <- t
}

func (s *Scheduler) Cancel() {
	s.ch <- 0
}
