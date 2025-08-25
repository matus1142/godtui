package main

import (
	"github.com/rivo/tview"
)

// MiniApp interface all mini apps must implement
type MiniApp interface {
	Name() string
	Widget(onExit func()) tview.Primitive
}

// App struct holds application state
type App struct {
	app      *tview.Application
	menu     *tview.List
	miniApps []MiniApp
}

// Constructor
func NewApp(app *tview.Application, miniApps []MiniApp) *App {
	menu := tview.NewList()
	menu.SetBorder(true)
	menu.SetTitle("Mini App Launcher")

	launcher := &App{
		app:      app,
		miniApps: miniApps,
		menu:     menu,
	}

	// Build menu items
	for index, m := range miniApps {
		mini := m
		launcher.menu.AddItem(mini.Name(), "", '1'+rune(index), func() {
			app.SetRoot(mini.Widget(func() {
				app.SetRoot(launcher.menu, true)
			}), true)
		})
	}
	launcher.menu.AddItem("Quit", "Exit program", 'q', func() {
		launcher.app.Stop()
	})

	return launcher
}

// Expose menu
func (a *App) Menu() *tview.List {
	return a.menu
}
