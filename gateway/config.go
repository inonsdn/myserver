package main

import "github.com/inonsdn/myserver/httpcon"

// define config path that supported
// TODO: move to config file and support route group
var ConfigPaths = []httpcon.HttpPath{
	{
		name:     "/ping",
		callback: pong,
		method:   "GET",
	},
	{
		name:     "/getUser",
		callback: getUser,
		method:   "GET",
	},
}
