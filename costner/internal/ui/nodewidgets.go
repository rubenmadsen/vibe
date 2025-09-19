package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"costner/pkg/types"
)

type NodeWidget struct {
	node      types.Node
	container *fyne.Container
	position  fyne.Position
}

func NewNodeWidget(node types.Node, position fyne.Position) *NodeWidget {
	w := &NodeWidget{
		node:     node,
		position: position,
	}
	w.createWidget()
	return w
}

func (w *NodeWidget) createWidget() {
	// Create node header
	header := widget.NewCard(w.node.Type(), w.node.Name(), nil)
	header.Resize(fyne.NewSize(200, 40))

	// Create inputs section
	inputs := w.createInputsSection()

	// Create outputs section
	outputs := w.createOutputsSection()

	// Create main container
	w.container = container.NewVBox(
		header,
		widget.NewSeparator(),
		inputs,
		widget.NewSeparator(),
		outputs,
	)

	w.container.Move(w.position)
	w.container.Resize(fyne.NewSize(220, 300))
}

func (w *NodeWidget) createInputsSection() *fyne.Container {
	inputsContainer := container.NewVBox()
	inputsContainer.Add(widget.NewLabel("Inputs:"))

	for _, input := range w.node.GetInputs() {
		inputWidget := w.createInputWidget(input)
		inputsContainer.Add(inputWidget)
	}

	return inputsContainer
}

func (w *NodeWidget) createOutputsSection() *fyne.Container {
	outputsContainer := container.NewVBox()
	outputsContainer.Add(widget.NewLabel("Outputs:"))

	for _, output := range w.node.GetOutputs() {
		outputWidget := w.createOutputWidget(output)
		outputsContainer.Add(outputWidget)
	}

	return outputsContainer
}

func (w *NodeWidget) createInputWidget(input types.NodeInput) *fyne.Container {
	label := widget.NewLabel(input.Name)
	label.Resize(fyne.NewSize(80, 20))

	var valueWidget fyne.CanvasObject

	switch input.Type {
	case "bool":
		check := widget.NewCheck("", func(checked bool) {
			w.node.SetInputValue(input.Name, checked)
		})
		if val, ok := input.Value.(bool); ok {
			check.SetChecked(val)
		}
		valueWidget = check

	case "string":
		entry := widget.NewEntry()
		if val, ok := input.Value.(string); ok {
			entry.SetText(val)
		}
		entry.OnChanged = func(text string) {
			w.node.SetInputValue(input.Name, text)
		}
		valueWidget = entry

	case "int":
		entry := widget.NewEntry()
		if val, ok := input.Value.(int); ok {
			entry.SetText(fmt.Sprintf("%d", val))
		}
		entry.OnChanged = func(text string) {
			// TODO: Parse int and set value
		}
		valueWidget = entry

	default:
		entry := widget.NewEntry()
		if input.Value != nil {
			entry.SetText(fmt.Sprintf("%v", input.Value))
		}
		entry.OnChanged = func(text string) {
			w.node.SetInputValue(input.Name, text)
		}
		valueWidget = entry
	}

	// Create connection point
	connectionPoint := widget.NewButton("◯", func() {
		// TODO: Handle connection creation
	})
	connectionPoint.Resize(fyne.NewSize(20, 20))

	return container.NewHBox(
		connectionPoint,
		label,
		valueWidget,
	)
}

func (w *NodeWidget) createOutputWidget(output types.NodeOutput) *fyne.Container {
	label := widget.NewLabel(output.Name)
	label.Resize(fyne.NewSize(80, 20))

	valueLabel := widget.NewLabel("")
	if output.Value != nil {
		valueLabel.SetText(fmt.Sprintf("%v", output.Value))
	}

	// Create connection point
	connectionPoint := widget.NewButton("◯", func() {
		// TODO: Handle connection creation
	})
	connectionPoint.Resize(fyne.NewSize(20, 20))

	return container.NewHBox(
		label,
		valueLabel,
		connectionPoint,
	)
}

func (w *NodeWidget) Container() *fyne.Container {
	return w.container
}

func (w *NodeWidget) Node() types.Node {
	return w.node
}

func (w *NodeWidget) SetPosition(pos fyne.Position) {
	w.position = pos
	w.container.Move(pos)
}