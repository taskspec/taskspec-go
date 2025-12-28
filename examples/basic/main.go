package main

import (
	"fmt"
	"log"

	"github.com/taskspec/taskspec-go"
)

func main() {
	parser := taskspec.NewParser()

	fmt.Println("=== Example 1: Simple TODO ===")
	task1, err := parser.Parse("TODO: Fix the login bug")
	if err != nil {
		log.Fatal(err)
	}
	if task1 != nil {
		fmt.Printf("Keyword: %s\n", task1.Keyword)
		fmt.Printf("Description: %s\n", task1.Description)
	}
	fmt.Println()

	fmt.Println("=== Example 2: TODO with Metadata ===")
	task2, err := parser.Parse("TODO: Implement new feature due:2026-03-01 @alice @bob #backend +ProjectX p:high")
	if err != nil {
		log.Fatal(err)
	}
	if task2 != nil {
		fmt.Printf("Description: %s\n", task2.Description)
		fmt.Printf("Due Date: %v\n", task2.DueDate)
		fmt.Printf("Assignees: %v\n", task2.Assignees)
		fmt.Printf("Tags: %v\n", task2.Tags)
		fmt.Printf("Projects: %v\n", task2.Projects)
		fmt.Printf("Priority: %s\n", task2.Priority)
	}
	fmt.Println()

	fmt.Println("=== Example 3: Emoji Format ===")
	task3, err := parser.Parse("FIXME: Refactor authentication module 📅 2026-03-15 🔺 @martin #security")
	if err != nil {
		log.Fatal(err)
	}
	if task3 != nil {
		fmt.Printf("Keyword: %s\n", task3.Keyword)
		fmt.Printf("Description: %s\n", task3.Description)
		fmt.Printf("Due Date: %v\n", task3.DueDate)
		fmt.Printf("Priority: %s\n", task3.Priority)
		fmt.Printf("Assignees: %v\n", task3.Assignees)
		fmt.Printf("Tags: %v\n", task3.Tags)
	}
	fmt.Println()

	fmt.Println("=== Example 4: Markdown Task List ===")
	task4, err := parser.Parse("- [ ] Buy groceries due:2026-03-01 estimate:1h")
	if err != nil {
		log.Fatal(err)
	}
	if task4 != nil {
		fmt.Printf("Is Markdown Task: %v\n", task4.IsMarkdownTask)
		fmt.Printf("Is Checked: %v\n", task4.IsChecked)
		fmt.Printf("Description: %s\n", task4.Description)
		fmt.Printf("Due Date: %v\n", task4.DueDate)
		fmt.Printf("Estimate: %s\n", task4.Estimate)
	}
	fmt.Println()

	fmt.Println("=== Example 5: Completed Markdown Task ===")
	task5, err := parser.Parse("- [x] Complete report 📅 2026-02-28 ✅")
	if err != nil {
		log.Fatal(err)
	}
	if task5 != nil {
		fmt.Printf("Is Markdown Task: %v\n", task5.IsMarkdownTask)
		fmt.Printf("Is Checked: %v\n", task5.IsChecked)
		fmt.Printf("Description: %s\n", task5.Description)
		fmt.Printf("Due Date: %v\n", task5.DueDate)
		fmt.Printf("Status: %s\n", task5.Status)
	}
	fmt.Println()

	fmt.Println("=== Example 6: Complex Task with Multiple Metadata ===")
	task6, err := parser.Parse("TODO: Deploy to production due:2026-03-01 scheduled:2026-02-28 @alice @bob #deployment #critical +Release1.0 id:DEPLOY-123 estimate:4h status:in-progress 🚧")
	if err != nil {
		log.Fatal(err)
	}
	if task6 != nil {
		fmt.Printf("Description: %s\n", task6.Description)
		fmt.Printf("Due Date: %v\n", task6.DueDate)
		fmt.Printf("Scheduled Date: %v\n", task6.ScheduledDate)
		fmt.Printf("Assignees: %v\n", task6.Assignees)
		fmt.Printf("Tags: %v\n", task6.Tags)
		fmt.Printf("Projects: %v\n", task6.Projects)
		fmt.Printf("Identifier: %s\n", task6.Identifier)
		fmt.Printf("Estimate: %s\n", task6.Estimate)
		fmt.Printf("Status: %s\n", task6.Status)
	}
	fmt.Println()

	fmt.Println("=== Example 7: Recurrence Pattern ===")
	task7, err := parser.Parse("TODO: Weekly team meeting repeat:every week @team #meeting")
	if err != nil {
		log.Fatal(err)
	}
	if task7 != nil {
		fmt.Printf("Description: %s\n", task7.Description)
		fmt.Printf("Recurrence: %s\n", task7.Recurrence)
		fmt.Printf("Assignees: %v\n", task7.Assignees)
		fmt.Printf("Tags: %v\n", task7.Tags)
	}
	fmt.Println()

	fmt.Println("=== Example 8: Custom Fields ===")
	task8, err := parser.Parse("TODO: Configure server env:production region:us-east-1 instance-type:t3.large")
	if err != nil {
		log.Fatal(err)
	}
	if task8 != nil {
		fmt.Printf("Description: %s\n", task8.Description)
		fmt.Printf("Custom Fields:\n")
		for key, value := range task8.CustomFields {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}
	fmt.Println()

	fmt.Println("=== Example 9: Escaped Characters ===")
	task9, err := parser.Parse(`TODO: Discuss the \#backend team structure and \@mentions policy`)
	if err != nil {
		log.Fatal(err)
	}
	if task9 != nil {
		fmt.Printf("Description: %s\n", task9.Description)
		fmt.Printf("Tags: %v (should be empty)\n", task9.Tags)
		fmt.Printf("Assignees: %v (should be empty)\n", task9.Assignees)
	}
	fmt.Println()
}
