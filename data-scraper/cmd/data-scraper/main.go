package main

import "github.com/kaniuse/kaniuse/data-scraper/pkg/cmds"

func main() {
	rootCommand, err := cmds.NewRootCommand()
	if err != nil {
		panic(err)
	}
	if err := rootCommand.Execute(); err != nil {
		panic(err)
	}
}
