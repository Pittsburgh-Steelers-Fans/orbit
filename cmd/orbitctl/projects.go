package main

import (
	"fmt"
	"os"
)

type project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func runProjects(client *Client, args []string) error {
	if err := parseListCommand("projects", args); err != nil {
		return err
	}
	var projects []project
	if err := client.get("/projects", &projects); err != nil {
		return err
	}
	for _, project := range projects {
		fmt.Fprintf(os.Stdout, "%s\t%s\n", project.ID, project.Name)
	}
	return nil
}
