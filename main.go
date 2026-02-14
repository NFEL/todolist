package main

import (
	"context"
	"flag"
	"graph-interview/cmd"

	"github.com/charmbracelet/log"
)

func main() {
	flag.Parse()
	var err error
	ctx := context.Background()

	switch flag.Arg(0) {
	case "api":
		err = cmd.ApiCmd(ctx, flag.Args()[1:])
	default:
		log.Error("expected 'api' subcommand")
		return
	}
	if err != nil {
		log.Fatal("runtime error happened", "err", err)
	}
}
