package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"tty-diary/keymap"
	"tty-diary/preview"
	"tty-diary/config"

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
		if fileExistsAndNotEmpty(m.cal.CurrentValue()) {
			// TODO: use defalul color as one of terminal colors
			// TODO: config
			m.cal.Colors[m.cal.CurrentValue()] = "#b16286"
		}

		m.preview.RenderFile(pathToMd(m.cal.CurrentValue()))
	case tea.KeyMsg:
		switch msg.String() {
		case "?":
			m.help.ShowAll = !m.help.ShowAll
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "enter":
			return m, openEditor(pathToMd(m.cal.CurrentValue()))
		default:
			m.cal.Update(msg)
			m.preview.RenderFile(pathToMd(m.cal.CurrentValue()))
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

func getDatesWithFiles(startYear, endYear int) []string {
	var dates []string

	for year := startYear; year <= endYear; year++ {
		for month := 1; month <= 12; month++ {
			daysInMonth := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.
				UTC).Day()
			for day := 1; day <= daysInMonth; day++ {
				date := fmt.Sprintf("%04d/%02d/%02d", year, month, day)
				if fileExistsAndNotEmpty(date) {
					dates = append(dates, date)
				}
			}
		}
	}

	return dates
}

func fileExistsAndNotEmpty(date string) bool {
	path := pathToMd(date)
	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return false
		}

		return len(string(data)) > 0
	}
	return false
}

func pathToMd(date string) string {
	// TODO: config
	diaryDir := os.Getenv("HOME") + "/code/util/notes/diary"
	if os.Getenv("DIARY_DIR") != "" {
		diaryDir = os.Getenv("DIARY_DIR")
	}

	// TODO: to config. file format (maybe someone want to use txt)
	return filepath.Join(diaryDir, date+".md")
}

func main() {
	datepickerConfig := datepicker_config.ValidateFlags()
	datepickerConfig.HideHelp = true

	colors := make(datepicker.Colors)
	for _, v := range getDatesWithFiles(time.Now().Year()-1, time.Now().Year()+1) {
		// TODO: config
		colors[v] = "#b16286"
	}

	cal := datepicker.InitModel(datepickerConfig, colors)

	fileForToday := pathToMd(cal.CurrentValue())
	preview, err := preview.NewModel(fileForToday)
	if err != nil {
		fmt.Println("Could not initialize Bubble Tea model:", err)
		os.Exit(1)
	}

	model := &mainModel{
		cal:     *cal,
		preview: *preview,
		help:    cal.Help,
	}
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
