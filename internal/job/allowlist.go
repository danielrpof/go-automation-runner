package job

var AllowedCommands = map[string][]string{
	"say-hello": {"echo", "hello"},
}

var AllowedJobs = map[string][]string{
	"say_hello": {
		"sh", "-c", "echo hello",
	},
	"list_files": {
		"sh", "-c", "ls -la",
	},
}
