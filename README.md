# taskspec-go

A Golang library for parsing [taskspec](https://github.com/taskspec/spec) annotations - a universal TODO annotation format.

## Overview

taskspec-go provides a robust parser for extracting structured task information from TODO comments and Markdown task lists. It supports both text-based and emoji-based metadata formats, making it easy to work with annotated tasks in any text-based environment.

## Installation

```bash
go get github.com/taskspec/taskspec-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/taskspec/taskspec-go"
)

func main() {
    parser := taskspec.NewParser()
    
    // Parse a simple TODO
    task, err := parser.Parse("TODO: Fix the login bug due:2026-03-01 @alice #backend p:high")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Description: %s\n", task.Description)
    fmt.Printf("Due Date: %v\n", task.DueDate)
    fmt.Printf("Assignee: %s\n", task.Assignees[0])
    fmt.Printf("Tag: %s\n", task.Tags[0])
    fmt.Printf("Priority: %s\n", task.Priority)
}
```

## Features

- ✅ Parse standard taskspec annotations (`TODO:`, `FIXME:`, etc.)
- ✅ Parse Markdown task lists (`- [ ]`, `- [x]`)
- ✅ Support both text and emoji metadata formats
- ✅ Extract dates (due, scheduled, start, created, completed)
- ✅ Parse priority levels (text, emoji, and numeric)
- ✅ Extract assignees, tags, and projects
- ✅ Parse recurrence patterns
- ✅ Handle identifiers, status, and time estimates
- ✅ Support custom metadata fields
- ✅ Handle escaped characters in descriptions

## Supported Keywords

The parser recognizes the following keywords (case-insensitive):
- `TODO`
- `FIXME`
- `BUG`
- `HACK`
- `NOTE`
- `INFO`
- `IDEA`
- `REFACTOR`
- `REMINDER`

## Metadata Fields

### Dates

All date fields support ISO 8601 format (`YYYY-MM-DD` or with time).

| Field | Text Format | Emoji Format | Example |
|-------|------------|--------------|---------|
| Due Date | `due:YYYY-MM-DD` | `📅 YYYY-MM-DD` | `due:2026-03-01` |
| Scheduled | `scheduled:YYYY-MM-DD` | `⏳ YYYY-MM-DD` | `scheduled:2026-02-15` |
| Start Date | `start:YYYY-MM-DD` | `🛫 YYYY-MM-DD` | `start:2026-02-01` |
| Created | `created:YYYY-MM-DD` | `➕ YYYY-MM-DD` | `created:2026-01-15` |
| Completed | `done:YYYY-MM-DD` | `✅ YYYY-MM-DD` | `done:2026-01-20` |

### Priority

| Level | Text | Emoji | Numeric |
|-------|------|-------|---------|
| Highest | `priority:highest`, `p:critical` | `🔺` | `p:1` |
| High | `priority:high` | `⏫` | `p:2` |
| Medium | `priority:medium`, `p:normal` | `🔼` | `p:3` |
| Low | `priority:low` | `🔽` | `p:4` |
| Lowest | `priority:lowest` | `⏬` | `p:5` |

### Other Metadata

| Field | Format | Example |
|-------|--------|---------|
| Assignee | `@username` or `👤username` | `@alice @bob` |
| Tags | `#tag` | `#backend #urgent` |
| Projects | `+project` | `+ProjectX` |
| Identifier | `id:ID` or `🆔 ID` | `id:TASK-123` |
| Status | `status:value` or emoji | `status:in-progress` or `🚧` |
| Recurrence | `repeat:pattern` or `🔁 pattern` | `repeat:every week` |
| Estimate | `estimate:duration` or `⏱️ duration` | `estimate:2h` |

### Status Values

| Status | Text | Emoji |
|--------|------|-------|
| Todo | `status:todo` | `⬜` |
| In Progress | `status:in-progress` | `🚧` |
| Done | `status:done` | `✅` |
| Cancelled | `status:cancelled` | `❌` |
| Blocked | `status:blocked` | `🚫` |

## Examples

### Standard Format

```go
parser := taskspec.NewParser()

// Simple TODO
task, _ := parser.Parse("TODO: Fix the bug")
// task.Keyword = "TODO"
// task.Description = "Fix the bug"

// TODO with metadata
task, _ = parser.Parse("TODO: Implement feature due:2026-03-01 @alice #backend p:high")
// task.Description = "Implement feature"
// task.DueDate = time for 2026-03-01
// task.Assignees = ["alice"]
// task.Tags = ["backend"]
// task.Priority = taskspec.PriorityHigh

// Using emoji format
task, _ = parser.Parse("FIXME: Refactor code 📅 2026-03-15 🔺 @bob")
// task.Description = "Refactor code"
// task.DueDate = time for 2026-03-15
// task.Priority = taskspec.PriorityHighest
// task.Assignees = ["bob"]
```

### Markdown Task Lists

```go
// Unchecked task
task, _ := parser.Parse("- [ ] Buy groceries due:2026-03-01")
// task.IsMarkdownTask = true
// task.IsChecked = false
// task.Description = "Buy groceries"
// task.DueDate = time for 2026-03-01

// Checked task
task, _ = parser.Parse("- [x] Complete report 📅 2026-02-28")
// task.IsMarkdownTask = true
// task.IsChecked = true
// task.Description = "Complete report"
// task.DueDate = time for 2026-02-28
```

### Multiple Metadata Fields

```go
task, _ := parser.Parse(`TODO: Deploy to production 
    due:2026-03-01 
    scheduled:2026-02-28 
    @alice @bob 
    #deployment #critical 
    +ProjectX 
    id:DEPLOY-123 
    estimate:4h 
    status:in-progress`)

// All metadata fields are populated
```

### Custom Fields

```go
task, _ := parser.Parse("TODO: Custom task env:production region:us-east-1")
// task.CustomFields["env"] = "production"
// task.CustomFields["region"] = "us-east-1"
```

### Escaping

```go
task, _ := parser.Parse(`TODO: This is about the \#backend team`)
// task.Description = "This is about the #backend team"
// task.Tags = [] (empty, # was escaped)
```

## API Reference

### Parser

```go
type Parser struct {}

// NewParser creates a new Parser instance
func NewParser() *Parser

// Parse parses a single taskspec annotation
func (p *Parser) Parse(text string) (*Task, error)

// ParseLines parses multiple lines and returns all found tasks
func (p *Parser) ParseLines(lines []string) ([]*Task, error)
```

### Task

```go
type Task struct {
    Keyword        string            // "TODO", "FIXME", etc. (empty for Markdown tasks)
    Description    string            // Task description
    DueDate        *time.Time        // When task is due
    ScheduledDate  *time.Time        // When task is scheduled
    StartDate      *time.Time        // When to start task
    CreatedDate    *time.Time        // When task was created
    CompletedDate  *time.Time        // When task was completed
    Priority       Priority          // Task priority level
    Recurrence     string            // Recurrence pattern
    Identifier     string            // External identifier
    Assignees      []string          // List of assignees
    Tags           []string          // List of tags
    Projects       []string          // List of projects
    Status         Status            // Current status
    Estimate       string            // Time estimate
    CustomFields   map[string]string // Custom metadata
    IsMarkdownTask bool              // True if parsed from Markdown
    IsChecked      bool              // True if Markdown task is checked
    Raw            string            // Original text
}
```

### Priority

```go
type Priority int

const (
    PriorityUnknown Priority = iota
    PriorityHighest
    PriorityHigh
    PriorityMedium
    PriorityLow
    PriorityLowest
)

func (p Priority) String() string
```

### Status

```go
type Status string

const (
    StatusTodo       Status = "todo"
    StatusInProgress Status = "in-progress"
    StatusDone       Status = "done"
    StatusCancelled  Status = "cancelled"
    StatusBlocked    Status = "blocked"
)
```

## Specification

This library implements the [taskspec specification](https://raw.githubusercontent.com/taskspec/spec/refs/heads/main/SPEC.md).

## Testing

This project includes comprehensive testing:
- **Unit tests** - Standard Go tests for all parser functionality
- **Fuzz tests** - Automated fuzzing to discover edge cases and crashes
- **Grammar-based tests** - Tests against samples generated from the ABNF grammar

See [TESTING.md](TESTING.md) for detailed information about running tests and the CI/CD pipeline.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

See the [LICENSE](LICENSE) file for details.