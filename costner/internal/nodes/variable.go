package nodes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"costner/pkg/types"
)

type VariableNode struct {
	types.BaseNode
}

func NewVariableNode(id string) *VariableNode {
	node := &VariableNode{
		BaseNode: types.BaseNode{
			NodeID:   id,
			NodeType: "variable",
			NodeName: "Variable Assignment",
			Inputs: []types.NodeInput{
				{Name: "source", Type: "any", Required: true, Description: "Source value"},
				{Name: "target_type", Type: "string", Required: true, Description: "Target type (header, query, path, body)"},
				{Name: "target_key", Type: "string", Required: true, Description: "Target key or parameter name"},
				{Name: "format", Type: "string", Required: false, Description: "Format template (optional)"},
			},
			Outputs: []types.NodeOutput{
				{Name: "assignment", Type: "map", Description: "Variable assignment for request"},
			},
			Config: make(map[string]interface{}),
		},
	}
	return node
}

func (n *VariableNode) Execute(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {
	source := inputs["source"]
	targetType, ok := inputs["target_type"].(string)
	if !ok || targetType == "" {
		return nil, fmt.Errorf("target_type is required")
	}

	targetKey, ok := inputs["target_key"].(string)
	if !ok || targetKey == "" {
		return nil, fmt.Errorf("target_key is required")
	}

	format, _ := inputs["format"].(string)

	// Validate target type
	validTypes := []string{"header", "query", "path", "body"}
	if !n.contains(validTypes, targetType) {
		return nil, fmt.Errorf("invalid target_type: %s. Must be one of: %s", targetType, strings.Join(validTypes, ", "))
	}

	// Format the source value if template provided
	var value interface{} = source
	if format != "" {
		formatted, err := n.formatValue(source, format)
		if err != nil {
			return nil, fmt.Errorf("formatting failed: %w", err)
		}
		value = formatted
	}

	assignment := map[string]interface{}{
		"type":  targetType,
		"key":   targetKey,
		"value": value,
	}

	// Update output value
	for i := range n.Outputs {
		if n.Outputs[i].Name == "assignment" {
			n.Outputs[i].Value = assignment
			break
		}
	}

	return map[string]interface{}{
		"assignment": assignment,
	}, nil
}

func (n *VariableNode) formatValue(source interface{}, template string) (string, error) {
	// Simple template formatting
	result := template

	switch v := source.(type) {
	case map[string]interface{}:
		for key, value := range v {
			placeholder := fmt.Sprintf("{%s}", key)
			result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
		}
	case string:
		result = strings.ReplaceAll(result, "{value}", v)
	case int, int64, float64, bool:
		result = strings.ReplaceAll(result, "{value}", fmt.Sprintf("%v", v))
	default:
		result = strings.ReplaceAll(result, "{value}", fmt.Sprintf("%v", v))
	}

	return result, nil
}

func (n *VariableNode) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (n *VariableNode) Clone() types.Node {
	clone := NewVariableNode(n.NodeID)
	clone.NodeName = n.NodeName
	clone.Position = n.Position
	clone.Config = make(map[string]interface{})
	for k, v := range n.Config {
		clone.Config[k] = v
	}
	return clone
}

func (n *VariableNode) Serialize() ([]byte, error) {
	return json.Marshal(n.BaseNode)
}

func (n *VariableNode) Deserialize(data []byte) error {
	return json.Unmarshal(data, &n.BaseNode)
}