package miniapps

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// DirectoryTreeApp represents the mini-app for displaying a directory tree.
type DirectoryTreeApp struct {
	app *tview.Application
}

// NewDirectoryTreeApp creates a new instance of the DirectoryTreeApp.
func NewDirectoryTreeApp(app *tview.Application) *DirectoryTreeApp {
	return &DirectoryTreeApp{
		app: app,
	}
}

// Name returns the name of the mini-app.
func (d *DirectoryTreeApp) Name() string {
	return "Directory Tree"
}

// Widget returns the tview.Primitive for the mini-app.
func (d *DirectoryTreeApp) Widget(onExit func()) tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	messageTextView := tview.NewTextView().SetTextColor(tview.Styles.PrimaryTextColor)
	messageTextView.SetBorder(true).SetTitle("Message")

	inputField := tview.NewInputField().
		SetLabel("Enter directory path: ").
		SetFieldWidth(50)

	// Function to show the input field screen
	showInputFieldScreen := func() {
		flex.Clear().
			AddItem(inputField, 3, 0, true).
			AddItem(messageTextView, 3, 0, false).
			AddItem(tview.NewTextView().SetText(""), 0, 1, false). // Placeholder for tree view
			AddItem(tview.NewTextView().SetText("Press ESC to return to menu.").SetTextAlign(tview.AlignCenter), 0, 1, false)
		d.app.SetFocus(inputField)
	}

	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			onExit()
		} else if event.Key() == tcell.KeyEnter {
			path := inputField.GetText()
			if path == "" {
				messageTextView.SetText("Path cannot be empty.")
				return nil
			}
			// Generate and display tree
			treeView, err := d.generateTree(path, showInputFieldScreen)
			if err != nil {
				messageTextView.SetText(fmt.Sprintf("Error: %v", err))
				return nil
			}
			flex.Clear().AddItem(treeView, 0, 1, true)
			d.app.SetFocus(treeView)
		}
		return event
	})

	inputField.SetBorder(true).SetTitle("  Directory Tree App  ")

	// Initial display of the input field screen
	showInputFieldScreen()

	return flex
}

// generateTree generates the directory tree and returns a tview.Primitive.
func (d *DirectoryTreeApp) generateTree(rootPath string, onExit func()) (tview.Primitive, error) {
	fileInfo, err := os.Stat(rootPath)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %s", rootPath)
	}

	tree := tview.NewTreeView()
	tree.SetBorder(true).SetTitle(fmt.Sprintf("Tree: %s", rootPath))

	rootNode := tview.NewTreeNode(rootPath).
		SetReference(rootPath).
		SetSelectable(true) // Make root selectable to expand/collapse

	tree.SetRoot(rootNode).
		SetCurrentNode(rootNode)

	// Function to add nodes recursively
	const maxDepth = 3 // How many levels to expand by default
	var addNodes func(path string, parent *tview.TreeNode, depth int)
	addNodes = func(path string, parent *tview.TreeNode, depth int) {
		files, err := os.ReadDir(path)
		if err != nil {
			// Display error within the tree view
			errorNode := tview.NewTreeNode(fmt.Sprintf("Error reading directory: %v", err)).SetColor(tview.Styles.SecondaryTextColor)
			parent.AddChild(errorNode)
			return
		}

		for _, file := range files {
			fullPath := filepath.Join(path, file.Name())
			node := tview.NewTreeNode(file.Name()).SetReference(fullPath)

			if file.IsDir() {
				node.SetColor(tview.Styles.PrimaryTextColor).SetSelectable(true)
				if depth < maxDepth {
					node.SetExpanded(true)
					addNodes(fullPath, node, depth+1)
				} else {
					// Add a dummy child to make it expandable
					node.AddChild(tview.NewTreeNode("").SetSelectable(false))
				}
				node.SetSelectedFunc(func() {
					// Toggle expansion
					node.SetExpanded(!node.IsExpanded())
					if node.IsExpanded() {
						// Clear existing children and re-add
						node.ClearChildren()
						addNodes(node.GetReference().(string), node, depth+1)
					} else {
						// Add dummy child back when collapsed
						node.AddChild(tview.NewTreeNode("").SetSelectable(false))
					}
				})
			} else {
				// if it's a file
				node.SetColor(tview.Styles.SecondaryTextColor).SetSelectable(false)
			}
			parent.AddChild(node)
		}
	}

	// Call addNodes to populate the tree
	addNodes(rootPath, rootNode, 0)

	// Create a TextView for messages related to opening folders
	messageDisplay := tview.NewTextView().SetTextColor(tview.Styles.PrimaryTextColor).SetTextAlign(tview.AlignCenter)
	messageDisplay.SetBorder(true).SetTitle("Action Message")

	// Create flex layout
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tree, 0, 1, true).
		AddItem(messageDisplay, 3, 0, false). // Add message display below the tree
		AddItem(tview.NewTextView().SetText("Press ESC to return to input field. Press 'o' to open selected folder.").SetTextAlign(tview.AlignCenter), 1, 0, false)

	// Handle input
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			onExit()
			return nil // Consume the event
		} else if event.Rune() == 'o' {
			currentNode := tree.GetCurrentNode()
			if currentNode == nil {
				messageDisplay.SetText("No item selected.")
				return nil
			}

			ref := currentNode.GetReference()
			if ref == nil {
				messageDisplay.SetText("Selected item has no path reference.")
				return nil
			}

			path, ok := ref.(string)
			if !ok {
				messageDisplay.SetText("Invalid path reference type.")
				return nil
			}

			fileInfo, err := os.Stat(path)
			if err != nil {
				messageDisplay.SetText(fmt.Sprintf("Error getting file info: %v", err))
				return nil
			}

			if !fileInfo.IsDir() {
				messageDisplay.SetText("Selected item is not a directory.")
				return nil
			}

			err = openPath(path)
			if err != nil {
				messageDisplay.SetText(fmt.Sprintf("Error opening folder: %v", err))
			} else {
				messageDisplay.SetText(fmt.Sprintf("Opened folder: %s", path))
			}
			return nil // Consume the event
		}
		return event
	})

	return flex, nil
}

// openPath opens the given path using the default file explorer for the OS.
func openPath(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", path)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", path)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Start()
}
