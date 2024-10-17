package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/imyounas/tuitaskapp/internal/bubbletasks"
	"github.com/imyounas/tuitaskapp/internal/database"
)

func main() {

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("./logs/debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	// in real app these will come from db
	tasks := []database.Task{
		{Id: 1, Name: "Code Review", Description: "Review the PR backend", Notes: "Focus on API design", Assignee: "Alice"},
		{Id: 2, Name: "Update Documentation", Description: "Update API docs", Notes: "Add examples, cases", Assignee: "Bob"},
		{Id: 3, Name: "Implement Feature X", Description: "Develop new feature", Notes: "Follow design doc", Assignee: "Charlie"},
		{Id: 4, Name: "Fix Bug Y", Description: "Fix memory leak", Notes: "Check logs, dumps", Assignee: "David"},
		{Id: 5, Name: "Set Up CI/CD", Description: "Create CI/CD pipeline", Notes: "Use Docker template", Assignee: "Eve"},
		{Id: 6, Name: "Optimize Queries", Description: "Optimize slow queries", Notes: "Review indexes, plan", Assignee: "Frank"},
		{Id: 7, Name: "Design New Microservice", Description: "Create service design", Notes: "Scalability, recovery", Assignee: "Grace"},
		{Id: 8, Name: "Security Audit", Description: "Audit authentication", Notes: "Check token expiry", Assignee: "Alice"},
	}

	p := tea.NewProgram(bubbletasks.InitialMainModel(tasks), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
