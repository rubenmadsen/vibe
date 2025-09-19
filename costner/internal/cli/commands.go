package cli

import (
	"flag"
	"fmt"
	"os"
)

type CLI struct {
	runner *Runner
}

func NewCLI() *CLI {
	return &CLI{
		runner: NewRunner(),
	}
}

func (c *CLI) Run() {
	if len(os.Args) < 2 {
		c.printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "run":
		c.runCommand()
	case "validate":
		c.validateCommand()
	case "list-nodes":
		c.listNodesCommand()
	case "help", "-h", "--help":
		c.printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		c.printUsage()
		os.Exit(1)
	}
}

func (c *CLI) runCommand() {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	verbose := fs.Bool("verbose", false, "Enable verbose output")
	fs.Usage = func() {
		fmt.Println("Usage: costner run [options] <project.costner>")
		fmt.Println("Options:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(os.Args[2:]); err != nil {
		os.Exit(1)
	}

	if fs.NArg() != 1 {
		fs.Usage()
		os.Exit(1)
	}

	projectPath := fs.Arg(0)
	if err := c.runner.RunProject(projectPath, *verbose); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func (c *CLI) validateCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: costner validate <project.costner>")
		os.Exit(1)
	}

	projectPath := os.Args[2]
	if err := c.runner.ValidateProject(projectPath); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func (c *CLI) listNodesCommand() {
	c.runner.ListNodeTypes()
}

func (c *CLI) printUsage() {
	fmt.Println("Costner - Graph-based API Testing Tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  costner <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  run <project.costner>     Execute a project file")
	fmt.Println("  validate <project.costner> Validate a project file")
	fmt.Println("  list-nodes               List available node types")
	fmt.Println("  help                     Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  costner run my-api-test.costner")
	fmt.Println("  costner run --verbose my-api-test.costner")
	fmt.Println("  costner validate my-api-test.costner")
}