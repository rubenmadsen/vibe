package nodes

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"costner/pkg/types"
)

type TransformNode struct {
	types.BaseNode
}

func NewTransformNode(id string) *TransformNode {
	node := &TransformNode{
		BaseNode: types.BaseNode{
			NodeID:   id,
			NodeType: "transform",
			NodeName: "Data Transform",
			Inputs: []types.NodeInput{
				{Name: "input", Type: "any", Required: true, Description: "Input data to transform"},
				{Name: "operation", Type: "string", Required: true, Description: "Transform operation (json_path, to_string, to_int, format)"},
				{Name: "expression", Type: "string", Required: false, Description: "Expression for the operation"},
			},
			Outputs: []types.NodeOutput{
				{Name: "output", Type: "any", Description: "Transformed data"},
			},
			Config: make(map[string]interface{}),
		},
	}
	return node
}

func (n *TransformNode) Execute(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {
	input, ok := inputs["input"]
	if !ok {
		return nil, fmt.Errorf("input is required")
	}

	operation, ok := inputs["operation"].(string)
	if !ok || operation == "" {
		return nil, fmt.Errorf("operation is required")
	}

	expression, _ := inputs["expression"].(string)

	var result interface{}
	var err error

	switch operation {
	case "json_path":
		result, err = n.jsonPath(input, expression)
	case "to_string":
		result = n.toString(input)
	case "to_int":
		result, err = n.toInt(input)
	case "format":
		result, err = n.format(input, expression)
	default:
		return nil, fmt.Errorf("unknown operation: %s", operation)
	}

	if err != nil {
		return nil, fmt.Errorf("transform failed: %w", err)
	}

	// Update output value
	for i := range n.Outputs {
		if n.Outputs[i].Name == "output" {
			n.Outputs[i].Value = result
			break
		}
	}

	return map[string]interface{}{
		"output": result,
	}, nil
}

func (n *TransformNode) jsonPath(input interface{}, path string) (interface{}, error) {
	if path == "" {
		return input, nil
	}

	// Simple JSON path implementation
	// For now, support simple dot notation like "field.subfield"
	parts := strings.Split(path, ".")
	current := input

	for _, part := range parts {
		switch v := current.(type) {
		case map[string]interface{}:
			var exists bool
			current, exists = v[part]
			if !exists {
				return nil, fmt.Errorf("path not found: %s", part)
			}
		case []interface{}:
			// Handle array indices
			if idx, err := strconv.Atoi(part); err == nil && idx >= 0 && idx < len(v) {
				current = v[idx]
			} else {
				return nil, fmt.Errorf("invalid array index: %s", part)
			}
		default:
			return nil, fmt.Errorf("cannot traverse path on type %T", v)
		}
	}

	return current, nil
}

func (n *TransformNode) toString(input interface{}) string {
	switch v := input.(type) {
	case string:
		return v
	case int, int64, float64:
		return fmt.Sprintf("%v", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case []interface{}, map[string]interface{}:
		if bytes, err := json.Marshal(v); err == nil {
			return string(bytes)
		}
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (n *TransformNode) toInt(input interface{}) (int, error) {
	switch v := input.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("cannot convert %T to int", v)
	}
}

func (n *TransformNode) format(input interface{}, template string) (string, error) {
	if template == "" {
		return n.toString(input), nil
	}

	// Simple template substitution
	result := template
	switch v := input.(type) {
	case map[string]interface{}:
		for key, value := range v {
			placeholder := fmt.Sprintf("{%s}", key)
			result = strings.ReplaceAll(result, placeholder, n.toString(value))
		}
	default:
		result = strings.ReplaceAll(result, "{input}", n.toString(input))
	}

	return result, nil
}

func (n *TransformNode) Clone() types.Node {
	clone := NewTransformNode(n.NodeID)
	clone.NodeName = n.NodeName
	clone.Position = n.Position
	clone.Config = make(map[string]interface{})
	for k, v := range n.Config {
		clone.Config[k] = v
	}
	return clone
}

func (n *TransformNode) Serialize() ([]byte, error) {
	return json.Marshal(n.BaseNode)
}

func (n *TransformNode) Deserialize(data []byte) error {
	return json.Unmarshal(data, &n.BaseNode)
}