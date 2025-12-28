package taskspec

import (
	"testing"
)

func BenchmarkParseSimpleTODO(b *testing.B) {
	parser := NewParser()
	input := "TODO: Fix the bug"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.Parse(input)
	}
}

func BenchmarkParseComplexTask(b *testing.B) {
	parser := NewParser()
	input := "TODO: Implement a new feature due:2026-03-01 @alice @bob #backend +ProjectX p:high"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.Parse(input)
	}
}

func BenchmarkParseWithAllMetadata(b *testing.B) {
	parser := NewParser()
	input := "TODO: Deploy to production due:2026-03-01 scheduled:2026-02-28 @alice @bob #deployment #critical +ProjectX id:DEPLOY-123 estimate:4h status:in-progress"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.Parse(input)
	}
}

func BenchmarkParseMarkdownTask(b *testing.B) {
	parser := NewParser()
	input := "- [ ] Buy groceries due:2026-03-01"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.Parse(input)
	}
}

func BenchmarkParseWithEmojis(b *testing.B) {
	parser := NewParser()
	input := "FIXME: Refactor code 📅 2026-03-15 🔺 👤bob"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.Parse(input)
	}
}

func BenchmarkParseLines(b *testing.B) {
	parser := NewParser()
	lines := []string{
		"TODO: Fix the bug",
		"FIXME: Refactor this code",
		"- [ ] Buy groceries",
		"TODO: Implement feature due:2026-03-01 @alice #backend",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.ParseLines(lines)
	}
}

func BenchmarkParseDates(b *testing.B) {
	parser := NewParser()
	input := "TODO: Fix bug due:2026-02-15 scheduled:2026-02-10 start:2026-02-01"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.Parse(input)
	}
}

func BenchmarkParsePriority(b *testing.B) {
	parser := NewParser()
	input := "TODO: Fix bug priority:highest p:high p:3"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.Parse(input)
	}
}

func BenchmarkParseAssignees(b *testing.B) {
	parser := NewParser()
	input := "TODO: Fix bug @alice @bob @charlie @diana"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.Parse(input)
	}
}

func BenchmarkParseTagsAndProjects(b *testing.B) {
	parser := NewParser()
	input := "TODO: Fix bug #backend #urgent #critical +ProjectX +ProjectY"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.Parse(input)
	}
}
