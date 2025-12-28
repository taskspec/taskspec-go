package taskspec

import (
	"testing"
)

// FuzzParser tests the parser with random inputs to find panics and crashes
func FuzzParser(f *testing.F) {
	parser := NewParser()

	// Seed corpus with valid examples
	seeds := []string{
		"TODO: Fix the bug",
		"FIXME: Refactor this code",
		"- [ ] Buy groceries",
		"- [x] Complete report",
		"TODO: Implement feature due:2026-03-01 @alice #backend p:high",
		"TODO: Fix bug 📅 2026-02-15 🔺 @bob",
		"TODO: Deploy to production due:2026-03-01 scheduled:2026-02-28 @alice @bob #deployment #critical +ProjectX id:DEPLOY-123 estimate:4h status:in-progress",
		"TODO: Custom task custom-field:value123",
		`TODO: This is about the \#backend team`,
		"",
		"Not a task",
		"TODO:",
		"TODO: ",
		"TODO: 📅",
		"TODO: due:",
		"TODO: @",
		"TODO: #",
		"TODO: +",
		"- [ ]",
		"- [x]",
		"BUG: Critical issue p:1",
		"HACK: Temporary fix rec:every day",
		"NOTE: Important 🆔 NOTE-456",
		"INFO: Status update status:done",
		"IDEA: New feature ⏱️ 2w",
		"REFACTOR: Code cleanup start:2026-01-01",
		"REMINDER: Meeting 🔁 every week",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		// The parser should not panic for any input
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Parser panicked on input %q: %v", input, r)
			}
		}()

		// Try to parse the input
		task, err := parser.Parse(input)

		// We allow errors, but we should not panic
		_ = task
		_ = err
	})
}

// FuzzParseLines tests the ParseLines function with random inputs
func FuzzParseLines(f *testing.F) {
	parser := NewParser()

	// Seed corpus with valid examples
	seeds := [][]string{
		{"TODO: First task", "FIXME: Second task"},
		{"- [ ] Task one", "- [x] Task two"},
		{"TODO: Task with metadata due:2026-01-01", "FIXME: Another task @alice"},
		{"", "TODO: Task after empty line"},
		{"Not a task", "TODO: Actual task"},
		{},
	}

	// Convert [][]string to []byte for fuzzing
	for _, seed := range seeds {
		combined := ""
		for i, line := range seed {
			if i > 0 {
				combined += "\n"
			}
			combined += line
		}
		f.Add(combined)
	}

	f.Fuzz(func(t *testing.T, input string) {
		// The parser should not panic for any input
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("ParseLines panicked on input %q: %v", input, r)
			}
		}()

		// Split input into lines and parse
		lines := []string{}
		if input != "" {
			// Simple line splitting
			currentLine := ""
			for _, ch := range input {
				if ch == '\n' {
					lines = append(lines, currentLine)
					currentLine = ""
				} else {
					currentLine += string(ch)
				}
			}
			if currentLine != "" || len(lines) > 0 {
				lines = append(lines, currentLine)
			}
		}

		// Try to parse the lines
		tasks, err := parser.ParseLines(lines)

		// We allow errors, but we should not panic
		_ = tasks
		_ = err
	})
}
