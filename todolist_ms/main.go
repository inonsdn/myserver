package main

import (
	"fmt"
	"sync"
	"time"
	"todolistms/internal/todo"
)

func todoFuncMock(i int) func() {
	return func() {
		fmt.Println("Execute todo function", i)
	}
}

func testArray() {
	s := []int{}
	a := [20]int{}
	now := time.Now()
	for i := range 20 {
		s = append(s, i)
		fmt.Println(len(s), cap(s))
	}
	fmt.Println("s takes", time.Since(now))
	now = time.Now()
	for i := range 20 {
		a[i] = i
	}
	fmt.Println("a takes", time.Since(now))
	fmt.Println("array a:", a)
}

func main() {
	opts := GetOptions()
	fmt.Println(opts.GetOptionsInfo())

	s1 := todo.NewScheduler(2)
	s2 := todo.NewScheduler(2)

	wg := sync.WaitGroup{}

	wg.Add(2)

	nowTime := time.Now()

	todo1 := todo.NewTodo("test1", todoFuncMock(1))
	todo2 := todo.NewTodo("test2", todoFuncMock(2))
	todo3 := todo.NewTodo("test3", todoFuncMock(3))
	todo4 := todo.NewTodo("test4", todoFuncMock(4))
	s1.AddTodo(todo1)
	s1.AddTodo(todo2)
	s2.AddTodo(todo3)
	s2.AddTodo(todo4)
	go s1.RunScheduler(&wg)
	go s2.RunScheduler(&wg)

	time.Sleep(300 * time.Millisecond)

	s1.Cancel()
	s2.Cancel()

	wg.Wait()

	fmt.Println(time.Since(nowTime))

	testArray()
}
