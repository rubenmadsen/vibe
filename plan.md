# Costner Implementation Plan

## Project Overview
Costner is a graph-based API testing application built with Go and Fyne, featuring a visual node-based interface for creating and executing API testing workflows.

## Phase 1: Core Architecture Setup

### 1.1 Project Structure
```
costner/
├── cmd/
│   └── costner/
│       └── main.go                 # Application entry point
├── internal/
│   ├── core/
│   │   ├── graph.go               # Graph management logic
│   │   ├── node.go                # Base node interface
│   │   └── executor.go            # Graph execution engine
│   ├── nodes/
│   │   ├── env.go                 # Environment variable node
│   │   ├── request.go             # HTTP request node
│   │   ├── transform.go           # Data transformation node
│   │   ├── conditional.go         # Conditional logic node
│   │   └── variable.go            # Variable assignment node
│   ├── ui/
│   │   ├── canvas.go              # Main canvas for node graph
│   │   ├── nodewidgets.go         # UI widgets for each node type
│   │   └── connections.go         # Visual connections between nodes
│   ├── persistence/
│   │   ├── json.go                # JSON serialization/deserialization
│   │   └── project.go             # Project file management
│   └── cli/
│       ├── runner.go              # CLI execution without GUI
│       └── commands.go            # CLI command definitions
├── pkg/
│   └── types/
│       ├── node.go                # Node type definitions
│       ├── connection.go          # Connection type definitions
│       └── project.go             # Project structure types
├── go.mod
├── go.sum
└── README.md
```

### 1.2 Core Types and Interfaces
- `Node` interface with common methods (Execute, GetInputs, GetOutputs, Serialize)
- `Connection` struct for linking node outputs to inputs
- `Graph` struct containing nodes and connections
- `Project` struct for persistence layer

## Phase 2: Node System Implementation

### 2.1 Base Node Interface
```go
type Node interface {
    ID() string
    Type() string
    Execute(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error)
    GetInputs() []NodeInput
    GetOutputs() []NodeOutput
    Serialize() ([]byte, error)
    Deserialize(data []byte) error
}
```

### 2.2 Node Types Implementation
1. **EnvNode**: Load environment variables from OS or .env files
2. **RequestNode**: Execute HTTP requests with configurable parameters
3. **TransformNode**: Apply data transformations (JSON path extraction, formatting)
4. **ConditionalNode**: Branch execution based on conditions
5. **VariableNode**: Define where variables should be injected in requests

### 2.3 Data Flow System
- Use structured Go types with JSON tags for serialization
- Implement type-safe data passing between nodes
- Support for dynamic typing where needed (interface{} with type assertions)

## Phase 3: Graph Engine

### 3.1 Graph Management
- Add/remove nodes from canvas
- Create/delete connections between nodes
- Validate graph structure (no cycles, valid connections)
- Topological sorting for execution order

### 3.2 Execution Engine
- Dependency resolution
- Parallel execution where possible
- Error handling and propagation
- Progress tracking and logging

## Phase 4: User Interface (Fyne)

### 4.1 Main Canvas
- Scrollable/pannable canvas for node placement
- Grid snap for precise positioning
- Context menus for adding nodes

### 4.2 Node Widgets
- Visual representation of each node type
- Input/output connection points
- Inline editing of node parameters
- Status indicators (idle, running, success, error)

### 4.3 Connection System
- Visual drag-and-drop connection creation
- Connection validation (type compatibility)
- Visual connection lines with proper routing

## Phase 5: Persistence & CLI

### 5.1 JSON Serialization
- Project file format (.costner files)
- Node state serialization
- Connection mapping persistence

### 5.2 CLI Interface
- `costner run project.costner` - Execute project without GUI
- `costner validate project.costner` - Validate project file
- `costner list-nodes` - Show available node types

## Phase 6: Testing & Documentation

### 6.1 Testing Strategy
- Unit tests for each node type
- Integration tests for graph execution
- UI tests for basic interactions

### 6.2 Documentation
- API documentation for node development
- User guide for creating graphs
- CLI usage examples

## Implementation Order

1. **Foundation** (Phase 1): Project structure, core types
2. **Core Logic** (Phase 2-3): Node system and graph engine
3. **Persistence** (Phase 5.1): JSON serialization for headless operation
4. **CLI** (Phase 5.2): Command-line interface
5. **GUI** (Phase 4): Fyne-based visual interface
6. **Polish** (Phase 6): Testing and documentation

## Key Design Decisions

- **Modularity**: Clean separation between core logic, UI, and CLI
- **CLI-first**: Core functionality works without GUI for automation
- **Extensibility**: Plugin-like node system for easy addition of new node types
- **Type Safety**: Strong typing where possible with escape hatches for dynamic data
- **Simplicity**: Start with basic features, build complexity incrementally

## Dependencies
- Go 1.21+
- Fyne v2.4+ for GUI
- Standard library for HTTP requests and JSON handling
- Optional: third-party libraries for advanced HTTP features (OAuth, etc.)