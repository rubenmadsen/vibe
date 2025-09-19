package nodes

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"costner/pkg/types"
)

type EnvNode struct {
	types.BaseNode
}

func NewEnvNode(id string) *EnvNode {
	node := &EnvNode{
		BaseNode: types.BaseNode{
			NodeID:   id,
			NodeType: "env",
			NodeName: "Environment Variables",
			Inputs: []types.NodeInput{
				{Name: "load_os", Type: "bool", Required: false, Description: "Load OS environment variables", Value: true},
				{Name: "env_file", Type: "string", Required: false, Description: "Path to .env file"},
			},
			Outputs: []types.NodeOutput{
				{Name: "variables", Type: "map", Description: "Environment variables"},
			},
			Config: make(map[string]interface{}),
		},
	}
	return node
}

func (n *EnvNode) Execute(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {
	variables := make(map[string]interface{})

	// Load OS environment variables if requested
	if loadOS, ok := inputs["load_os"].(bool); ok && loadOS {
		for _, env := range os.Environ() {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				variables[parts[0]] = parts[1]
			}
		}
	}

	// Load .env file if specified
	if envFile, ok := inputs["env_file"].(string); ok && envFile != "" {
		if err := n.loadEnvFile(envFile, variables); err != nil {
			return nil, err
		}
	}

	// Update output value
	for i := range n.Outputs {
		if n.Outputs[i].Name == "variables" {
			n.Outputs[i].Value = variables
			break
		}
	}

	return map[string]interface{}{
		"variables": variables,
	}, nil
}

func (n *EnvNode) loadEnvFile(filename string, variables map[string]interface{}) error {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	content, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove quotes if present
			if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
				(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
				value = value[1 : len(value)-1]
			}
			variables[key] = value
		}
	}

	return nil
}

func (n *EnvNode) Clone() types.Node {
	clone := NewEnvNode(n.NodeID)
	clone.NodeName = n.NodeName
	clone.Position = n.Position
	clone.Config = make(map[string]interface{})
	for k, v := range n.Config {
		clone.Config[k] = v
	}
	return clone
}

func (n *EnvNode) Serialize() ([]byte, error) {
	return json.Marshal(n.BaseNode)
}

func (n *EnvNode) Deserialize(data []byte) error {
	return json.Unmarshal(data, &n.BaseNode)
}