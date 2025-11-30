package httpcon

import (
	"fmt"
	"net/http"
)

type Path struct {
	Name     string
	Callback func(w http.ResponseWriter, r *http.Request)
}

type PathHandler interface {
	GetPathName() string
	Execute(w http.ResponseWriter, r *http.Request)
}

func (p *Path) GetPathName() string {
	return p.Name
}

func (p *Path) Execute(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w.Header(), "print", r.Body)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w.Header(), "Home", r.Body)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w.Header(), "Test", r.Body)
}
