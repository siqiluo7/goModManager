package modManager

import (
	"bufio"
	"bytes"
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
	Name          string
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
			dependencies = append(dependencies, Dependency{moduleName: moduleName, Name: matches[0], path: matches[1], PseudoVersion: matches[2]})
		}

	}
	return dependencies

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
