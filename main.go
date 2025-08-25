package main

import (
	"godtui/miniapps"

	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// Register mini apps
	miniApps := []MiniApp{
		miniapps.NewHelloApp(app),
		miniapps.NewCounterApp(app),
		miniapps.NewDirectoryTreeApp(app),
	}

	// Create launcher app
	launcher := NewApp(app, miniApps)

	// Run program
	if err := app.SetRoot(launcher.Menu(), true).Run(); err != nil {
		panic(err)
	}
}
