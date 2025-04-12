package modManager

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	test "github.com/siqiluo7/goModManager/testModule"
)

type Dependency struct {
	path          string
	Version       string
	PseudoVersion string
	URL           string
	moduleName    string
}

// Regex to match pseudo versions (v0.0.0-<YYYYMMDDHHMMSS>-<commit-hash>)
var pseudoVersionRegex = regexp.MustCompile(`([a-zA-Z0-9./_-]+) (v0.0.0-\d{14}-[a-f0-9]+)`)

func GetAllModFiles() []string {
	test.Test()
	modPath := []string{}
	rootPath := getRepoRoot()
	filepath.Walk(rootPath+"/", func(path string, info fs.FileInfo, err error) error {

		if err != nil {
			return err
		}
		if info.Name() == "go.mod" {
			modPath = append(modPath, path)
		}
		return nil
	})
	return modPath
}

func GetDependenciesFromModFile(modPath string) []Dependency {
	file, err := os.Open(modPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var moduleName string
	var dependencies []Dependency

	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "module") {
			moduleName = strings.Split(text, " ")[1]
		}

		matches := pseudoVersionRegex.FindStringSubmatch(text)
		if matches != nil {
			dependencies = append(dependencies, Dependency{moduleName: moduleName, URL: matches[0], path: matches[1], PseudoVersion: matches[2]})
		}

	}
	return dependencies

}

func CheckIfPseudoVersionValid(dep Dependency) (bool, error) {
	harshCommit := exactPseudoVersion(dep)
	repoUrl := extractRepoURL(dep)
	valid, err := checkIfCommitExists(repoUrl, harshCommit)
	if err != nil {
		return false, err
	}
	if !valid {
		return false, nil
	}

	isCommitMerged, err := isCommitMerged(repoUrl, harshCommit)
	if err != nil {
		return false, err
	}
	return !isCommitMerged, nil
}

func checkIfCommitExists(repoUrl string, commitHash string) (bool, error) {
	cmd := exec.Command("git", "ls-remote", "https://"+repoUrl)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}

	res := strings.TrimSpace(string(output))
	for _, line := range strings.Split(res, "\n") {
		if strings.Contains(line, commitHash) {
			fmt.Println("Valid pseudo version for dependency")
			return true, nil
		}
	}
	return false, nil
}

func exactPseudoVersion(dep Dependency) string {
	regex := regexp.MustCompile(`v0.0.0-(\d{14})-([a-f0-9]+)`)
	matches := regex.FindStringSubmatch(dep.PseudoVersion)
	if matches == nil {
		fmt.Println("Invalid pseudo version for dependency", dep.URL)
	}
	// fmt.Println("checking version....", dep.URL, matches)
	return matches[2]
}

func extractRepoURL(dep Dependency) string {
	// Define regex pattern to capture "github.com/user/repo"
	re := regexp.MustCompile(`^(github\.com/[^/]+/[^/]+)`)

	matches := re.FindStringSubmatch(dep.URL)
	// fmt.Println("matches.......", matches)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func getRepoRoot() string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(out.String())
}

func isCommitMerged(repoURL, commitHash string) (bool, error) {
	tmpDir := "/tmp/go-mod-check"
	exec.Command("rm", "-rf", tmpDir).Run() // Clean up any previous temp dir
	err := exec.Command("git", "clone", "--bare", repoURL, tmpDir).Run()
	if err != nil {
		return false, fmt.Errorf("failed to clone repository: %v", err)
	}

	cmd := exec.Command("git", "--git-dir="+tmpDir, "merge-base", "--is-ancestor", commitHash, "origin/main")
	err = cmd.Run()

	exec.Command("rm", "-rf", tmpDir).Run()

	// If exit code is 0, commit is merged; if non-zero, it is not
	if err == nil {
		return true, nil
	}
	return false, nil
}
