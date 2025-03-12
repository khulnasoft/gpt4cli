package fs

import (
	"os"
	"os/exec"
	"path/filepath"
	"gpt4cli-cli/term"
)

var Cwd string
var Gpt4cliDir string
var ProjectRoot string
var HomeGpt4cliDir string
var CacheDir string

var HomeDir string
var HomeAuthPath string
var HomeAccountsPath string

func init() {
	var err error
	Cwd, err = os.Getwd()
	if err != nil {
		term.OutputErrorAndExit("Error getting current working directory: %v", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		term.OutputErrorAndExit("Couldn't find home dir: %v", err.Error())
	}
	HomeDir = home

	if os.Getenv("GPT4CLI_ENV") == "development" {
		HomeGpt4cliDir = filepath.Join(home, ".gpt4cli-home-dev-v2")
	} else {
		HomeGpt4cliDir = filepath.Join(home, ".gpt4cli-home-v2")
	}

	// Create the home gpt4cli directory if it doesn't exist
	err = os.MkdirAll(HomeGpt4cliDir, os.ModePerm)
	if err != nil {
		term.OutputErrorAndExit(err.Error())
	}

	CacheDir = filepath.Join(HomeGpt4cliDir, "cache")
	HomeAuthPath = filepath.Join(HomeGpt4cliDir, "auth.json")
	HomeAccountsPath = filepath.Join(HomeGpt4cliDir, "accounts.json")

	err = os.MkdirAll(filepath.Join(CacheDir, "tiktoken"), os.ModePerm)
	if err != nil {
		term.OutputErrorAndExit(err.Error())
	}
	err = os.Setenv("TIKTOKEN_CACHE_DIR", CacheDir)
	if err != nil {
		term.OutputErrorAndExit(err.Error())
	}

	FindGpt4cliDir()
	if Gpt4cliDir != "" {
		ProjectRoot = Cwd
	}
}

func FindOrCreateGpt4cli() (string, bool, error) {
	FindGpt4cliDir()
	if Gpt4cliDir != "" {
		ProjectRoot = Cwd
		return Gpt4cliDir, false, nil
	}

	// Determine the directory path
	var dir string
	if os.Getenv("GPT4CLI_ENV") == "development" {
		dir = filepath.Join(Cwd, ".gpt4cli-dev-v2")
	} else {
		dir = filepath.Join(Cwd, ".gpt4cli-v2")
	}

	err := os.Mkdir(dir, os.ModePerm)
	if err != nil {
		return "", false, err
	}
	Gpt4cliDir = dir
	ProjectRoot = Cwd

	return dir, true, nil
}

func ProjectRootIsGitRepo() bool {
	if ProjectRoot == "" {
		return false
	}

	return IsGitRepo(ProjectRoot)
}

func IsGitRepo(dir string) bool {
	isGitRepo := false

	if isCommandAvailable("git") {
		// check whether we're in a git repo
		cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")

		cmd.Dir = dir

		err := cmd.Run()

		if err == nil {
			isGitRepo = true
		}
	}

	return isGitRepo
}

func FindGpt4cliDir() {
	Gpt4cliDir = findGpt4cli(Cwd)
}

func findGpt4cli(baseDir string) string {
	var dir string
	if os.Getenv("GPT4CLI_ENV") == "development" {
		dir = filepath.Join(baseDir, ".gpt4cli-dev-v2")
	} else {
		dir = filepath.Join(baseDir, ".gpt4cli-v2")
	}
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return dir
	}

	return ""
}

func isCommandAvailable(name string) bool {
	cmd := exec.Command(name, "--version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
