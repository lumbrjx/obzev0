package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func printHelp() {
	fmt.Println(`Usage: obzev0 <command> [options]

Commands:
  init    Initialize the service
  apply   Apply the configuration

Run 'obzev0 help <command>' for more information on a specific command.`)
}
func readArgs() {
	if len(os.Args) < 2 {
		fmt.Println("expected 'init' or 'apply' command")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "help":
		printHelp()
	case "init":
		initCmd := flag.NewFlagSet("init", flag.ExitOnError)
		dist := initCmd.String(
			"dist",
			"",
			"Specify the distination directory (default is current directory)",
		)
		addr := initCmd.String(
			"Addr",
			"127.0.0.1:50051",
			"Specify the gRPC server address (default is 127.0.0.1:50051)",
		)

		initCmd.Parse(os.Args[2:])

		fmt.Printf("Initializing with dist=%s and Addr=%s\n", *dist, *addr)

		generateYaml(*addr, *dist)

	case "apply":
		applyCmd := flag.NewFlagSet("apply", flag.ExitOnError)
		path := applyCmd.String(
			"path",
			"",
			"Specify the path to apply (default is current directory)",
		)

		applyCmd.Parse(os.Args[2:])

		fmt.Printf("Applying configurations from path=%s\n", *path)
		c, err := LoadConfig(*path)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		apply(c)

	default:
		fmt.Println("expected 'init', 'apply' or 'help' command")
		os.Exit(1)
	}
}
