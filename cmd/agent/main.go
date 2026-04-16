package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type Task struct {
	ID     string `json:"id"`
	Query  string `json:"query"`
	Status string `json:"status"`
}

type Result struct {
	ID        string    `json:"id"`
	Query     string    `json:"query"`
	Output    string    `json:"output"`
	Diff      string    `json:"diff"`
	Timestamp time.Time `json:"timestamp"`
	Success   bool      `json:"success"`
}

func main() {
	todoPath := os.Getenv("TODO_PATH")
	outputPath := os.Getenv("OUTPUT_PATH")

	fmt.Println("Agent Runner started. Watching", todoPath)

	for {
		// 1. Read tasks from todo.json
		data, _ := os.ReadFile(todoPath)
		var tasks []Task
		json.Unmarshal(data, &tasks)

		if len(tasks) > 0 {
			current := tasks[0]
			fmt.Printf("Processing Task: %s\n", current.ID)

			// 2. Execute Gemini CLI and capture output
			// Using CombinedOutput to get both Stdout and Stderr for the log
			fmt.Printf("Running: gemini --yolo -p %s\n", current.Query)
			cmd := exec.Command("gemini", "--yolo", "-p", current.Query)
			rawOutput, err := cmd.CombinedOutput()

			success := true
			if err != nil {
				fmt.Printf("Task %s failed: %v\n", current.ID, err)
				success = false
			}

			var diffOutput []byte
			if success {
				// 1. Capture the diff
				diffCmd := exec.Command("git", "-C", os.Getenv("WORKSPACE_DIR"), "diff")
				diffOutput, _ = diffCmd.Output()
			}

			// 3. Record to output.json
			recordResult(outputPath, Result{
				ID:        current.ID,
				Query:     current.Query,
				Output:    string(rawOutput),
				Diff:      string(diffOutput),
				Timestamp: time.Now(),
				Success:   success,
			})

			// 4. Update todo.json (Remove the completed task)
			remaining := tasks[1:]
			newData, _ := json.MarshalIndent(remaining, "", "  ")
			os.WriteFile(todoPath, newData, 0644)

			fmt.Printf("Task %s complete. Results written to %s\n", current.ID, outputPath)
		}

		time.Sleep(5 * time.Second)
	}
}

func recordResult(path string, res Result) {
	var results []Result

	// Read existing results
	data, err := os.ReadFile(path)
	if err == nil {
		json.Unmarshal(data, &results)
	}

	// Append new result
	results = append(results, res)

	// Write back to file
	newData, _ := json.MarshalIndent(results, "", "  ")
	os.WriteFile(path, newData, 0644)
}
