package main

import (
	"fmt"

	"github.com/andydunstall/minc/pkg/cli"
)

func main() {
	if err := cli.Start(); err != nil {
		fmt.Println(err)
	}
}
