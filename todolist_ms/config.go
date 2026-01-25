package main

import "fmt"

type Options struct {
	RunningPort int
	Host        string
	Port        int
	User        string
	Password    string
	DBName      string
}
type OptsFunc func(*Options)

func defaultOptions() Options {
	return Options{
		RunningPort: 8082,
		Host:        "localhost",
		Port:        3306,
		User:        "root",
		Password:    "rootPass",
		DBName:      "todo_db",
	}
}

func (o *Options) GetOptionsInfo() string {
	return fmt.Sprintf("Options: %+v\n", o)
}

func GetOptions(funcs ...OptsFunc) *Options {
	o := defaultOptions()
	for _, f := range funcs {
		f(&o)
	}
	return &o
}

func setRunningPort(port int) OptsFunc {
	return func(o *Options) {
		o.RunningPort = port
	}
}
