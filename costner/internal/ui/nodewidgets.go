package ui

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"costner/pkg/types"
)

type NodeWidget struct {
	node         types.Node
	container    *DraggableWidget
	position     fyne.Position
	onRun        func(nodeID string)
	onMove       func(nodeID string, pos fyne.Position)
	lastResult   *types.ExecutionResult
}

func NewNodeWidget(node types.Node, position fyne.Position) *NodeWidget {
	w := &NodeWidget{
		node:     node,
		position: position,
	}
	w.createWidget()
	return w
}

func (w *NodeWidget) SetCallbacks(onRun func(string), onMove func(string, fyne.Position)) {
	w.onRun = onRun
	w.onMove = onMove
}

func (w *NodeWidget) createWidget() {
	// Create background rectangle
	bg := canvas.NewRectangle(theme.ButtonColor())
	bg.StrokeWidth = 2
	bg.StrokeColor = theme.PrimaryColor()

	// Create draggable header
	header := widget.NewLabel(fmt.Sprintf("%s - %s", w.node.Type(), w.node.Name()))
	header.Alignment = fyne.TextAlignCenter

	// Create run button
	runBtn := widget.NewButton("▶ Run", func() {
		if w.onRun != nil {
			w.onRun(w.node.ID())
		}
	})

	// Create inputs section
	inputs := w.createInputsSection()

	// Create outputs section
	outputs := w.createOutputsSection()

	// Create main content
	content := container.NewVBox(
		container.NewHBox(header, runBtn),
		widget.NewSeparator(),
		inputs,
		widget.NewSeparator(),
		outputs,
	)

	// Create container with background and content
	containerContent := container.NewWithoutLayout(bg, content)
	containerContent.Resize(fyne.NewSize(250, 320))

	// Size background to match container
	bg.Resize(fyne.NewSize(250, 320))

	// Size content to fit inside with padding
	content.Move(fyne.NewPos(5, 5))
	content.Resize(fyne.NewSize(240, 310))

	// Create draggable container
	w.container = NewDraggableWidget(containerContent, func(pos fyne.Position) {
		w.position = pos
		if w.onMove != nil {
			w.onMove(w.node.ID(), pos)
		}
	})

	w.container.Move(w.position)
	w.container.Resize(fyne.NewSize(250, 320))
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

func (w *NodeWidget) Container() fyne.CanvasObject {
	return w.container
}

func (w *NodeWidget) Node() types.Node {
	return w.node
}


func parseFloat(s string, defaultVal float32) float32 {
	if s == "" {
		return defaultVal
	}
	if f, err := strconv.ParseFloat(s, 32); err == nil {
		return float32(f)
	}
	return defaultVal
}


func (w *NodeWidget) UpdateResult(result types.ExecutionResult) {
	w.lastResult = &result
	// Update visual status based on result
	// TODO: Change border color based on success/failure
}

func (w *NodeWidget) SetPosition(pos fyne.Position) {
	w.position = pos
	w.container.Move(pos)
	if w.onMove != nil {
		w.onMove(w.node.ID(), pos)
	}
}

// DraggableWidget wraps any widget to make it draggable
type DraggableWidget struct {
	widget.BaseWidget
	content    fyne.CanvasObject
	onMove     func(fyne.Position)
	isDragging bool
	dragStart  fyne.Position
}

func NewDraggableWidget(content fyne.CanvasObject, onMove func(fyne.Position)) *DraggableWidget {
	d := &DraggableWidget{
		content: content,
		onMove:  onMove,
	}
	d.ExtendBaseWidget(d)
	return d
}

func (d *DraggableWidget) CreateRenderer() fyne.WidgetRenderer {
	return &draggableRenderer{
		draggable: d,
		content:   d.content,
	}
}

func (d *DraggableWidget) MouseDown(event *desktop.MouseEvent) {
	if event.Button == desktop.MouseButtonPrimary {
		d.isDragging = true
		d.dragStart = event.Position
	}
}

func (d *DraggableWidget) MouseUp(event *desktop.MouseEvent) {
	d.isDragging = false
}

func (d *DraggableWidget) MouseIn(event *desktop.MouseEvent) {}
func (d *DraggableWidget) MouseOut()                        {}

func (d *DraggableWidget) MouseMoved(event *desktop.MouseEvent) {
	if d.isDragging {
		// Calculate new position
		deltaX := event.Position.X - d.dragStart.X
		deltaY := event.Position.Y - d.dragStart.Y

		currentPos := d.Position()
		newPos := fyne.NewPos(currentPos.X+deltaX, currentPos.Y+deltaY)

		d.Move(newPos)
		if d.onMove != nil {
			d.onMove(newPos)
		}
	}
}

type draggableRenderer struct {
	draggable *DraggableWidget
	content   fyne.CanvasObject
}

func (r *draggableRenderer) Layout(size fyne.Size) {
	r.content.Resize(size)
}

func (r *draggableRenderer) MinSize() fyne.Size {
	return r.content.MinSize()
}

func (r *draggableRenderer) Refresh() {
	r.content.Refresh()
}

func (r *draggableRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.content}
}

func (r *draggableRenderer) Destroy() {}