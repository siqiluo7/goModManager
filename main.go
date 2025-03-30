package main

import (
	"fmt"

	"github.com/goModManager/modManager"
	"github.com/goModManager/testModule"
)

func main() {
	fmt.Println("Hello World")
	modManager.GetDependencies()
	testModule.Test()
}
