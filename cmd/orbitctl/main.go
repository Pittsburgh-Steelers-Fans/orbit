package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: orbitctl <projects|tasks> <command>")
	}
	client, err := NewClient()
	if err != nil {
		return err
	}
	switch args[0] {
	case "projects":
		return runProjects(client, args[1:])
	case "tasks":
		return runTasks(client, args[1:])
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func parseListCommand(name string, args []string) error {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 || flags.Arg(0) != "list" {
		return fmt.Errorf("usage: orbitctl %s list", name)
	}
	return nil
}
