package preview

import (
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	viewport viewport.Model
	renderer *glamour.TermRenderer
}

func NewModel(initFile string) (*Model, error) {
	const width = 44

	vp := viewport.New(width, 18)
	vp.Style = lipgloss.NewStyle()

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithPreservedNewLines(),
		glamour.WithWordWrap(width-vp.Style.GetHorizontalFrameSize()-2),
	)
	if err != nil {
		return nil, err
	}

	model := &Model{
		viewport: vp,
		renderer: renderer,
	}

	err = model.RenderFile(initFile)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (e Model) Init() tea.Cmd {
	return nil
}

func (e Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return e, nil
}

func (e Model) View() string {
	return e.viewport.View()
}

func (e *Model) RenderFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		data = []byte("")
	}
	fileContent, err := e.renderer.Render(string(data))
	if err != nil {
		return err
	}
	e.viewport.SetContent(fileContent)
	return nil
}
