package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"

	"texteditor/internal/editor"
	"texteditor/internal/keys"
	"texteditor/internal/render"

	"golang.org/x/term"
)

func main() {
	path := ""
	initialText := ""
	if len(os.Args) >= 2 {
		path = os.Args[1]
		initialText, _ = readFromFile(path)
	}

	statusMsg := ""
	statusTTL := 0

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	terminalWidth, terminalHeight, _ := term.GetSize(int(os.Stdout.Fd()))
	editorState := editor.New(terminalWidth, terminalHeight, initialText)
	render.Draw(os.Stdout, editorState, keys.KeyEvent{}, statusMsg)

	for {
		terminalWidth, terminalHeight, _ := term.GetSize(int(os.Stdout.Fd()))
		editorState.UpdateSize(terminalWidth, terminalHeight)

		keyEvent := keys.ReadKey(os.Stdin)
		action := editorState.HandleKey(keyEvent)
		if action == editor.ActionQuit {
			break
		}

		if action == editor.ActionSave {
			statusMsg, err = saveToFile(path, editorState.Text())
			statusTTL = 10
			if err == nil {
				editorState.MarkSaved()
			}
		}

		if statusTTL > 0 {
			statusTTL--
		} else {
			statusMsg = ""
		}

		render.Draw(os.Stdout, editorState, keyEvent, statusMsg)
	}
}

func saveToFile(path string, text string) (string, error) {
	if path == "" {
		return "No filename (Save As not implemented)", errors.New("No filename provided in startup arguments. Save As function not implemented yet.")
	}

	err := writeToFile(path, text)
	if err != nil {
		return "Error in saving", err
	} else {
		return "Saved successfuly", nil
	}
}

func readFromFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var builder strings.Builder
	textFileReader := bufio.NewReader(file)

	for {
		line, err := textFileReader.ReadString('\n')

		if err == io.EOF {
			if len(line) > 0 {
				builder.WriteString(line)
			}
			break
		}

		if err != nil {
			return "", err
		}

		builder.WriteString(line)
	}

	return builder.String(), nil
}

func writeToFile(path string, text string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	textFileWriter := bufio.NewWriter(file)

	_, err = textFileWriter.WriteString(text)

	if err != nil {
		return err
	}

	err = textFileWriter.Flush()

	if err != nil {
		return err
	}

	return nil
}
