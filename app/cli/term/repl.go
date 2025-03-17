package term

import "os"

var IsRepl = os.Getenv("GPT4CLI_REPL") != ""

func SetIsRepl(value bool) {
	IsRepl = value
}
