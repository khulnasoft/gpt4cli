package types

import ignore "github.com/sabhiram/go-gitignore"

type ProjectPaths struct {
	ActivePaths    map[string]bool
	AllPaths       map[string]bool
	Gpt4cliIgnored *ignore.GitIgnore
	IgnoredPaths   map[string]string
}
