package main

import (
	"fmt"

	"github.com/netsells/katsu/cmd"
)

var GitCommit string = "dev"

func main() {

	fmt.Printf("Hello world, version: %s\n", GitCommit)
	cmd.Execute()
}
