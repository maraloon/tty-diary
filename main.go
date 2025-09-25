package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"tty-diary/config"
	"tty-diary/filer"
	"tty-diary/keymap"
	"tty-diary/preview"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maraloon/datepicker"
	"golang.org/x/term"
)

type mainModel struct {
	cal     datepicker.Model
	preview preview.Model
	help    help.Model
	config  config.Config
	filer   *filer.Filer
}

type editorFinishedMsg struct{ err error }

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(m.cal.Init(), m.preview.Init())
}

func openEditor(path string) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	c := exec.Command(editor, path)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case editorFinishedMsg:
		if m.filer.FileExistsAndNotEmpty(m.cal.CurrentValue()) {
			m.cal.Colors[m.cal.CurrentValue()] = m.config.FileColor
		}

		m.preview.RenderFile(m.filer.Filepath(m.cal.CurrentValue()))
	case tea.KeyMsg:
		switch msg.String() {
		case "?":
			m.help.ShowAll = !m.help.ShowAll
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "enter":
			return m, openEditor(m.filer.Filepath(m.cal.CurrentValue()))
		default:
			m.cal.Update(msg)
			m.preview.RenderFile(m.filer.Filepath(m.cal.CurrentValue()))
		}
	}
	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	var s string

	s += lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.NewStyle().Render(m.cal.View()),
		lipgloss.NewStyle().Render("--------------------------------------------"),
		lipgloss.NewStyle().Render(m.preview.View()),
		lipgloss.NewStyle().Render("--------------------------------------------"),
	)

	s += "\n" + m.help.View(keymap.Keys) + "\n"

	w, h, _ := term.GetSize(int(os.Stderr.Fd()))
	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, s)
}

func main() {
	config, err := config.ValidateFlags()
	if err != nil {
		fmt.Println("Something wrong with config file", err)
		os.Exit(1)
	}
	
	config.DatepickerConfig.HideHelp = true
	filer := filer.NewFiler(config.DiaryDir, config.FileFormat)

	colors := make(datepicker.Colors)
	for _, v := range filer.GetDatesWithFiles(time.Now().Year()-1, time.Now().Year()+1) {
		colors[v] = config.FileColor
	}

	cal := datepicker.InitModel(config.DatepickerConfig, colors)

	fileForToday := filer.Filepath(cal.CurrentValue())
	preview, err := preview.NewModel(fileForToday)
	if err != nil {
		fmt.Println("Could not initialize Bubble Tea model:", err)
		os.Exit(1)
	}

	model := &mainModel{
		cal:     *cal,
		preview: *preview,
		help:    cal.Help,
		config:  config,
		filer:   filer,
	}
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
