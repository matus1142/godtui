package miniapps

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type HelloApp struct {
	app *tview.Application
}

// Constructor
func NewHelloApp(app *tview.Application) *HelloApp {
	return &HelloApp{app: app}
}

func (h *HelloApp) Name() string {
	return "Hello App"
}

func (h *HelloApp) Widget(onExit func()) tview.Primitive {
	text := tview.NewTextView().
		SetText("Welcome to Hello App!\n\nPress ESC to return to menu.").
		SetTextAlign(tview.AlignCenter)

	text.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			onExit() // call back to launcher
		}
		return event
	})

	text.SetBorder(true).SetTitle("Hello App")
	return text
}
