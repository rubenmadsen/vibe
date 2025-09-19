# Build Instructions for Costner

## Current Status

The GUI version (`costner.exe`) was cross-compiled from WSL/Linux, which has limitations with Fyne's OpenGL dependencies.

## Working Versions

✅ **CLI Version Works**: `costner-cli.exe` - Fully functional command-line interface

## To Get GUI Working on Windows

### Option 1: Build Directly on Windows
1. Install Go on Windows machine
2. Copy the entire `costner/` source folder to Windows
3. Run: `go build -o costner-gui.exe ./cmd/costner`

### Option 2: Use CLI for Now
The CLI version is fully functional:
```cmd
costner-cli.exe run example.costner
costner-cli.exe validate example.costner
costner-cli.exe list-nodes
```

### Option 3: Web Interface (Future)
Consider implementing a web-based GUI that runs locally and opens in browser.

## Files to Copy to Windows

**Essential:**
- `costner-cli.exe` (working CLI version)
- `example.costner` (test file)

**For GUI development:**
- Entire `costner/` source folder
- Build on Windows with: `go build ./cmd/costner`

## Testing the CLI

```cmd
# Test the CLI version works
costner-cli.exe help
costner-cli.exe run example.costner
```

The CLI provides all the core functionality - the GUI is just a visual interface on top of the same engine.