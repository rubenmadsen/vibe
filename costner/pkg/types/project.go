package types

import "time"

type Project struct {
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Nodes       []NodeData             `json:"nodes"`
	Connections []Connection           `json:"connections"`
	Variables   map[string]interface{} `json:"variables"`
}

type NodeData struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Name     string                 `json:"name"`
	Position Position               `json:"position"`
	Config   map[string]interface{} `json:"config"`
	Inputs   []NodeInput            `json:"inputs"`
	Outputs  []NodeOutput           `json:"outputs"`
}

type ExecutionResult struct {
	NodeID    string                 `json:"node_id"`
	Success   bool                   `json:"success"`
	Error     string                 `json:"error,omitempty"`
	Outputs   map[string]interface{} `json:"outputs"`
	Duration  time.Duration          `json:"duration"`
	Timestamp time.Time              `json:"timestamp"`
}