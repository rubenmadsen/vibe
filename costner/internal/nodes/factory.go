package nodes

import (
	"fmt"

	"costner/pkg/types"
)

type NodeFactory struct{}

func NewNodeFactory() *NodeFactory {
	return &NodeFactory{}
}

func (f *NodeFactory) CreateNode(nodeType, id string) (types.Node, error) {
	switch nodeType {
	case "env":
		return NewEnvNode(id), nil
	case "request":
		return NewRequestNode(id), nil
	case "transform":
		return NewTransformNode(id), nil
	case "conditional":
		return NewConditionalNode(id), nil
	case "variable":
		return NewVariableNode(id), nil
	default:
		return nil, fmt.Errorf("unknown node type: %s", nodeType)
	}
}

func (f *NodeFactory) GetAvailableNodeTypes() []string {
	return []string{"env", "request", "transform", "conditional", "variable"}
}

func (f *NodeFactory) CreateNodeFromData(data types.NodeData) (types.Node, error) {
	node, err := f.CreateNode(data.Type, data.ID)
	if err != nil {
		return nil, err
	}

	node.SetName(data.Name)

	// Set input values
	for _, input := range data.Inputs {
		node.SetInputValue(input.Name, input.Value)
	}

	return node, nil
}