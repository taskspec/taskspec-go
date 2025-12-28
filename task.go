package taskspec

import "time"

// Priority represents the priority level of a task.
type Priority int

const (
	PriorityUnknown Priority = iota
	PriorityHighest          // 1, critical, 🔺
	PriorityHigh             // 2, high, ⏫
	PriorityMedium           // 3, medium, normal, 🔼
	PriorityLow              // 4, low, 🔽
	PriorityLowest           // 5, lowest, ⏬
)

// String returns the string representation of a Priority.
func (p Priority) String() string {
	switch p {
	case PriorityHighest:
		return "highest"
	case PriorityHigh:
		return "high"
	case PriorityMedium:
		return "medium"
	case PriorityLow:
		return "low"
	case PriorityLowest:
		return "lowest"
	default:
		return ""
	}
}

// Status represents the current state of a task.
type Status string

const (
	StatusTodo       Status = "todo"        // ⬜
	StatusInProgress Status = "in-progress" // 🚧
	StatusDone       Status = "done"        // ✅
	StatusCancelled  Status = "cancelled"   // ❌
	StatusBlocked    Status = "blocked"     // 🚫
)

// Task represents a parsed taskspec annotation.
type Task struct {
	// Keyword is the annotation keyword (TODO, FIXME, BUG, etc.)
	// Empty for Markdown task lists.
	Keyword string

	// Description is the task description.
	Description string

	// Metadata fields (all optional)
	DueDate       *time.Time
	ScheduledDate *time.Time
	StartDate     *time.Time
	CreatedDate   *time.Time
	CompletedDate *time.Time
	Priority      Priority
	Recurrence    string
	Identifier    string
	Assignees     []string
	Tags          []string
	Projects      []string
	Status        Status
	Estimate      string

	// CustomFields stores any non-standard metadata fields.
	CustomFields map[string]string

	// IsMarkdownTask indicates if this task was parsed from a Markdown task list.
	IsMarkdownTask bool

	// IsChecked indicates if a Markdown task is checked (only valid if IsMarkdownTask is true).
	IsChecked bool

	// Raw is the original unparsed text.
	Raw string
}
