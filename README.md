# hahooh-worker 🤖📦

### The Headless Coding Agent that lives in your file system.

hahooh-worker is a minimalist, Unix-style background agent designed for autonomous, unattended coding tasks. Built on the "Batch-JSON" philosophy, it eliminates the need for complex state management or fragile HTTP connections. You simply drop a task into a JSON file, and the worker grinds through it in a secure, sandboxed environment.

## 📜 The Philosophy

- Local-First: Your code, your context, and your secrets stay on your machine.

- Resilient by Design: By using the filesystem as a message bus, the system is immortal. If the container crashes, the state is preserved on your host.

- Minimalist: No heavy frameworks. Just a Go binary, a Docker container, and the Gemini CLI.

- Unattended Execution: Designed for deep work. Feed it a list of 20 refactors, turn off your monitor, and return to a completed output.json with full git diffs.

## 🛠 How it Works

- The Inbox: You append tasks to todo.json.

- The Worker: A Go-based consumer picks up the task, prepares the workspace, and invokes the Gemini CLI.

- The Sandbox: All execution happens inside a Docker container, protecting your host (Arch Linux/macOS/Ubuntu) from unintended side effects.

- The Receipt: Results, logs, and exact code changes are written to output.json for your review.

## 🏗️ Tech Stack

- Engine: Go 1.26 (Optimized for cgo performance)

- Brain: Gemini 2.5 Pro (via Gemini CLI)

- Isolation: Docker + Docker Compose

- Context: Markdown-driven memory via .gemini/GEMINI.md

## ⚠️ Security Warning Note
- This agent is designed for local use. Never give it a PAT with more permissions than it needs, and always review the output.json before merging changes to your main branch.

## 🚀 Quick Start

```bash
# 1. Set your keys in .env

echo "GEMINI_API_KEY=your_key\nGITHUB_PAT=your_pat" > .env

# 2. Add a task

echo '[{"id": "fix_01", "query": "Add unit tests to main.go", "status": "pending"}]' > shared/todo.json

# 3. Spin up the worker

docker-compose up --build
```
