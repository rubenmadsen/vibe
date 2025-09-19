package nodes

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"costner/pkg/types"
)

type ConditionalNode struct {
	types.BaseNode
}

func NewConditionalNode(id string) *ConditionalNode {
	node := &ConditionalNode{
		BaseNode: types.BaseNode{
			NodeID:   id,
			NodeType: "conditional",
			NodeName: "Conditional",
			Inputs: []types.NodeInput{
				{Name: "value", Type: "any", Required: true, Description: "Value to evaluate"},
				{Name: "condition", Type: "string", Required: true, Description: "Condition (eq, ne, gt, lt, contains, exists)"},
				{Name: "compare_value", Type: "any", Required: false, Description: "Value to compare against"},
				{Name: "true_output", Type: "any", Required: false, Description: "Output when condition is true"},
				{Name: "false_output", Type: "any", Required: false, Description: "Output when condition is false"},
			},
			Outputs: []types.NodeOutput{
				{Name: "result", Type: "bool", Description: "Condition result"},
				{Name: "output", Type: "any", Description: "Selected output based on condition"},
			},
			Config: make(map[string]interface{}),
		},
	}
	return node
}

func (n *ConditionalNode) Execute(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {
	value := inputs["value"]
	condition, ok := inputs["condition"].(string)
	if !ok || condition == "" {
		return nil, fmt.Errorf("condition is required")
	}

	compareValue := inputs["compare_value"]
	trueOutput := inputs["true_output"]
	falseOutput := inputs["false_output"]

	result, err := n.evaluateCondition(value, condition, compareValue)
	if err != nil {
		return nil, fmt.Errorf("condition evaluation failed: %w", err)
	}

	var output interface{}
	if result {
		output = trueOutput
	} else {
		output = falseOutput
	}

	// Update output values
	for i := range n.Outputs {
		switch n.Outputs[i].Name {
		case "result":
			n.Outputs[i].Value = result
		case "output":
			n.Outputs[i].Value = output
		}
	}

	return map[string]interface{}{
		"result": result,
		"output": output,
	}, nil
}

func (n *ConditionalNode) evaluateCondition(value interface{}, condition string, compareValue interface{}) (bool, error) {
	switch condition {
	case "exists":
		return value != nil, nil
	case "eq":
		return n.equals(value, compareValue), nil
	case "ne":
		return !n.equals(value, compareValue), nil
	case "gt":
		return n.greaterThan(value, compareValue)
	case "lt":
		return n.lessThan(value, compareValue)
	case "contains":
		return n.contains(value, compareValue)
	default:
		return false, fmt.Errorf("unknown condition: %s", condition)
	}
}

func (n *ConditionalNode) equals(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func (n *ConditionalNode) greaterThan(a, b interface{}) (bool, error) {
	aFloat, aErr := n.toFloat64(a)
	bFloat, bErr := n.toFloat64(b)
	if aErr != nil || bErr != nil {
		return false, fmt.Errorf("cannot compare non-numeric values")
	}
	return aFloat > bFloat, nil
}

func (n *ConditionalNode) lessThan(a, b interface{}) (bool, error) {
	aFloat, aErr := n.toFloat64(a)
	bFloat, bErr := n.toFloat64(b)
	if aErr != nil || bErr != nil {
		return false, fmt.Errorf("cannot compare non-numeric values")
	}
	return aFloat < bFloat, nil
}

func (n *ConditionalNode) contains(haystack, needle interface{}) (bool, error) {
	haystackStr, ok := haystack.(string)
	if !ok {
		return false, fmt.Errorf("contains operation requires string haystack")
	}
	needleStr, ok := needle.(string)
	if !ok {
		return false, fmt.Errorf("contains operation requires string needle")
	}
	return strings.Contains(haystackStr, needleStr), nil
}

func (n *ConditionalNode) toFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", v)
	}
}

func (n *ConditionalNode) Clone() types.Node {
	clone := NewConditionalNode(n.NodeID)
	clone.NodeName = n.NodeName
	clone.Position = n.Position
	clone.Config = make(map[string]interface{})
	for k, v := range n.Config {
		clone.Config[k] = v
	}
	return clone
}

func (n *ConditionalNode) Serialize() ([]byte, error) {
	return json.Marshal(n.BaseNode)
}

func (n *ConditionalNode) Deserialize(data []byte) error {
	return json.Unmarshal(data, &n.BaseNode)
}