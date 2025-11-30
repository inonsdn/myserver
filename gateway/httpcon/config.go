package httpcon

var HttpPathConfigs = []HttpGroupPath{
	{
		name: "",
		paths: []HttpPath{
			{
				name:     "/ping",
				callback: pong,
				method:   routeMethod_GET,
			},
			{
				name:     "/getUser",
				callback: getUser,
				method:   routeMethod_GET,
			},
		},
	},
	{
		name: "/admin",
		paths: []HttpPath{
			{
				name:     "/getUser",
				callback: getAllUser,
				method:   routeMethod_GET,
			},
		},
	},
}
