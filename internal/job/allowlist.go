package job

var AllowedCommands = map[string][]string{
	"say-hello": {"echo", "hello"},
	"sleep":     {"sleep", "40"},
	"timeout":   {"sh", "-c", "sleep 40"},
}

var AllowedJobs = map[string]string{
	"say_hello": "echo hello",
	"sleep":     "sleep 40",
	"timeout":   "sleep 40",
}
