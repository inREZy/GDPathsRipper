package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	appVersion     string = "2.0"
	officialServer string = "www.boomlings.com/database"
)

type Model struct {
	filePicker filepicker.Model
	textInput  textinput.Model

	selectedFile    string
	finishedMessage string
	isSelecting     bool
	isTyping        bool
	isFinished      bool
}

func replacePaths(filePath string, customServer string) string {
	var server string

	fileArray := strings.Split(filePath, "\\")
	fileName := fileArray[len(fileArray)-1]

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "Cannot read the file"
	}

	if strings.Contains(string(fileData), officialServer) {
		if strings.Contains(string(fileData), fmt.Sprintf("https://%s", officialServer)) {
			server = fmt.Sprintf("https://%s", officialServer)
		} else {
			server = fmt.Sprintf("http://%s", officialServer)
		}
	} else {
		return fmt.Sprintf("This file don't contains '%s'", officialServer)
	}

	if len(server) != len(customServer) {
		return "The length of the characters doesn't match"
	}

	data := strings.ReplaceAll(string(fileData), server, customServer)
	data = strings.ReplaceAll(data, base64.StdEncoding.EncodeToString([]byte(server)), base64.StdEncoding.EncodeToString([]byte(customServer)))

	os.Mkdir("dist", 0755)

	newFile, err := os.OpenFile(fmt.Sprintf("dist/%s.%s", strings.Split(fileName, ".")[0]+"_modified", strings.Split(fileName, ".")[1]), os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return "Cannot create the file"
	}

	newFile.Write([]byte(data))
	newFile.Close()

	return fmt.Sprintf("%s is successfully modified!", fileName)
}

func initialModel() Model {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".exe", ".so"}
	fp.CurrentDirectory, _ = os.Getwd()

	ti := textinput.New()
	ti.Placeholder = "Your server path"
	ti.Focus()
	ti.CharLimit = 34
	ti.Width = 34

	return Model{
		filePicker: fp,
		textInput:  ti,

		isSelecting: true,
	}
}

func (m Model) Init() tea.Cmd {
	windowTitle := tea.SetWindowTitle(fmt.Sprintf("GDPathsRipper v%s", appVersion))
	return tea.Batch(m.filePicker.Init(), textinput.Blink, windowTitle)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "esc":
			return m, tea.Quit

		case "enter":
			if m.isTyping && (len(m.textInput.Value()) == 33 || len(m.textInput.Value()) == 34) {
				m.finishedMessage = replacePaths(m.selectedFile, m.textInput.Value())
				m.isTyping = false
				m.isFinished = true
			}
		}

	}

	if m.isSelecting {
		var cmd tea.Cmd

		m.filePicker, cmd = m.filePicker.Update(msg)

		if didSelect, path := m.filePicker.DidSelectFile(msg); didSelect {
			m.selectedFile = path
			m.isSelecting = false
			m.isTyping = true
		}

		return m, cmd
	}

	if m.isTyping {
		var cmd tea.Cmd

		m.textInput, cmd = m.textInput.Update(msg)

		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	var s strings.Builder

	s.WriteString("[ GDPathsRipper v" + appVersion + " created by inREZy ]\n\n")

	if m.isSelecting {
		s.WriteString("Pick a file: \n\n")
		s.WriteString(m.filePicker.View() + "\n")
	} else {
		s.WriteString("Selected file: " + m.filePicker.Styles.Selected.Render(m.selectedFile) + "\n")
	}

	if m.isTyping {
		s.WriteString("Write a server address: \n" + m.textInput.View() + "\n\n")
	} else if !m.isTyping && !m.isSelecting {
		s.WriteString("Server address: " + m.textInput.Value() + "\n\n")
	}

	if m.isFinished {
		s.WriteString("Status: " + m.finishedMessage + "\n\n")
	}

	s.WriteString("Press 'ESC' to quit.")

	return s.String()
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		os.Exit(1)
	}
}
