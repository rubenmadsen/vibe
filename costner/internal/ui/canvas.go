package ui

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"costner/internal/core"
	"costner/internal/nodes"
	"costner/pkg/types"
)

type Canvas struct {
	container    *fyne.Container
	content      *fyne.Container
	background   *canvas.Rectangle
	graph        *core.Graph
	factory      *nodes.NodeFactory
	nodeWidgets  map[string]*NodeWidget
	nextPosition fyne.Position
}

func NewCanvas() *Canvas {
	c := &Canvas{
		graph:        core.NewGraph(),
		factory:      nodes.NewNodeFactory(),
		nodeWidgets:  make(map[string]*NodeWidget),
		nextPosition: fyne.NewPos(50, 50),
	}

	// Create background
	c.background = canvas.NewRectangle(theme.BackgroundColor())
	c.background.Resize(fyne.NewSize(2000, 2000))

	c.content = container.NewWithoutLayout()
	c.content.Add(c.background) // Add background first

	// Create toolbar
	toolbar := c.createToolbar()

	// Create main layout
	c.container = container.NewBorder(
		toolbar, // top
		nil,     // bottom
		nil,     // left
		nil,     // right
		container.NewScroll(c.content), // center
	)

	return c
}

func (c *Canvas) createToolbar() *fyne.Container {
	addBtn := widget.NewButton("Add Node", func() {
		c.showAddNodeDialog()
	})

	runBtn := widget.NewButton("Run", func() {
		c.executeGraph()
	})

	saveBtn := widget.NewButton("Save", func() {
		c.saveProject()
	})

	loadBtn := widget.NewButton("Load", func() {
		c.loadProject()
	})

	return container.NewHBox(addBtn, runBtn, saveBtn, loadBtn)
}

func (c *Canvas) showAddNodeDialog() {
	nodeTypes := c.factory.GetAvailableNodeTypes()

	typeSelect := widget.NewSelect(nodeTypes, nil)
	nameEntry := widget.NewEntry()
	nameEntry.SetText("New Node")

	content := container.NewVBox(
		widget.NewLabel("Select node type:"),
		typeSelect,
		widget.NewLabel("Node name:"),
		nameEntry,
	)

	dialog := widget.NewModalPopUp(content, fyne.CurrentApp().Driver().AllWindows()[0].Canvas())

	createBtn := widget.NewButton("Create", func() {
		if typeSelect.Selected == "" {
			return
		}

		nodeID := fmt.Sprintf("node_%d", len(c.nodeWidgets)+1)
		node, err := c.factory.CreateNode(typeSelect.Selected, nodeID)
		if err != nil {
			return
		}

		node.SetName(nameEntry.Text)
		c.addNodeToCanvas(node, c.getNextPosition())
		dialog.Hide()
	})

	cancelBtn := widget.NewButton("Cancel", func() {
		dialog.Hide()
	})

	buttons := container.NewHBox(createBtn, cancelBtn)
	content.Add(buttons)

	dialog.Resize(fyne.NewSize(300, 200))
	dialog.Show()
}

func (c *Canvas) addNodeToCanvas(node types.Node, position fyne.Position) {
	widget := NewNodeWidget(node, position)
	widget.SetCallbacks(
		func(nodeID string) { c.executeNode(nodeID) },
		func(nodeID string, pos fyne.Position) { /* handle move */ },
	)
	c.nodeWidgets[node.ID()] = widget
	c.graph.AddNode(node)
	c.content.Add(widget.Container())
	c.content.Refresh()
}

func (c *Canvas) executeGraph() {
	executor := core.NewExecutor(c.graph)
	results, err := executor.ExecuteGraph(context.Background())

	if err != nil {
		c.showError("Execution Error", err.Error())
		return
	}

	c.showExecutionResults(results)
}

func (c *Canvas) saveProject() {
	// TODO: Show file dialog and save project
	fmt.Println("Save project functionality to be implemented")
}

func (c *Canvas) loadProject() {
	// TODO: Show file dialog and load project
	fmt.Println("Load project functionality to be implemented")
}

func (c *Canvas) showError(title, message string) {
	dialog := widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel(title),
			widget.NewLabel(message),
			widget.NewButton("OK", func() {}),
		),
		fyne.CurrentApp().Driver().AllWindows()[0].Canvas(),
	)
	dialog.Show()
}

func (c *Canvas) showExecutionResults(results []types.ExecutionResult) {
	content := ""
	for _, result := range results {
		status := "✓"
		if !result.Success {
			status = "✗"
		}
		content += fmt.Sprintf("%s %s (%v)\n", status, result.NodeID, result.Duration)
		if !result.Success {
			content += fmt.Sprintf("  Error: %s\n", result.Error)
		}
	}

	dialog := widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel("Execution Results"),
			widget.NewLabel(content),
			widget.NewButton("OK", func() {}),
		),
		fyne.CurrentApp().Driver().AllWindows()[0].Canvas(),
	)
	dialog.Show()
}

func (c *Canvas) getNextPosition() fyne.Position {
	pos := c.nextPosition

	// Move to next position (grid layout)
	c.nextPosition.X += 250
	if c.nextPosition.X > 800 {
		c.nextPosition.X = 50
		c.nextPosition.Y += 200
	}

	return pos
}

func (c *Canvas) executeNode(nodeID string) {
	executor := core.NewExecutor(c.graph)

	// Get dependencies and execute them first
	deps := c.graph.GetDependencies(nodeID)
	for _, depID := range deps {
		if _, err := executor.ExecuteNode(context.Background(), depID); err != nil {
			c.showError("Dependency Error", fmt.Sprintf("Failed to execute dependency %s: %v", depID, err))
			return
		}
	}

	// Execute the target node
	result, err := executor.ExecuteNode(context.Background(), nodeID)
	if err != nil {
		c.showError("Execution Error", err.Error())
		return
	}

	// Update node widget with result
	if widget, exists := c.nodeWidgets[nodeID]; exists {
		widget.UpdateResult(result)
	}

	// Show result
	c.showNodeResult(result)
}

func (c *Canvas) showNodeResult(result types.ExecutionResult) {
	status := "Success"
	if !result.Success {
		status = "Failed"
	}

	content := fmt.Sprintf("Node: %s\nStatus: %s\nDuration: %v",
		result.NodeID, status, result.Duration)

	if !result.Success {
		content += fmt.Sprintf("\nError: %s", result.Error)
	} else {
		content += "\nOutputs:\n"
		for key, value := range result.Outputs {
			content += fmt.Sprintf("  %s: %v\n", key, value)
		}
	}

	dialog := widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel("Execution Result"),
			widget.NewLabel(content),
			widget.NewButton("OK", func() {}),
		),
		fyne.CurrentApp().Driver().AllWindows()[0].Canvas(),
	)
	dialog.Show()
}

func (c *Canvas) Container() *fyne.Container {
	return c.container
}