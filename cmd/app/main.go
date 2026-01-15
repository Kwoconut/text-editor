package main

import (
	"os"

	"golang.org/x/term"
	"texteditor/internal/keys"
	"texteditor/internal/editor"
	"texteditor/internal/render"
)

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	terminalWidth, terminalHeight, _ := term.GetSize(int(os.Stdout.Fd()))
	editorState := editor.New(terminalWidth, terminalHeight)
	render.Draw(os.Stdout, editorState, keys.KeyEvent{})

	for {
		terminalWidth, terminalHeight, _ := term.GetSize(int(os.Stdout.Fd()))
		editorState.UpdateSize(terminalWidth, terminalHeight)

		keyEvent := keys.ReadKey(os.Stdin)
		action := editorState.HandleKey(keyEvent)
		if action == editor.ActionQuit {
			break
		}

		render.Draw(os.Stdout, editorState, keyEvent)
	}
}
