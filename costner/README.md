# Costner

A graph-based API testing application that uses a visual node-based interface for creating and executing API testing workflows.

## Features

- **Graph-based workflow**: Connect nodes to create API testing flows
- **5 Node types**: Environment, Request, Transform, Conditional, and Variable nodes
- **CLI-first**: Run tests from command line without GUI
- **JSON persistence**: Save and load projects as `.costner` files

## Installation

Download the appropriate binary for your platform:

**For Windows users:**
- `costner.exe` - Full GUI application (double-click to open GUI, or use from command line)
- `costner-cli.exe` - CLI-only version

**For Linux users:**
- `costner-cli` - CLI version (GUI requires proper display server)

## GUI vs CLI

- **GUI Mode**: Double-click `costner.exe` on Windows to open the visual node editor
- **CLI Mode**: Run `costner <command>` from terminal for automated execution

## Usage

### CLI Commands

```bash
# Execute a project
costner run project.costner

# Execute with verbose output
costner run --verbose project.costner

# Validate a project file
costner validate project.costner

# List available node types
costner list-nodes

# Show help
costner help
```

### Node Types

1. **EnvNode**: Load environment variables from OS or .env files
2. **RequestNode**: Execute HTTP requests with configurable parameters
3. **TransformNode**: Apply data transformations (JSON path extraction, formatting)
4. **ConditionalNode**: Branch execution based on conditions
5. **VariableNode**: Define where variables should be injected in requests

## Example Project File

See `example.costner` for a basic project that tests httpbin.org.

## Project File Format

Projects are saved as JSON files with `.costner` extension containing:
- Project metadata (name, version, description)
- Node definitions with inputs/outputs
- Connections between nodes
- Global variables

## Building from Source

```bash
go build ./cmd/costner
```

For Windows:
```bash
GOOS=windows GOARCH=amd64 go build -o costner.exe ./cmd/costner
```