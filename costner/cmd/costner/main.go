//go:build !cli
// +build !cli

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
	// If no arguments provided (clicked), open GUI
	// If arguments provided, use CLI
	if len(os.Args) == 1 {
		runGUI()
		return
	}

	// Run CLI for commands
	cli := cli.NewCLI()
	cli.Run()
}

func runGUI() {
	defer func() {
		if r := recover(); r != nil {
			// Create a simple error window if GUI fails
			fmt.Printf("GUI Error: %v\n", r)
			os.Exit(1)
		}
	}()

	myApp := app.New()
	myApp.Metadata().ID = "com.costner.app"
	myApp.Metadata().Name = "Costner"

	window := myApp.NewWindow("Costner - Graph-based API Testing")
	window.Resize(fyne.NewSize(1200, 800))

	canvas := ui.NewCanvas()
	window.SetContent(canvas.Container())

	window.ShowAndRun()
}