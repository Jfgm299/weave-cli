package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Jfgm299/weave-cli/internal/cli"
)

func main() {
	code, err := cli.Run(context.Background(), os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	os.Exit(code)
}
