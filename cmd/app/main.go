package main

import (
	"bufio"
	"io"
	"os"
	"strings"

	"texteditor/internal/editor"
	"texteditor/internal/keys"
	"texteditor/internal/render"

	"golang.org/x/term"
)

func main() {
	path := os.Args[1]

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	initialText, _ := readFromFile(path)
	terminalWidth, terminalHeight, _ := term.GetSize(int(os.Stdout.Fd()))
	editorState := editor.New(terminalWidth, terminalHeight, initialText)
	render.Draw(os.Stdout, editorState, keys.KeyEvent{}, "")

	for {
		terminalWidth, terminalHeight, _ := term.GetSize(int(os.Stdout.Fd()))
		editorState.UpdateSize(terminalWidth, terminalHeight)

		keyEvent := keys.ReadKey(os.Stdin)
		action := editorState.HandleKey(keyEvent)
		if action == editor.ActionQuit {
			break
		}

		saveState := ""
		if action == editor.ActionSave {
			err := writeToFile(path, editorState.Text())
			if err != nil {
				saveState = "Error in saving"
			} else {
				saveState = "Saved successfuly"
				editorState.MarkSaved()
			}
		}

		render.Draw(os.Stdout, editorState, keyEvent, saveState)
	}
}

func readFromFile(path string) (string, error) {
	file, _ := os.Open(path)
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
	file, _ := os.Create(path)
	defer file.Close()

	textFileWriter := bufio.NewWriter(file)

	_, err := textFileWriter.WriteString(text)

	if err != nil {
		return err
	}

	err = textFileWriter.Flush()

	if err != nil {
		return err
	}

	return nil
}
