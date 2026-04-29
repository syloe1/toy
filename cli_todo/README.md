# CLI Todo

A simple command-line todo application written in Go.

This project stores tasks in a local `tasks.json` file and supports basic CRUD operations from the terminal:

- Add a task
- List tasks
- Complete a task
- Delete a task

## Features

- Store tasks in JSON
- Auto-increment task IDs
- Show only incomplete tasks by default
- Show all tasks with `-a` or `--all`
- Friendly error messages for invalid IDs or broken data files

## Project Structure

```text
cli_todo/
├── internal/
│   ├── model/
│   ├── repository/
│   └── service/
├── main.go
├── go.mod
└── tasks.json
```

## Requirements

- Go 1.25+

## Build

In Windows `cmd`:

```bat
cd /d F:\cli_todo
go build -o task.exe
```

After building, you can run the program with:

```bat
task add "My new task"
```

If you prefer, you can also build it as `tasks.exe`:

```bat
go build -o tasks.exe
```

## Usage

```bat
task add "Tidy my desk"
task list
task list -a
task complete 1
task delete 1
task help
```

## Commands

### Add a task

```bat
task add "Write README"
```

Adds a new task to `tasks.json`.

### List incomplete tasks

```bat
task list
```

Shows only tasks where `done == false`.

### List all tasks

```bat
task list -a
```

Or:

```bat
task list --all
```

Shows all tasks, including completed ones.

### Complete a task

```bat
task complete 1
```

Marks task `1` as completed.

### Delete a task

```bat
task delete 1
```

Removes task `1` from storage.

### Show help

```bat
task help
```

## Data Format

Tasks are stored in `tasks.json` like this:

```json
[
  {
    "id": 1,
    "content": "Write README",
    "done": false,
    "created_at": "2026-04-29T18:40:27+08:00"
  }
]
```

## Example Workflow

```bat
task add "Tidy my desk"
task add "Write documentation"
task list
task complete 1
task list
task list -a
task delete 2
task list -a
```

## Notes

- The app reads and writes `tasks.json` in the current working directory.
- If `tasks.json` is missing, the app will create it when you add the first task.
- If `tasks.json` is invalid JSON, the app will return an error instead of silently overwriting the file.

## Future Improvements

- Better table formatting for terminal output
- Edit/update task content
- Support for custom data file paths
- Sorting and filtering options

