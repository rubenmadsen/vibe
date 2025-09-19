package cli

import (
	"context"
	"fmt"

	"costner/internal/persistence"
	"costner/internal/core"
	"costner/pkg/types"
)

type Runner struct {
	persistence *persistence.ProjectPersistence
}

func NewRunner() *Runner {
	return &Runner{
		persistence: persistence.NewProjectPersistence(),
	}
}

func (r *Runner) RunProject(projectPath string, verbose bool) error {
	// Load project
	project, err := r.persistence.LoadProject(projectPath)
	if err != nil {
		return fmt.Errorf("failed to load project: %w", err)
	}

	if verbose {
		fmt.Printf("Loaded project: %s\n", project.Name)
		fmt.Printf("Description: %s\n", project.Description)
		fmt.Printf("Nodes: %d, Connections: %d\n\n", len(project.Nodes), len(project.Connections))
	}

	// Convert to graph
	graph, err := r.persistence.ProjectToGraph(project)
	if err != nil {
		return fmt.Errorf("failed to create graph: %w", err)
	}

	// Execute graph
	executor := core.NewExecutor(graph)
	ctx := context.Background()

	if verbose {
		fmt.Println("Executing graph...")
	}

	results, err := executor.ExecuteGraph(ctx)
	if err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	// Display results
	r.displayResults(results, verbose)

	return nil
}

func (r *Runner) ValidateProject(projectPath string) error {
	project, err := r.persistence.LoadProject(projectPath)
	if err != nil {
		return fmt.Errorf("failed to load project: %w", err)
	}

	graph, err := r.persistence.ProjectToGraph(project)
	if err != nil {
		return fmt.Errorf("failed to create graph: %w", err)
	}

	// Check for cycles
	_, err = graph.GetTopologicalOrder()
	if err != nil {
		return fmt.Errorf("graph validation failed: %w", err)
	}

	fmt.Printf("Project %s is valid\n", project.Name)
	fmt.Printf("- %d nodes\n", len(project.Nodes))
	fmt.Printf("- %d connections\n", len(project.Connections))

	return nil
}

func (r *Runner) ListNodeTypes() {
	nodeTypes := []string{"env", "request", "transform", "conditional", "variable"}

	fmt.Println("Available node types:")
	for _, nodeType := range nodeTypes {
		fmt.Printf("- %s\n", nodeType)
	}
}

func (r *Runner) displayResults(results []types.ExecutionResult, verbose bool) {
	fmt.Println("Execution Results:")
	fmt.Println("==================")

	for _, result := range results {
		status := "✓"
		if !result.Success {
			status = "✗"
		}

		fmt.Printf("%s Node: %s (Duration: %v)\n", status, result.NodeID, result.Duration)

		if !result.Success {
			fmt.Printf("  Error: %s\n", result.Error)
		} else if verbose {
			fmt.Printf("  Outputs:\n")
			for key, value := range result.Outputs {
				fmt.Printf("    %s: %v\n", key, value)
			}
		}
		fmt.Println()
	}

	successCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
		}
	}

	fmt.Printf("Summary: %d/%d nodes executed successfully\n", successCount, len(results))
}