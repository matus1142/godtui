package miniapps

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type CounterApp struct {
	app   *tview.Application
	count int
}

// Constructor
func NewCounterApp(app *tview.Application) *CounterApp {
	return &CounterApp{app: app, count: 0}
}

func (c *CounterApp) Name() string {
	return "Counter App"
}

func (c *CounterApp) Widget(onExit func()) tview.Primitive {
	text := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	// Handle input
	text.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case '+':
			c.count++
		}
		if event.Key() == tcell.KeyEscape {
			onExit()
		}
		return event
	})

	// Update UI in goroutine
	go func() {
		for {
			c.app.QueueUpdateDraw(func() {
				text.SetText(fmt.Sprintf("[green]Counter App[-]\n\nCount: %d\n\nPress + to increment, ESC to return.", c.count))
			})
		}
	}()

	text.SetBorder(true).SetTitle("Counter App")
	return text
}
