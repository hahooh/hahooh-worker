package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

const (
	todoPath     = "/shared/todo.json"
	outputPath   = "/shared/output.json"
	workspaceDir = "/app/workspace"
)

func main() {
	fmt.Println("Agent Runner started. Watching", todoPath)
	setupGit()

	for {
		// 1. Read tasks from todo.json
		data, err := os.ReadFile(todoPath)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		var tasks []Task
		json.Unmarshal(data, &tasks)

		if len(tasks) > 0 {
			current := tasks[0]

			// 4. Update todo.json (Remove the completed task)
			remaining := tasks[1:]
			newData, _ := json.MarshalIndent(remaining, "", "  ")
			os.WriteFile(todoPath, newData, 0644)

			fmt.Printf("\n--- Processing Task: %s ---\n", current.ID)

			// 2. Setup Real-time Output + Buffer
			var buf bytes.Buffer
			// MultiWriter sends output to the Docker logs (Stdout) AND our buffer
			multiWriter := io.MultiWriter(os.Stdout, &buf)

			cmd := exec.Command("gemini", "--yolo", "-p", "Task: "+current.Query)
			cmd.Stdout = multiWriter
			cmd.Stderr = multiWriter

			// Start the execution
			err := cmd.Run()
			rawOutput := buf.Bytes()

			success := true
			if err != nil {
				fmt.Printf("\nTask %s failed: %v\n", current.ID, err)
				success = false
			}

			var diffOutput []byte
			if success {
				// Capture the diff using the constant workspaceDir
				diffCmd := exec.Command("git", "-C", workspaceDir, "diff")
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

			fmt.Printf("\n--- Task %s Complete ---\n", current.ID)
		}

		time.Sleep(5 * time.Second)
	}
}

func setupGit() {
	exec.Command("git", "config", "--global", "--add", "safe.directory", "/app").Run()
	exec.Command("git", "config", "--global", "--add", "safe.directory", workspaceDir).Run()

	pat := os.Getenv("GITHUB_PAT")
	if pat != "" {
		authUrl := fmt.Sprintf("https://%s@github.com/", pat)
		exec.Command("git", "config", "--global", "url."+authUrl+".insteadOf", "https://github.com/").Run()
	}
}

func recordResult(path string, res Result) {
	var results []Result
	data, err := os.ReadFile(path)
	if err == nil {
		json.Unmarshal(data, &results)
	}
	results = append(results, res)
	newData, _ := json.MarshalIndent(results, "", "  ")
	os.WriteFile(path, newData, 0644)
}
