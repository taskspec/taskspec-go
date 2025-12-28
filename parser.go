package taskspec

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Parser is responsible for parsing taskspec annotations.
type Parser struct {
	// Options can be added here for future extensibility
}

// NewParser creates a new Parser instance.
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses a single line or multiline taskspec annotation.
func (p *Parser) Parse(text string) (*Task, error) {
	if text == "" {
		return nil, nil
	}

	task := &Task{
		Raw:          text,
		CustomFields: make(map[string]string),
	}

	// Trim leading/trailing whitespace
	text = strings.TrimSpace(text)

	// Check if it's a Markdown task list
	if p.isMarkdownTask(text) {
		return p.parseMarkdownTask(text, task)
	}

	// Parse standard format: KEYWORD: description [metadata]
	return p.parseStandardTask(text, task)
}

// isMarkdownTask checks if the text is a Markdown task list item.
func (p *Parser) isMarkdownTask(text string) bool {
	return strings.HasPrefix(text, "- [ ]") || strings.HasPrefix(text, "- [x]") ||
		strings.HasPrefix(text, "- [X]")
}

// parseMarkdownTask parses a Markdown task list item.
func (p *Parser) parseMarkdownTask(text string, task *Task) (*Task, error) {
	task.IsMarkdownTask = true

	// Check if task is checked
	if strings.HasPrefix(text, "- [x]") || strings.HasPrefix(text, "- [X]") {
		task.IsChecked = true
		text = strings.TrimPrefix(text, "- [x]")
		text = strings.TrimPrefix(text, "- [X]")
	} else {
		text = strings.TrimPrefix(text, "- [ ]")
	}

	text = strings.TrimSpace(text)

	// Parse description and metadata
	return p.parseDescriptionAndMetadata(text, task)
}

// parseStandardTask parses a standard taskspec annotation.
func (p *Parser) parseStandardTask(text string, task *Task) (*Task, error) {
	// Match keyword pattern (case-insensitive)
	keywordRegex := regexp.MustCompile(`^(?i)(TODO|FIXME|BUG|HACK|NOTE|INFO|IDEA|REFACTOR|REMINDER):\s*(.*)$`)
	matches := keywordRegex.FindStringSubmatch(text)

	if len(matches) < 3 {
		// Not a valid taskspec annotation
		return nil, nil
	}

	task.Keyword = strings.ToUpper(matches[1])
	remaining := matches[2]

	// Parse description and metadata
	return p.parseDescriptionAndMetadata(remaining, task)
}

// parseDescriptionAndMetadata extracts description and metadata from the remaining text.
func (p *Parser) parseDescriptionAndMetadata(text string, task *Task) (*Task, error) {
	// Find the first metadata tag
	metadataStart := p.findFirstMetadataTag(text)

	if metadataStart == -1 {
		// No metadata, everything is description
		task.Description = p.unescapeDescription(strings.TrimSpace(text))
		return task, nil
	}

	// Split description and metadata
	task.Description = p.unescapeDescription(strings.TrimSpace(text[:metadataStart]))
	metadataText := text[metadataStart:]

	// Parse metadata
	p.parseMetadata(metadataText, task)

	return task, nil
}

// findFirstMetadataTag finds the position of the first metadata tag in the text.
func (p *Parser) findFirstMetadataTag(text string) int {
	// List of metadata patterns to search for
	patterns := []string{
		`due:`, `📅`,
		`scheduled:`, `⏳`,
		`start:`, `🛫`,
		`priority:`, `p:`, `🔺`, `⏫`, `🔼`, `🔽`, `⏬`,
		`repeat:`, `rec:`, `🔁`,
		`id:`, `🆔`,
		`@`, `👤`,
		`#`, `+`,
		`status:`, `✅`, `🚧`, `❌`, `⬜`, `🚫`,
		`created:`, `➕`,
		`done:`,
		`estimate:`, `⏱️`,
	}

	minPos := -1
	for _, pattern := range patterns {
		// Check for non-escaped occurrences
		pos := 0
		for {
			idx := strings.Index(text[pos:], pattern)
			if idx == -1 {
				break
			}
			actualPos := pos + idx
			// Check if it's escaped
			if actualPos > 0 && text[actualPos-1] == '\\' {
				pos = actualPos + 1
				continue
			}
			if minPos == -1 || actualPos < minPos {
				minPos = actualPos
			}
			break
		}
	}

	// Also check for custom key:value patterns ([\w-]+:)
	// This catches custom fields like custom-field:value
	re := regexp.MustCompile(`([\w-]+):\S`)
	matches := re.FindAllStringIndex(text, -1)
	for _, match := range matches {
		pos := match[0]
		// Check if it's escaped
		if pos > 0 && text[pos-1] == '\\' {
			continue
		}
		// Skip if it's part of a URL (has // before it)
		if pos > 1 && text[pos-2:pos] == "//" {
			continue
		}
		if minPos == -1 || pos < minPos {
			minPos = pos
		}
	}

	return minPos
}

// parseMetadata parses all metadata fields from the text.
func (p *Parser) parseMetadata(text string, task *Task) {
	// Remove escaped characters first
	text = p.unescapeText(text)

	// Parse various metadata fields
	p.parseDates(text, task)
	p.parsePriority(text, task)
	p.parseRecurrence(text, task)
	p.parseIdentifier(text, task)
	p.parseAssignees(text, task)
	p.parseTagsAndProjects(text, task)
	p.parseStatus(text, task)
	p.parseEstimate(text, task)
	p.parseCustomFields(text, task)
}

// parseDates parses date fields from the metadata text.
func (p *Parser) parseDates(text string, task *Task) {
	// Due date
	if date := p.extractDate(text, []string{`due:`, `📅`}); date != nil {
		task.DueDate = date
	}

	// Scheduled date
	if date := p.extractDate(text, []string{`scheduled:`, `⏳`}); date != nil {
		task.ScheduledDate = date
	}

	// Start date
	if date := p.extractDate(text, []string{`start:`, `🛫`}); date != nil {
		task.StartDate = date
	}

	// Created date
	if date := p.extractDate(text, []string{`created:`, `➕`}); date != nil {
		task.CreatedDate = date
	}

	// Completed date
	if date := p.extractDate(text, []string{`done:`, `✅`}); date != nil {
		task.CompletedDate = date
	}
}

// extractDate extracts a date value following one of the given prefixes.
func (p *Parser) extractDate(text string, prefixes []string) *time.Time {
	for _, prefix := range prefixes {
		pattern := regexp.QuoteMeta(prefix) + `\s*(\d{4}-\d{2}-\d{2}(?:T[\d:]+(?:Z|[+-]\d{2}:\d{2})?)?)`
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(text)
		if len(matches) >= 2 {
			dateStr := matches[1]
			// Try parsing with different formats
			formats := []string{
				"2006-01-02",
				time.RFC3339,
				"2006-01-02T15:04:05Z",
				"2006-01-02T15:04:05",
			}
			for _, format := range formats {
				if t, err := time.Parse(format, dateStr); err == nil {
					return &t
				}
			}
		}
	}
	return nil
}

// parsePriority parses priority from the metadata text.
func (p *Parser) parsePriority(text string, task *Task) {
	// Check emoji priority
	emojiMap := map[string]Priority{
		`🔺`: PriorityHighest,
		`⏫`: PriorityHigh,
		`🔼`: PriorityMedium,
		`🔽`: PriorityLow,
		`⏬`: PriorityLowest,
	}

	for emoji, priority := range emojiMap {
		if strings.Contains(text, emoji) {
			task.Priority = priority
			return
		}
	}

	// Check text priority
	prefixes := []string{`priority:`, `p:`}
	for _, prefix := range prefixes {
		pattern := regexp.QuoteMeta(prefix) + `\s*(\w+)`
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(text)
		if len(matches) >= 2 {
			priorityStr := strings.ToLower(matches[1])
			switch priorityStr {
			case "highest", "critical", "1":
				task.Priority = PriorityHighest
			case "high", "2":
				task.Priority = PriorityHigh
			case "medium", "normal", "3":
				task.Priority = PriorityMedium
			case "low", "4":
				task.Priority = PriorityLow
			case "lowest", "5":
				task.Priority = PriorityLowest
			}
			return
		}
	}
}

// parseRecurrence parses recurrence patterns from the metadata text.
func (p *Parser) parseRecurrence(text string, task *Task) {
	prefixes := []string{`repeat:`, `rec:`, `🔁`}
	for _, prefix := range prefixes {
		pattern := regexp.QuoteMeta(prefix) + `\s*([^\s]+(?:\s+[^\s]+)*?)(?:\s+(?:due:|scheduled:|start:|priority:|p:|id:|@|#|\+|status:|created:|done:|estimate:|🔺|⏫|🔼|🔽|⏬|📅|⏳|🛫|🔁|🆔|👤|✅|🚧|❌|⬜|🚫|➕|⏱️)|$)`
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(text)
		if len(matches) >= 2 {
			task.Recurrence = strings.TrimSpace(matches[1])
			return
		}
	}
}

// parseIdentifier parses the task identifier from the metadata text.
func (p *Parser) parseIdentifier(text string, task *Task) {
	prefixes := []string{`id:`, `🆔`}
	for _, prefix := range prefixes {
		pattern := regexp.QuoteMeta(prefix) + `\s*(\S+)`
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(text)
		if len(matches) >= 2 {
			task.Identifier = matches[1]
			return
		}
	}
}

// parseAssignees parses assignees from the metadata text.
func (p *Parser) parseAssignees(text string, task *Task) {
	// Match @username or 👤username
	re := regexp.MustCompile(`(?:@|👤)(\w+)`)
	matches := re.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			task.Assignees = append(task.Assignees, match[1])
		}
	}
}

// parseTagsAndProjects parses tags and projects from the metadata text.
func (p *Parser) parseTagsAndProjects(text string, task *Task) {
	// Match #tag
	tagRe := regexp.MustCompile(`#(\w+)`)
	tagMatches := tagRe.FindAllStringSubmatch(text, -1)
	for _, match := range tagMatches {
		if len(match) >= 2 {
			task.Tags = append(task.Tags, match[1])
		}
	}

	// Match +project (need to escape + as it's a regex metacharacter)
	projectRe := regexp.MustCompile(`\+(\w+)`)
	projectMatches := projectRe.FindAllStringSubmatch(text, -1)
	for _, match := range projectMatches {
		if len(match) >= 2 {
			task.Projects = append(task.Projects, match[1])
		}
	}
}

// parseStatus parses status from the metadata text.
func (p *Parser) parseStatus(text string, task *Task) {
	// Check emoji status
	if strings.Contains(text, `⬜`) {
		task.Status = StatusTodo
		return
	}
	if strings.Contains(text, `🚧`) {
		task.Status = StatusInProgress
		return
	}
	if strings.Contains(text, `✅`) && !strings.Contains(text, `done:`) {
		// Only if not part of done: date field
		task.Status = StatusDone
		return
	}
	if strings.Contains(text, `❌`) {
		task.Status = StatusCancelled
		return
	}
	if strings.Contains(text, `🚫`) {
		task.Status = StatusBlocked
		return
	}

	// Check text status
	re := regexp.MustCompile(`status:\s*(\S+)`)
	matches := re.FindStringSubmatch(text)
	if len(matches) >= 2 {
		task.Status = Status(strings.ToLower(matches[1]))
	}
}

// parseEstimate parses time estimate from the metadata text.
func (p *Parser) parseEstimate(text string, task *Task) {
	prefixes := []string{`estimate:`, `⏱️`}
	for _, prefix := range prefixes {
		pattern := regexp.QuoteMeta(prefix) + `\s*(\S+)`
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(text)
		if len(matches) >= 2 {
			task.Estimate = matches[1]
			return
		}
	}
}

// parseCustomFields parses any custom metadata fields.
func (p *Parser) parseCustomFields(text string, task *Task) {
	// Match key:value patterns that aren't standard fields (allow hyphens in key names)
	re := regexp.MustCompile(`([\w-]+):(\S+)`)
	matches := re.FindAllStringSubmatch(text, -1)

	standardFields := map[string]bool{
		"due": true, "scheduled": true, "start": true, "priority": true, "p": true,
		"repeat": true, "rec": true, "id": true, "status": true, "created": true,
		"done": true, "estimate": true,
	}

	for _, match := range matches {
		if len(match) >= 3 {
			key := strings.ToLower(match[1])
			if !standardFields[key] {
				task.CustomFields[key] = match[2]
			}
		}
	}
}

// unescapeText removes escape characters from the text.
func (p *Parser) unescapeText(text string) string {
	// Replace \# with # etc.
	replacements := map[string]string{
		`\#`: `#`,
		`\@`: `@`,
		`\+`: `+`,
	}
	for escaped, unescaped := range replacements {
		text = strings.ReplaceAll(text, escaped, unescaped)
	}
	return text
}

// unescapeDescription removes escape characters from description text.
func (p *Parser) unescapeDescription(text string) string {
	return p.unescapeText(text)
}

// ParseLines parses multiple lines of text and returns all found tasks.
func (p *Parser) ParseLines(lines []string) ([]*Task, error) {
	var tasks []*Task
	var currentLines []string
	var inMultiline bool

	for i, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines
		if line == "" {
			if inMultiline && len(currentLines) > 0 {
				// End of multiline comment
				combined := strings.Join(currentLines, " ")
				if task, err := p.Parse(combined); err == nil && task != nil {
					tasks = append(tasks, task)
				}
				currentLines = nil
				inMultiline = false
			}
			continue
		}

		// Try to parse as a single line first
		task, err := p.Parse(line)
		if err != nil {
			return nil, err
		}

		if task != nil {
			tasks = append(tasks, task)
		} else if i < len(lines)-1 {
			// Check if this might be a multiline continuation
			if p.looksLikeComment(line) {
				currentLines = append(currentLines, line)
				inMultiline = true
			}
		}
	}

	// Handle any remaining multiline content
	if len(currentLines) > 0 {
		combined := strings.Join(currentLines, " ")
		if task, err := p.Parse(combined); err == nil && task != nil {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

// looksLikeComment checks if a line looks like it could be part of a comment.
func (p *Parser) looksLikeComment(line string) bool {
	commentPrefixes := []string{"//", "#", "/*", "*", "--", "<!--"}
	for _, prefix := range commentPrefixes {
		if strings.HasPrefix(line, prefix) {
			return true
		}
	}
	return false
}

// ParseInt is a helper to convert priority string to int
func ParseInt(s string) (int, error) {
	return strconv.Atoi(s)
}
