package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"tty-diary/config"
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
		if fileExistsAndNotEmpty(m.cal.CurrentValue(), m.config) {
			m.cal.Colors[m.cal.CurrentValue()] = m.config.NotesColor
		}

		m.preview.RenderFile(pathToMd(m.cal.CurrentValue(), m.config))
	case tea.KeyMsg:
		switch msg.String() {
		case "?":
			m.help.ShowAll = !m.help.ShowAll
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "enter":
			return m, openEditor(pathToMd(m.cal.CurrentValue(), m.config))
		default:
			m.cal.Update(msg)
			m.preview.RenderFile(pathToMd(m.cal.CurrentValue(), m.config))
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

// TODO: refactor. this 3 methods to module `fileScrapper` with init of dir and file format
func getDatesWithFiles(startYear, endYear int, config config.Config) []string {
	var dates []string

	for year := startYear; year <= endYear; year++ {
		for month := 1; month <= 12; month++ {
			daysInMonth := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.
				UTC).Day()
			for day := 1; day <= daysInMonth; day++ {
				date := fmt.Sprintf("%04d/%02d/%02d", year, month, day)
				if fileExistsAndNotEmpty(date, config) {
					dates = append(dates, date)
				}
			}
		}
	}

	return dates
}

func fileExistsAndNotEmpty(date string, config config.Config) bool {
	path := pathToMd(date, config)
	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return false
		}

		return len(string(data)) > 0
	}
	return false
}

func pathToMd(date string, config config.Config) string {
	return filepath.Join(config.DiaryDir, date+"."+config.FileFormat)
}

func main() {
	config := config.ValidateFlags()
	config.DatepickerConfig.HideHelp = true

	colors := make(datepicker.Colors)
	for _, v := range getDatesWithFiles(time.Now().Year()-1, time.Now().Year()+1, config) {
		colors[v] = config.NotesColor
	}

	cal := datepicker.InitModel(config.DatepickerConfig, colors)

	fileForToday := pathToMd(cal.CurrentValue(), config)
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
	}
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
