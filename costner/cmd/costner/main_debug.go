//go:build debug
// +build debug

package main

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"costner/internal/cli"
	"costner/internal/ui"
)

func main() {
	fmt.Println("Costner starting...")
	fmt.Printf("Arguments: %v\n", os.Args)

	// If no arguments provided (clicked), open GUI
	// If arguments provided, use CLI
	if len(os.Args) == 1 {
		fmt.Println("Starting GUI mode...")
		runGUI()
		return
	}

	fmt.Println("Starting CLI mode...")
	// Run CLI for commands
	cli := cli.NewCLI()
	cli.Run()
}

func runGUI() {
	fmt.Println("Creating Fyne app...")

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("GUI Panic: %v\n", r)
			fmt.Println("Press Enter to exit...")
			fmt.Scanln()
			os.Exit(1)
		}
	}()

	myApp := app.New()
	fmt.Println("App created...")

	fmt.Println("Creating window...")
	window := myApp.NewWindow("Costner - Graph-based API Testing")
	window.Resize(fyne.NewSize(1200, 800))

	fmt.Println("Creating canvas...")
	canvas := ui.NewCanvas()
	window.SetContent(canvas.Container())

	fmt.Println("Showing window...")
	window.ShowAndRun()
}