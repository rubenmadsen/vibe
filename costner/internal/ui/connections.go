package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"

	"costner/pkg/types"
)

type ConnectionWidget struct {
	connection types.Connection
	line       *canvas.Line
}

func NewConnectionWidget(conn types.Connection, sourcePos, targetPos fyne.Position) *ConnectionWidget {
	line := canvas.NewLine(theme.PrimaryColor())
	line.StrokeWidth = 2
	line.Position1 = sourcePos
	line.Position2 = targetPos

	return &ConnectionWidget{
		connection: conn,
		line:       line,
	}
}

func (w *ConnectionWidget) UpdatePositions(sourcePos, targetPos fyne.Position) {
	w.line.Position1 = sourcePos
	w.line.Position2 = targetPos
	w.line.Refresh()
}

func (w *ConnectionWidget) Line() *canvas.Line {
	return w.line
}

func (w *ConnectionWidget) Connection() types.Connection {
	return w.connection
}

type ConnectionManager struct {
	connections map[string]*ConnectionWidget
	canvas      *fyne.Container
}

func NewConnectionManager(canvas *fyne.Container) *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*ConnectionWidget),
		canvas:      canvas,
	}
}

func (m *ConnectionManager) AddConnection(conn types.Connection, sourcePos, targetPos fyne.Position) {
	widget := NewConnectionWidget(conn, sourcePos, targetPos)
	m.connections[conn.ID] = widget
	m.canvas.Add(widget.Line())
}

func (m *ConnectionManager) RemoveConnection(connectionID string) {
	if widget, exists := m.connections[connectionID]; exists {
		m.canvas.Remove(widget.Line())
		delete(m.connections, connectionID)
	}
}

func (m *ConnectionManager) UpdateConnection(connectionID string, sourcePos, targetPos fyne.Position) {
	if widget, exists := m.connections[connectionID]; exists {
		widget.UpdatePositions(sourcePos, targetPos)
	}
}

func (m *ConnectionManager) GetConnections() map[string]*ConnectionWidget {
	return m.connections
}