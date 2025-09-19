package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"costner/pkg/types"
)

type Executor struct {
	graph   *Graph
	results map[string]map[string]interface{}
	mutex   sync.RWMutex
}

func NewExecutor(graph *Graph) *Executor {
	return &Executor{
		graph:   graph,
		results: make(map[string]map[string]interface{}),
	}
}

func (e *Executor) ExecuteGraph(ctx context.Context) ([]types.ExecutionResult, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// Clear previous results
	e.results = make(map[string]map[string]interface{})

	// Get execution order
	order, err := e.graph.GetTopologicalOrder()
	if err != nil {
		return nil, fmt.Errorf("failed to get execution order: %w", err)
	}

	results := make([]types.ExecutionResult, 0, len(order))

	// Execute nodes in topological order
	for _, nodeID := range order {
		node, exists := e.graph.GetNode(nodeID)
		if !exists {
			return nil, fmt.Errorf("node not found during execution: %s", nodeID)
		}

		result, err := e.executeNode(ctx, node)
		results = append(results, result)

		if err != nil {
			// Stop execution on error
			return results, fmt.Errorf("execution stopped at node %s: %w", nodeID, err)
		}
	}

	return results, nil
}

func (e *Executor) ExecuteNode(ctx context.Context, nodeID string) (types.ExecutionResult, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	node, exists := e.graph.GetNode(nodeID)
	if !exists {
		return types.ExecutionResult{}, types.ErrNodeNotFound
	}

	return e.executeNode(ctx, node)
}

func (e *Executor) executeNode(ctx context.Context, node types.Node) (types.ExecutionResult, error) {
	start := time.Now()

	result := types.ExecutionResult{
		NodeID:    node.ID(),
		Timestamp: start,
	}

	// Prepare inputs from connections
	inputs, err := e.prepareInputs(node)
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		result.Duration = time.Since(start)
		return result, err
	}

	// Execute the node
	outputs, err := node.Execute(ctx, inputs)
	result.Duration = time.Since(start)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result, err
	}

	// Store outputs for dependent nodes
	e.results[node.ID()] = outputs
	result.Success = true
	result.Outputs = outputs

	return result, nil
}

func (e *Executor) prepareInputs(node types.Node) (map[string]interface{}, error) {
	inputs := make(map[string]interface{})

	// Get node's input definitions
	nodeInputs := node.GetInputs()

	// Set default values from node inputs
	for _, input := range nodeInputs {
		if input.Value != nil {
			inputs[input.Name] = input.Value
		}
	}

	// Override with connected values
	connections := e.graph.GetConnections()
	for _, conn := range connections {
		if conn.TargetNode == node.ID() {
			// Get value from source node's output
			sourceOutputs, exists := e.results[conn.SourceNode]
			if !exists {
				if e.isInputRequired(nodeInputs, conn.TargetPort) {
					return nil, fmt.Errorf("required input %s not available from %s", conn.TargetPort, conn.SourceNode)
				}
				continue
			}

			value, exists := sourceOutputs[conn.SourcePort]
			if !exists {
				if e.isInputRequired(nodeInputs, conn.TargetPort) {
					return nil, fmt.Errorf("required output %s not found in %s", conn.SourcePort, conn.SourceNode)
				}
				continue
			}

			inputs[conn.TargetPort] = value
		}
	}

	// Validate required inputs
	for _, input := range nodeInputs {
		if input.Required {
			if _, exists := inputs[input.Name]; !exists {
				return nil, fmt.Errorf("required input %s not provided for node %s", input.Name, node.ID())
			}
		}
	}

	return inputs, nil
}

func (e *Executor) isInputRequired(inputs []types.NodeInput, inputName string) bool {
	for _, input := range inputs {
		if input.Name == inputName {
			return input.Required
		}
	}
	return false
}

func (e *Executor) GetNodeResult(nodeID string) (map[string]interface{}, bool) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	result, exists := e.results[nodeID]
	return result, exists
}

func (e *Executor) ClearResults() {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.results = make(map[string]map[string]interface{})
}