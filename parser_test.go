package taskspec

import (
	"testing"
	"time"
)

func TestParseStandardTask(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name        string
		input       string
		wantKeyword string
		wantDesc    string
		wantNil     bool
	}{
		{
			name:        "Simple TODO",
			input:       "TODO: Fix the bug",
			wantKeyword: "TODO",
			wantDesc:    "Fix the bug",
		},
		{
			name:        "FIXME with description",
			input:       "FIXME: Refactor this code",
			wantKeyword: "FIXME",
			wantDesc:    "Refactor this code",
		},
		{
			name:        "Case insensitive",
			input:       "todo: lowercase keyword",
			wantKeyword: "TODO",
			wantDesc:    "lowercase keyword",
		},
		{
			name:    "Not a task",
			input:   "This is just a comment",
			wantNil: true,
		},
		{
			name:    "Empty string",
			input:   "",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if tt.wantNil {
				if task != nil {
					t.Errorf("Parse() = %v, want nil", task)
				}
				return
			}

			if task == nil {
				t.Fatal("Parse() = nil, want task")
			}

			if task.Keyword != tt.wantKeyword {
				t.Errorf("Keyword = %v, want %v", task.Keyword, tt.wantKeyword)
			}

			if task.Description != tt.wantDesc {
				t.Errorf("Description = %v, want %v", task.Description, tt.wantDesc)
			}
		})
	}
}

func TestParseMarkdownTask(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name        string
		input       string
		wantDesc    string
		wantChecked bool
	}{
		{
			name:        "Unchecked task",
			input:       "- [ ] Buy groceries",
			wantDesc:    "Buy groceries",
			wantChecked: false,
		},
		{
			name:        "Checked task lowercase",
			input:       "- [x] Complete the report",
			wantDesc:    "Complete the report",
			wantChecked: true,
		},
		{
			name:        "Checked task uppercase",
			input:       "- [X] Call the client",
			wantDesc:    "Call the client",
			wantChecked: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if task == nil {
				t.Fatal("Parse() = nil, want task")
			}

			if !task.IsMarkdownTask {
				t.Error("IsMarkdownTask = false, want true")
			}

			if task.IsChecked != tt.wantChecked {
				t.Errorf("IsChecked = %v, want %v", task.IsChecked, tt.wantChecked)
			}

			if task.Description != tt.wantDesc {
				t.Errorf("Description = %v, want %v", task.Description, tt.wantDesc)
			}

			if task.Keyword != "" {
				t.Errorf("Keyword = %v, want empty", task.Keyword)
			}
		})
	}
}

func TestParseDates(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name     string
		input    string
		wantDate string
		checkDue bool
	}{
		{
			name:     "Due date text format",
			input:    "TODO: Fix bug due:2026-02-15",
			wantDate: "2026-02-15",
			checkDue: true,
		},
		{
			name:     "Due date emoji format",
			input:    "TODO: Fix bug 📅 2026-02-15",
			wantDate: "2026-02-15",
			checkDue: true,
		},
		{
			name:     "Due date with time",
			input:    "TODO: Fix bug due:2026-02-15T18:30:00Z",
			wantDate: "2026-02-15T18:30:00Z",
			checkDue: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if task == nil {
				t.Fatal("Parse() = nil, want task")
			}

			if tt.checkDue {
				if task.DueDate == nil {
					t.Fatal("DueDate = nil, want date")
				}

				// Parse expected date
				var expectedTime time.Time
				formats := []string{"2006-01-02", time.RFC3339}
				for _, format := range formats {
					if parsed, err := time.Parse(format, tt.wantDate); err == nil {
						expectedTime = parsed
						break
					}
				}

				if !task.DueDate.Equal(expectedTime) {
					t.Errorf("DueDate = %v, want %v", task.DueDate, expectedTime)
				}
			}
		})
	}
}

func TestParsePriority(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name         string
		input        string
		wantPriority Priority
	}{
		{
			name:         "Priority text highest",
			input:        "TODO: Fix bug priority:highest",
			wantPriority: PriorityHighest,
		},
		{
			name:         "Priority text critical",
			input:        "TODO: Fix bug priority:critical",
			wantPriority: PriorityHighest,
		},
		{
			name:         "Priority short form",
			input:        "TODO: Fix bug p:high",
			wantPriority: PriorityHigh,
		},
		{
			name:         "Priority numeric",
			input:        "TODO: Fix bug p:3",
			wantPriority: PriorityMedium,
		},
		{
			name:         "Priority emoji highest",
			input:        "TODO: Fix bug 🔺",
			wantPriority: PriorityHighest,
		},
		{
			name:         "Priority emoji high",
			input:        "TODO: Fix bug ⏫",
			wantPriority: PriorityHigh,
		},
		{
			name:         "Priority emoji high (up arrow)",
			input:        "TODO: Fix bug 🔼",
			wantPriority: PriorityHigh,
		},
		{
			name:         "Priority emoji low",
			input:        "TODO: Fix bug 🔽",
			wantPriority: PriorityLow,
		},
		{
			name:         "Priority emoji lowest",
			input:        "TODO: Fix bug ⏬",
			wantPriority: PriorityLowest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if task == nil {
				t.Fatal("Parse() = nil, want task")
			}

			if task.Priority != tt.wantPriority {
				t.Errorf("Priority = %v, want %v", task.Priority, tt.wantPriority)
			}
		})
	}
}

func TestParseAssignees(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name          string
		input         string
		wantAssignees []string
	}{
		{
			name:          "Single assignee",
			input:         "TODO: Fix bug @martin",
			wantAssignees: []string{"martin"},
		},
		{
			name:          "Multiple assignees",
			input:         "TODO: Fix bug @alice @bob",
			wantAssignees: []string{"alice", "bob"},
		},
		{
			name:          "Emoji assignee",
			input:         "TODO: Fix bug 👤john",
			wantAssignees: []string{"john"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if task == nil {
				t.Fatal("Parse() = nil, want task")
			}

			if len(task.Assignees) != len(tt.wantAssignees) {
				t.Errorf("Assignees count = %v, want %v", len(task.Assignees), len(tt.wantAssignees))
			}

			for i, assignee := range tt.wantAssignees {
				if i >= len(task.Assignees) || task.Assignees[i] != assignee {
					t.Errorf("Assignees[%d] = %v, want %v", i, task.Assignees, assignee)
				}
			}
		})
	}
}

func TestParseTagsAndProjects(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name         string
		input        string
		wantTags     []string
		wantProjects []string
	}{
		{
			name:     "Single tag",
			input:    "TODO: Fix bug #backend",
			wantTags: []string{"backend"},
		},
		{
			name:     "Multiple tags",
			input:    "TODO: Fix bug #backend #urgent",
			wantTags: []string{"backend", "urgent"},
		},
		{
			name:         "Single project",
			input:        "TODO: Fix bug +ProjectX",
			wantProjects: []string{"ProjectX"},
		},
		{
			name:         "Tags and projects",
			input:        "TODO: Fix bug #backend +ProjectX",
			wantTags:     []string{"backend"},
			wantProjects: []string{"ProjectX"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if task == nil {
				t.Fatal("Parse() = nil, want task")
			}

			if len(task.Tags) != len(tt.wantTags) {
				t.Errorf("Tags count = %v, want %v", len(task.Tags), len(tt.wantTags))
			}

			for i, tag := range tt.wantTags {
				if i >= len(task.Tags) || task.Tags[i] != tag {
					t.Errorf("Tags[%d] = %v, want %v", i, task.Tags, tag)
				}
			}

			if len(task.Projects) != len(tt.wantProjects) {
				t.Errorf("Projects count = %v, want %v", len(task.Projects), len(tt.wantProjects))
			}

			for i, project := range tt.wantProjects {
				if i >= len(task.Projects) || task.Projects[i] != project {
					t.Errorf("Projects[%d] = %v, want %v", i, task.Projects, project)
				}
			}
		})
	}
}

func TestParseStatus(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name       string
		input      string
		wantStatus Status
	}{
		{
			name:       "Status todo",
			input:      "TODO: Fix bug status:todo",
			wantStatus: StatusTodo,
		},
		{
			name:       "Status in-progress",
			input:      "TODO: Fix bug status:in-progress",
			wantStatus: StatusInProgress,
		},
		{
			name:       "Status done",
			input:      "TODO: Fix bug status:done",
			wantStatus: StatusDone,
		},
		{
			name:       "Status emoji in-progress",
			input:      "TODO: Fix bug 🚧",
			wantStatus: StatusInProgress,
		},
		{
			name:       "Status emoji blocked",
			input:      "TODO: Fix bug 🚫",
			wantStatus: StatusBlocked,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if task == nil {
				t.Fatal("Parse() = nil, want task")
			}

			if task.Status != tt.wantStatus {
				t.Errorf("Status = %v, want %v", task.Status, tt.wantStatus)
			}
		})
	}
}

func TestParseRecurrence(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name           string
		input          string
		wantRecurrence string
	}{
		{
			name:           "Recurrence text",
			input:          "TODO: Daily standup repeat:every day",
			wantRecurrence: "every day",
		},
		{
			name:           "Recurrence rec shorthand",
			input:          "TODO: Weekly review rec:every week",
			wantRecurrence: "every week",
		},
		{
			name:           "Recurrence emoji",
			input:          "TODO: Monthly report 🔁 every month",
			wantRecurrence: "every month",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if task == nil {
				t.Fatal("Parse() = nil, want task")
			}

			if task.Recurrence != tt.wantRecurrence {
				t.Errorf("Recurrence = %v, want %v", task.Recurrence, tt.wantRecurrence)
			}
		})
	}
}

func TestParseIdentifier(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name           string
		input          string
		wantIdentifier string
	}{
		{
			name:           "ID text format",
			input:          "TODO: Fix bug id:TODO-1234",
			wantIdentifier: "TODO-1234",
		},
		{
			name:           "ID emoji format",
			input:          "TODO: Fix bug 🆔 ISSUE-567",
			wantIdentifier: "ISSUE-567",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if task == nil {
				t.Fatal("Parse() = nil, want task")
			}

			if task.Identifier != tt.wantIdentifier {
				t.Errorf("Identifier = %v, want %v", task.Identifier, tt.wantIdentifier)
			}
		})
	}
}

func TestParseEstimate(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name         string
		input        string
		wantEstimate string
	}{
		{
			name:         "Estimate text format",
			input:        "TODO: Fix bug estimate:2h",
			wantEstimate: "2h",
		},
		{
			name:         "Estimate emoji format",
			input:        "TODO: Fix bug ⏱️ 3d",
			wantEstimate: "3d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if task == nil {
				t.Fatal("Parse() = nil, want task")
			}

			if task.Estimate != tt.wantEstimate {
				t.Errorf("Estimate = %v, want %v", task.Estimate, tt.wantEstimate)
			}
		})
	}
}

func TestParseComplexTask(t *testing.T) {
	parser := NewParser()

	input := "TODO: Implement new feature due:2026-03-01 @alice @bob #backend +ProjectX p:high"
	task, err := parser.Parse(input)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if task == nil {
		t.Fatal("Parse() = nil, want task")
	}

	if task.Keyword != "TODO" {
		t.Errorf("Keyword = %v, want TODO", task.Keyword)
	}

	if task.Description != "Implement new feature" {
		t.Errorf("Description = %v, want 'Implement new feature'", task.Description)
	}

	if task.DueDate == nil {
		t.Fatal("DueDate = nil, want date")
	}

	expectedDate, _ := time.Parse("2006-01-02", "2026-03-01")
	if !task.DueDate.Equal(expectedDate) {
		t.Errorf("DueDate = %v, want %v", task.DueDate, expectedDate)
	}

	if len(task.Assignees) != 2 || task.Assignees[0] != "alice" || task.Assignees[1] != "bob" {
		t.Errorf("Assignees = %v, want [alice bob]", task.Assignees)
	}

	if len(task.Tags) != 1 || task.Tags[0] != "backend" {
		t.Errorf("Tags = %v, want [backend]", task.Tags)
	}

	if len(task.Projects) != 1 || task.Projects[0] != "ProjectX" {
		t.Errorf("Projects = %v, want [ProjectX]", task.Projects)
	}

	if task.Priority != PriorityHigh {
		t.Errorf("Priority = %v, want PriorityHigh", task.Priority)
	}
}

func TestParseEscaping(t *testing.T) {
	parser := NewParser()

	input := `TODO: This is about the \#backend team`
	task, err := parser.Parse(input)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if task == nil {
		t.Fatal("Parse() = nil, want task")
	}

	// The escaped # should be part of the description, not parsed as a tag
	if len(task.Tags) != 0 {
		t.Errorf("Tags = %v, want empty", task.Tags)
	}
}

func TestParseCustomFields(t *testing.T) {
	parser := NewParser()

	input := "TODO: Custom task custom-field:value123"
	task, err := parser.Parse(input)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if task == nil {
		t.Fatal("Parse() = nil, want task")
	}

	if task.CustomFields == nil {
		t.Fatal("CustomFields = nil, want map")
	}

	if val, ok := task.CustomFields["custom-field"]; !ok || val != "value123" {
		t.Errorf("CustomFields[custom-field] = %v, want value123", val)
	}
}

func TestPriorityString(t *testing.T) {
	tests := []struct {
		priority Priority
		want     string
	}{
		{PriorityHighest, "highest"},
		{PriorityHigh, "high"},
		{PriorityMedium, "medium"},
		{PriorityLow, "low"},
		{PriorityLowest, "lowest"},
		{PriorityUnknown, ""},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.priority.String(); got != tt.want {
				t.Errorf("Priority.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
