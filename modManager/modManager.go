package modManager

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	test "github.com/siqiluo7/goModManager/testModule"
)

type dependency struct {
	path    string
	version string
}

func GetDependencies() []dependency {
	test.Test()
	cmd := exec.Command("go", "list", "-m", "-json", "all")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error", err, out)
		panic(err)
	}
	decoder := json.NewDecoder(strings.NewReader(string(out)))
	fmt.Println("mod", string(out))
	regex, err := regexp.Compile(`v\d+\.\d+\.\d+(-[a-zA-Z0-9.]+)?-\d{14}-[a-f0-9]+`)
	if err != nil {
		panic(err)
	}
	var deps []dependency
	for decoder.More() {
		var mod map[string]interface{}
		if err := decoder.Decode(&mod); err != nil {
			panic(err)
		}
		fmt.Println("Path", mod)
		if version, ok := mod["Version"]; ok {
			fmt.Println("matched", version, regex.MatchString(version.(string)))
			deps = append(deps, dependency{path: mod["Path"].(string), version: version.(string)})
		}
	}
	fmt.Println("deps", deps)
	return deps
}
