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
	regex := regexp.MustCompile(`v0.0.0-\d{14}-[a-f0-9]+`)
	matches := regex.FindStringSubmatch(dep.PseudoVersion)
	if matches == nil {
		fmt.Println("Invalid pseudo version for dependency", dep.URL)
	}
	harshCommit := matches[1]
	cmd := exec.Command("git", "ls-remote", "https://"+dep.URL, harshCommit)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	cmd.search.NewMatch(arg)
	if err != nil {
		return false, err
	}

	return out.String(), nil
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
