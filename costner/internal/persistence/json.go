package persistence

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"costner/internal/core"
	"costner/internal/nodes"
	"costner/pkg/types"
)

type ProjectPersistence struct {
	factory *nodes.NodeFactory
}

func NewProjectPersistence() *ProjectPersistence {
	return &ProjectPersistence{
		factory: nodes.NewNodeFactory(),
	}
}

func (p *ProjectPersistence) SaveProject(project *types.Project, filePath string) error {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	project.UpdatedAt = time.Now()

	data, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal project: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write project file: %w", err)
	}

	return nil
}

func (p *ProjectPersistence) LoadProject(filePath string) (*types.Project, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read project file: %w", err)
	}

	var project types.Project
	if err := json.Unmarshal(data, &project); err != nil {
		return nil, fmt.Errorf("failed to unmarshal project: %w", err)
	}

	return &project, nil
}

func (p *ProjectPersistence) ProjectToGraph(project *types.Project) (*core.Graph, error) {
	graph := core.NewGraph()

	// Create nodes
	for _, nodeData := range project.Nodes {
		node, err := p.factory.CreateNodeFromData(nodeData)
		if err != nil {
			return nil, fmt.Errorf("failed to create node %s: %w", nodeData.ID, err)
		}
		graph.AddNode(node)
	}

	// Add connections
	for _, conn := range project.Connections {
		if err := graph.AddConnection(conn); err != nil {
			return nil, fmt.Errorf("failed to add connection %s: %w", conn.ID, err)
		}
	}

	return graph, nil
}

func (p *ProjectPersistence) GraphToProject(graph *core.Graph, name, description string) *types.Project {
	nodes := make([]types.NodeData, 0)

	// Convert nodes
	for _, node := range graph.GetAllNodes() {
		nodeData := types.NodeData{
			ID:       node.ID(),
			Type:     node.Type(),
			Name:     node.Name(),
			Inputs:   node.GetInputs(),
			Outputs:  node.GetOutputs(),
		}

		// Get position and config from BaseNode if possible
		if serializable, ok := node.(interface{ Serialize() ([]byte, error) }); ok {
			if data, err := serializable.Serialize(); err == nil {
				var baseNode types.BaseNode
				if json.Unmarshal(data, &baseNode) == nil {
					nodeData.Position = baseNode.Position
					nodeData.Config = baseNode.Config
				}
			}
		}

		nodes = append(nodes, nodeData)
	}

	return &types.Project{
		Name:        name,
		Version:     "1.0.0",
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Nodes:       nodes,
		Connections: graph.GetConnections(),
		Variables:   make(map[string]interface{}),
	}
}

func (p *ProjectPersistence) CreateNewProject(name, description string) *types.Project {
	return &types.Project{
		Name:        name,
		Version:     "1.0.0",
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Nodes:       make([]types.NodeData, 0),
		Connections: make([]types.Connection, 0),
		Variables:   make(map[string]interface{}),
	}
}