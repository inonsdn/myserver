package httpcon

var HttpPaths = []Path{
	{
		Name:     "/",
		Callback: homeHandler,
	},
	{
		Name:     "/test",
		Callback: testHandler,
	},
}
