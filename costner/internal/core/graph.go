package core

import (
	"fmt"
	"sync"

	"costner/pkg/types"
)

type Graph struct {
	nodes       map[string]types.Node
	connections []types.Connection
	mutex       sync.RWMutex
}

func NewGraph() *Graph {
	return &Graph{
		nodes:       make(map[string]types.Node),
		connections: make([]types.Connection, 0),
	}
}

func (g *Graph) AddNode(node types.Node) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.nodes[node.ID()] = node
}

func (g *Graph) RemoveNode(nodeID string) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if _, exists := g.nodes[nodeID]; !exists {
		return types.ErrNodeNotFound
	}

	// Remove all connections involving this node
	newConnections := make([]types.Connection, 0)
	for _, conn := range g.connections {
		if conn.SourceNode != nodeID && conn.TargetNode != nodeID {
			newConnections = append(newConnections, conn)
		}
	}
	g.connections = newConnections

	delete(g.nodes, nodeID)
	return nil
}

func (g *Graph) GetNode(nodeID string) (types.Node, bool) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	node, exists := g.nodes[nodeID]
	return node, exists
}

func (g *Graph) GetAllNodes() map[string]types.Node {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	result := make(map[string]types.Node)
	for id, node := range g.nodes {
		result[id] = node
	}
	return result
}

func (g *Graph) AddConnection(conn types.Connection) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// Validate connection
	if err := g.validateConnection(conn); err != nil {
		return err
	}

	// Check for cycles
	if g.wouldCreateCycle(conn) {
		return types.ErrCyclicGraph
	}

	g.connections = append(g.connections, conn)
	return nil
}

func (g *Graph) RemoveConnection(connectionID string) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	for i, conn := range g.connections {
		if conn.ID == connectionID {
			g.connections = append(g.connections[:i], g.connections[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("connection not found: %s", connectionID)
}

func (g *Graph) GetConnections() []types.Connection {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	result := make([]types.Connection, len(g.connections))
	copy(result, g.connections)
	return result
}

func (g *Graph) validateConnection(conn types.Connection) error {
	// Check if source and target nodes exist
	sourceNode, exists := g.nodes[conn.SourceNode]
	if !exists {
		return fmt.Errorf("source node not found: %s", conn.SourceNode)
	}

	targetNode, exists := g.nodes[conn.TargetNode]
	if !exists {
		return fmt.Errorf("target node not found: %s", conn.TargetNode)
	}

	// Check if source port exists
	sourceOutputs := sourceNode.GetOutputs()
	sourcePortExists := false
	for _, output := range sourceOutputs {
		if output.Name == conn.SourcePort {
			sourcePortExists = true
			break
		}
	}
	if !sourcePortExists {
		return fmt.Errorf("source port not found: %s.%s", conn.SourceNode, conn.SourcePort)
	}

	// Check if target port exists
	targetInputs := targetNode.GetInputs()
	targetPortExists := false
	for _, input := range targetInputs {
		if input.Name == conn.TargetPort {
			targetPortExists = true
			break
		}
	}
	if !targetPortExists {
		return fmt.Errorf("target port not found: %s.%s", conn.TargetNode, conn.TargetPort)
	}

	return nil
}

func (g *Graph) wouldCreateCycle(newConn types.Connection) bool {
	// Create adjacency list including the new connection
	adjList := make(map[string][]string)

	// Add existing connections
	for _, conn := range g.connections {
		adjList[conn.SourceNode] = append(adjList[conn.SourceNode], conn.TargetNode)
	}

	// Add new connection
	adjList[newConn.SourceNode] = append(adjList[newConn.SourceNode], newConn.TargetNode)

	// DFS to detect cycle
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for nodeID := range g.nodes {
		if !visited[nodeID] {
			if g.hasCycleDFS(nodeID, adjList, visited, recStack) {
				return true
			}
		}
	}

	return false
}

func (g *Graph) hasCycleDFS(nodeID string, adjList map[string][]string, visited, recStack map[string]bool) bool {
	visited[nodeID] = true
	recStack[nodeID] = true

	for _, neighbor := range adjList[nodeID] {
		if !visited[neighbor] {
			if g.hasCycleDFS(neighbor, adjList, visited, recStack) {
				return true
			}
		} else if recStack[neighbor] {
			return true
		}
	}

	recStack[nodeID] = false
	return false
}

func (g *Graph) GetTopologicalOrder() ([]string, error) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	// Create adjacency list and in-degree count
	adjList := make(map[string][]string)
	inDegree := make(map[string]int)

	// Initialize in-degree for all nodes
	for nodeID := range g.nodes {
		inDegree[nodeID] = 0
		adjList[nodeID] = make([]string, 0)
	}

	// Build adjacency list and calculate in-degrees
	for _, conn := range g.connections {
		adjList[conn.SourceNode] = append(adjList[conn.SourceNode], conn.TargetNode)
		inDegree[conn.TargetNode]++
	}

	// Kahn's algorithm for topological sorting
	queue := make([]string, 0)
	for nodeID, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, nodeID)
		}
	}

	result := make([]string, 0)
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		for _, neighbor := range adjList[current] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// Check if all nodes were processed (no cycles)
	if len(result) != len(g.nodes) {
		return nil, types.ErrCyclicGraph
	}

	return result, nil
}

func (g *Graph) GetDependencies(nodeID string) []string {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	dependencies := make([]string, 0)
	for _, conn := range g.connections {
		if conn.TargetNode == nodeID {
			dependencies = append(dependencies, conn.SourceNode)
		}
	}
	return dependencies
}

func (g *Graph) GetDependents(nodeID string) []string {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	dependents := make([]string, 0)
	for _, conn := range g.connections {
		if conn.SourceNode == nodeID {
			dependents = append(dependents, conn.TargetNode)
		}
	}
	return dependents
}