package main

import (
	"fmt"

	"github.com/siqiluo7/goModManager/modManager"
)

func main() {
	path := modManager.GetAllModFiles()
	for _, p := range path {
		deps := modManager.GetDependenciesFromModFile(p)
		for _, dep := range deps {
			fmt.Println("dep", dep.Name, dep.Version)

		}
	}

}
