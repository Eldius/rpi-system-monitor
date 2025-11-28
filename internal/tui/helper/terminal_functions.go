package helper

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

type WindowSize struct {
	Width  int
	Height int
}

func GetTerminalSize() (WindowSize, error) {
	ws := WindowSize{
		Width:  36,
		Height: 8,
	}

	fd := int(os.Stdin.Fd())

	// Check if the file descriptor is a terminal
	if term.IsTerminal(fd) {
		// Get the terminal size
		w, h, err := term.GetSize(fd)
		if err != nil {
			fmt.Printf("Error getting terminal size: %v\n", err)
			return ws, err
		}

		ws.Width = w
		ws.Height = h
		fmt.Printf("Terminal dimensions: Width = %d, Height = %d\n", w, h)
	}

	return ws, nil
}
