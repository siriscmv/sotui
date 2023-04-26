package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var tui *tea.Program

type (
	errMsg  error
	State   int
	LogType int
)

const (
	WaitingForInput State = iota
	WaitingForResponse
	DisplayingAllQuestions
	DisplayingQuestion
	DisplayingAllComments
	DisplayingHelpScreen
)

const (
	Info LogType = iota
	Warning
	Error
)

type model struct {
	table    table.Model
	textarea textarea.Model
	viewport viewport.Model
	spinner  spinner.Model
	mouse    bool
	response SEResponse
	state    State
	err      error
}

var (
	BoldStyle    = lipgloss.NewStyle().Bold(true)
	InfoLogStyle = BoldStyle.Copy().
			Foreground(lipgloss.Color("#a6da95"))
	WarningLogStyle = BoldStyle.Copy().
			Foreground(lipgloss.Color("#eed49f"))
	ErrorLogStyle = BoldStyle.Copy().
			Foreground(lipgloss.Color("#ed8796"))
	UserStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#c6a0f6"))
	FadedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#999999"))
	HelpStyle      = FadedStyle.Copy().Italic(true).Padding(0, 1).Margin(0, 1)
	ContainerStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#c6a0f6")).Padding(1).Margin(1)
)

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "What is your question?"
	ta.Focus()

	ta.Prompt = UserStyle.Render("❯ ")
	ta.CharLimit = 2000

	ta.SetWidth(30)
	ta.SetHeight(1)

	ta.FocusedStyle.CursorLine = ta.FocusedStyle.CursorLine.Copy().UnsetBackground()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	vp := viewport.New(30, 3)
	vp.MouseWheelEnabled = true

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = UserStyle

	table := table.New()

	return model{
		table:    table,
		textarea: ta,
		viewport: vp,
		spinner:  sp,
		response: SEResponse{},
		state:    WaitingForInput,
		err:      nil,
		mouse:    true,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
		spCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	m.spinner, spCmd = m.spinner.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Height = msg.Height - 10
		m.viewport.Width = msg.Width - 6
		m.textarea.SetWidth(msg.Width - 6)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			if m.mouse {
				m.mouse = false
				return m, tea.Sequence(tea.DisableMouse, getLogCmd("Disabled mouse scroll/clicks", Info))
			} else {
				m.mouse = true
				return m, tea.Sequence(tea.EnableMouseCellMotion, getLogCmd("Enabled mouse scroll/clicks", Info))
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:

			return m, tea.Batch(tiCmd, vpCmd, spinner.Tick)
		}

	case errMsg:
		m.err = msg
		return m, nil

	}

	return m, tea.Batch(tiCmd, vpCmd, spCmd)
}

func (m model) View() string {
	var bottomView string

	view := ContainerStyle.Render(
		fmt.Sprintf(
			"%s\n\n%s",
			m.viewport.View(),
			lipgloss.NewStyle().Width(m.viewport.Width).Render(bottomView),
		),
	)

	return view
}

func RunTUI() {
	var m = initialModel()
	tui = tea.NewProgram(m, tea.WithMouseCellMotion())

	if _, err := tui.Run(); err != nil {
		panic(err)
	}
}

func getLogCmd(msg string, logType LogType) tea.Cmd {
	return func() tea.Msg {
		switch logType {
		case Info:
		case Warning:
		case Error:
		}

		return nil
	}
}
