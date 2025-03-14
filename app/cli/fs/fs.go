package fs

import (
	"context"
	"encoding/json"
	"fmt"
	"gpt4cli/term"
	"gpt4cli/types"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/khulnasoft/gpt4cli/shared"
	ignore "github.com/sabhiram/go-gitignore"
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
		HomeGpt4cliDir = filepath.Join(home, ".gpt4cli-home-dev")
	} else {
		HomeGpt4cliDir = filepath.Join(home, ".gpt4cli-home")
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

	Gpt4cliDir = findGpt4cli(Cwd)
	if Gpt4cliDir != "" {
		ProjectRoot = Cwd
	}
}

func FindOrCreateGpt4cli() (string, bool, error) {
	Gpt4cliDir = findGpt4cli(Cwd)
	if Gpt4cliDir != "" {
		ProjectRoot = Cwd
		return Gpt4cliDir, false, nil
	}

	// Determine the directory path
	var dir string
	if os.Getenv("GPT4CLI_ENV") == "development" {
		dir = filepath.Join(Cwd, ".gpt4cli-dev")
	} else {
		dir = filepath.Join(Cwd, ".gpt4cli")
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

type ProjectPaths struct {
	ActivePaths    map[string]bool
	AllPaths       map[string]bool
	Gpt4cliIgnored *ignore.GitIgnore
	IgnoredPaths   map[string]string
}

func GetProjectPaths(baseDir string) (*ProjectPaths, error) {
	if ProjectRoot == "" {
		return nil, fmt.Errorf("no project root found")
	}

	return GetPaths(baseDir, ProjectRoot)
}

func GetPaths(baseDir, currentDir string) (*ProjectPaths, error) {
	ignored, err := GetGpt4cliIgnore(currentDir)

	if err != nil {
		return nil, err
	}

	allPaths := map[string]bool{}
	activePaths := map[string]bool{}

	allDirs := map[string]bool{}
	activeDirs := map[string]bool{}

	isGitRepo := IsGitRepo(baseDir)

	errCh := make(chan error)
	var mu sync.Mutex
	numRoutines := 0

	deletedFiles := map[string]bool{}

	if isGitRepo {

		// Use git status to find deleted files
		numRoutines++
		go func() {
			cmd := exec.Command("git", "rev-parse", "--show-toplevel")
			output, err := cmd.Output()
			if err != nil {
				errCh <- fmt.Errorf("error getting git root: %s", err)
				return
			}
			repoRoot := strings.TrimSpace(string(output))

			cmd = exec.Command("git", "status", "--porcelain")
			cmd.Dir = baseDir
			out, err := cmd.Output()
			if err != nil {
				errCh <- fmt.Errorf("error getting git status: %s", err)
			}

			lines := strings.Split(string(out), "\n")

			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "D ") {
					path := strings.TrimSpace(line[2:])
					absPath := filepath.Join(repoRoot, path)
					relPath, err := filepath.Rel(currentDir, absPath)
					if err != nil {
						errCh <- fmt.Errorf("error getting relative path: %s", err)
						return
					}
					deletedFiles[relPath] = true
				}
			}

			errCh <- nil
		}()

		// combine `git ls-files` and `git ls-files --others --exclude-standard`
		// to get all files in the repo

		numRoutines++
		go func() {
			// get all tracked files in the repo
			cmd := exec.Command("git", "ls-files")
			cmd.Dir = baseDir
			out, err := cmd.Output()

			if err != nil {
				errCh <- fmt.Errorf("error getting files in git repo: %s", err)
				return
			}

			files := strings.Split(string(out), "\n")

			mu.Lock()
			defer mu.Unlock()
			for _, file := range files {
				absFile := filepath.Join(baseDir, file)
				relFile, err := filepath.Rel(currentDir, absFile)

				if err != nil {
					errCh <- fmt.Errorf("error getting relative path: %s", err)
					return
				}

				if ignored != nil && ignored.MatchesPath(relFile) {
					continue
				}

				activePaths[relFile] = true

				parentDir := relFile
				for parentDir != "." && parentDir != "/" && parentDir != "" {
					parentDir = filepath.Dir(parentDir)
					activeDirs[parentDir] = true
				}
			}

			errCh <- nil
		}()

		// get all untracked non-ignored files in the repo
		numRoutines++
		go func() {
			cmd := exec.Command("git", "ls-files", "--others", "--exclude-standard")
			cmd.Dir = baseDir
			out, err := cmd.Output()

			if err != nil {
				errCh <- fmt.Errorf("error getting untracked files in git repo: %s", err)
				return
			}

			files := strings.Split(string(out), "\n")

			mu.Lock()
			defer mu.Unlock()
			for _, file := range files {
				absFile := filepath.Join(baseDir, file)
				relFile, err := filepath.Rel(currentDir, absFile)

				if err != nil {
					errCh <- fmt.Errorf("error getting relative path: %s", err)
					return
				}

				if ignored != nil && ignored.MatchesPath(relFile) {
					continue
				}

				activePaths[relFile] = true

				parentDir := relFile
				for parentDir != "." && parentDir != "/" && parentDir != "" {
					parentDir = filepath.Dir(parentDir)
					activeDirs[parentDir] = true
				}
			}

			errCh <- nil
		}()
	}

	// get all paths in the directory
	numRoutines++
	go func() {
		err = filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				if info.Name() == ".git" {
					return filepath.SkipDir
				}
				if info.Name() == ".gpt4cli" || info.Name() == ".gpt4cli-dev" {
					return filepath.SkipDir
				}

				relPath, err := filepath.Rel(currentDir, path)
				if err != nil {
					return err
				}

				allDirs[relPath] = true

				if ignored != nil && ignored.MatchesPath(relPath) {
					return filepath.SkipDir
				}
			} else {
				relPath, err := filepath.Rel(currentDir, path)
				if err != nil {
					return err
				}

				allPaths[relPath] = true

				if ignored != nil && ignored.MatchesPath(relPath) {
					return nil
				}

				if !isGitRepo {
					mu.Lock()
					defer mu.Unlock()
					activePaths[relPath] = true

					parentDir := relPath
					for parentDir != "." && parentDir != "/" && parentDir != "" {
						parentDir = filepath.Dir(parentDir)
						activeDirs[parentDir] = true
					}
				}
			}

			return nil
		})

		if err != nil {
			errCh <- fmt.Errorf("error walking directory: %s", err)
			return
		}

		errCh <- nil
	}()

	for i := 0; i < numRoutines; i++ {
		err := <-errCh
		if err != nil {
			return nil, err
		}
	}

	for dir := range allDirs {
		allPaths[dir] = true
	}

	for dir := range activeDirs {
		activePaths[dir] = true
	}

	// remove deleted files from active paths
	for path := range deletedFiles {
		delete(activePaths, path)
	}

	ignoredPaths := map[string]string{}
	for path := range allPaths {
		if _, ok := activePaths[path]; !ok {
			if ignored != nil && ignored.MatchesPath(path) {
				ignoredPaths[path] = "gpt4cli"
			} else {
				ignoredPaths[path] = "git"
			}
		}
	}

	return &ProjectPaths{
		ActivePaths:    activePaths,
		AllPaths:       allPaths,
		Gpt4cliIgnored: ignored,
		IgnoredPaths:   ignoredPaths,
	}, nil
}

func GetGpt4cliIgnore(dir string) (*ignore.GitIgnore, error) {
	ignorePath := filepath.Join(dir, ".gpt4cliignore")

	if _, err := os.Stat(ignorePath); err == nil {
		ignored, err := ignore.CompileIgnoreFile(ignorePath)

		if err != nil {
			return nil, fmt.Errorf("error reading .gpt4cliignore file: %s", err)
		}

		return ignored, nil
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("error checking for .gpt4cliignore file: %s", err)
	}

	return nil, nil
}

func GetParentProjectIdsWithPaths() ([][2]string, error) {
	var parentProjectIds [][2]string
	currentDir := filepath.Dir(Cwd)

	for currentDir != "/" {
		gpt4cliDir := findGpt4cli(currentDir)
		projectSettingsPath := filepath.Join(gpt4cliDir, "project.json")
		if _, err := os.Stat(projectSettingsPath); err == nil {
			bytes, err := os.ReadFile(projectSettingsPath)
			if err != nil {
				return nil, fmt.Errorf("error reading projectId file: %s", err)
			}

			var settings types.CurrentProjectSettings
			err = json.Unmarshal(bytes, &settings)

			if err != nil {
				term.OutputErrorAndExit("error unmarshalling project.json: %v", err)
			}

			projectId := string(settings.Id)
			parentProjectIds = append(parentProjectIds, [2]string{currentDir, projectId})
		}
		currentDir = filepath.Dir(currentDir)
	}

	return parentProjectIds, nil
}

func GetChildProjectIdsWithPaths(ctx context.Context) ([][2]string, error) {
	var childProjectIds [][2]string

	err := filepath.Walk(Cwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// if permission denied, skip the path
			if os.IsPermission(err) {
				if info.IsDir() {
					return filepath.SkipDir
				} else {
					return nil
				}
			}

			return err
		}

		if strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("context timeout")
		default:
		}

		if info.IsDir() && path != Cwd {
			gpt4cliDir := findGpt4cli(path)
			projectSettingsPath := filepath.Join(gpt4cliDir, "project.json")
			if _, err := os.Stat(projectSettingsPath); err == nil {
				bytes, err := os.ReadFile(projectSettingsPath)
				if err != nil {
					return fmt.Errorf("error reading projectId file: %s", err)
				}
				var settings types.CurrentProjectSettings
				err = json.Unmarshal(bytes, &settings)

				if err != nil {
					term.OutputErrorAndExit("error unmarshalling project.json: %v", err)
				}

				projectId := string(settings.Id)
				childProjectIds = append(childProjectIds, [2]string{path, projectId})
			}
		}
		return nil
	})

	if err != nil {
		if err.Error() == "context timeout" {
			return childProjectIds, nil
		}

		return nil, fmt.Errorf("error walking the path %s: %s", Cwd, err)
	}

	return childProjectIds, nil
}

func GetBaseDirForContexts(contexts []*shared.Context) string {
	var paths []string

	for _, context := range contexts {
		if context.FilePath != "" {
			paths = append(paths, context.FilePath)
		}
	}

	return GetBaseDirForFilePaths(paths)
}

func GetBaseDirForFilePaths(paths []string) string {
	baseDir := ProjectRoot
	dirsUp := 0

	for _, path := range paths {
		currentDir := ProjectRoot

		pathSplit := strings.Split(path, string(os.PathSeparator))

		n := 0
		for _, p := range pathSplit {
			if p == ".." {
				n++
				currentDir = filepath.Dir(currentDir)
			} else {
				break
			}
		}

		if n > dirsUp {
			dirsUp = n
			baseDir = currentDir
		}
	}

	return baseDir
}

func findGpt4cli(baseDir string) string {
	var dir string
	if os.Getenv("GPT4CLI_ENV") == "development" {
		dir = filepath.Join(baseDir, ".gpt4cli-dev")
	} else {
		dir = filepath.Join(baseDir, ".gpt4cli")
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
