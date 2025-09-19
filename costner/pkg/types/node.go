package types

import (
	"context"
	"encoding/json"
)

type NodeInput struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Required    bool        `json:"required"`
	Description string      `json:"description"`
	Value       interface{} `json:"value,omitempty"`
}

type NodeOutput struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Value       interface{} `json:"value,omitempty"`
}

type Node interface {
	ID() string
	Type() string
	Name() string
	SetName(name string)
	Execute(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error)
	GetInputs() []NodeInput
	GetOutputs() []NodeOutput
	SetInputValue(name string, value interface{}) error
	GetOutputValue(name string) (interface{}, bool)
	Serialize() ([]byte, error)
	Deserialize(data []byte) error
	Clone() Node
}

type BaseNode struct {
	NodeID    string                 `json:"id"`
	NodeType  string                 `json:"type"`
	NodeName  string                 `json:"name"`
	Inputs    []NodeInput            `json:"inputs"`
	Outputs   []NodeOutput           `json:"outputs"`
	Config    map[string]interface{} `json:"config"`
	Position  Position               `json:"position"`
}

type Position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

func (b *BaseNode) ID() string {
	return b.NodeID
}

func (b *BaseNode) Type() string {
	return b.NodeType
}

func (b *BaseNode) Name() string {
	return b.NodeName
}

func (b *BaseNode) SetName(name string) {
	b.NodeName = name
}

func (b *BaseNode) GetInputs() []NodeInput {
	return b.Inputs
}

func (b *BaseNode) GetOutputs() []NodeOutput {
	return b.Outputs
}

func (b *BaseNode) SetInputValue(name string, value interface{}) error {
	for i := range b.Inputs {
		if b.Inputs[i].Name == name {
			b.Inputs[i].Value = value
			return nil
		}
	}
	return ErrInputNotFound
}

func (b *BaseNode) GetOutputValue(name string) (interface{}, bool) {
	for _, output := range b.Outputs {
		if output.Name == name {
			return output.Value, true
		}
	}
	return nil, false
}

func (b *BaseNode) Serialize() ([]byte, error) {
	return json.Marshal(b)
}

func (b *BaseNode) Deserialize(data []byte) error {
	return json.Unmarshal(data, b)
}