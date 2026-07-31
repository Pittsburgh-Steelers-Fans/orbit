package main

import (
	"fmt"
	"os"
)

type task struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

func runTasks(client *Client, args []string) error {
	if err := parseListCommand("tasks", args); err != nil {
		return err
	}
	var tasks []task
	if err := client.get("/tasks", &tasks); err != nil {
		return err
	}
	for _, task := range tasks {
		fmt.Fprintf(os.Stdout, "%s\t%s\t%s\n", task.ID, task.Title, task.Status)
	}
	return nil
}
