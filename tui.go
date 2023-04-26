package main

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var tui *tea.Program

type Log struct {
	Msg  string
	Type LogType
}

type (
	errMsg  error
	State   int
	LogType int
	logMsg  Log
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
	WhiteTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	BaseLogStyle   = WhiteTextStyle.Copy().AlignVertical(lipgloss.Center).AlignHorizontal(lipgloss.Center)
	InfoLogStyle   = BaseLogStyle.Copy().
			Background(lipgloss.Color("#a6da9580"))
	WarningLogStyle = BaseLogStyle.Copy().
			Background(lipgloss.Color("#eed49f80"))
	ErrorLogStyle = BaseLogStyle.Copy().
			Background(lipgloss.Color("#ed879680"))
	AccentStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#c6a0f6"))
	FadedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#999999"))
)

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "What is your question?"
	ta.Focus()

	ta.Prompt = AccentStyle.Render("❯ ")
	ta.CharLimit = 200

	ta.SetWidth(30)
	ta.SetHeight(1)

	ta.FocusedStyle.CursorLine = ta.FocusedStyle.CursorLine.Copy().UnsetBackground()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	vp := viewport.New(30, 3)
	vp.MouseWheelEnabled = true

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = AccentStyle

	tb := table.New()
	tb.SetHeight(10)
	tb.SetWidth(30)

	m := model{
		table:    tb,
		textarea: ta,
		viewport: vp,
		spinner:  sp,
		response: SEResponse{},
		state:    WaitingForInput,
		err:      nil,
		mouse:    true,
	}

	m.SetTableHeaders()
	return m
}

func (m model) SetTableHeaders() {
	columns := []table.Column{
		{
			Title: "Score",
			Width: int(0.1 * float32(m.table.Width())),
		},
		{
			Title: "Title",
			Width: int(0.6 * float32(m.table.Width())),
		},
		{
			Title: "Views",
			Width: int(0.1 * float32(m.table.Width())),
		},
		{
			Title: "Last Activity",
			Width: int(0.2 * float32(m.table.Width())),
		},
	}

	m.table.SetColumns(columns)
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		taCmd tea.Cmd
		vpCmd tea.Cmd
		spCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.table, taCmd = m.table.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	m.spinner, spCmd = m.spinner.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.table.SetHeight(msg.Height - 10)
		m.table.SetWidth(msg.Width - 6)
		m.SetTableHeaders()

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
			if m.state == WaitingForInput {
				go func() {
					question := m.textarea.Value()
					m.textarea.SetValue("")
					resp := Search(question, "", "", "", "") //TODO: Add the other params here

					tui.Send(resp)
				}()
				m.state = WaitingForResponse
				return m, tea.Batch(tiCmd, taCmd, vpCmd, spinner.Tick)
			}
		}

	case SEResponse:
		if len(msg.Items) == 0 {
			m.state = WaitingForInput
			m.textarea.Blur()
			m.textarea.SetValue("")
			return m, tea.Batch(tiCmd, taCmd, vpCmd, spCmd, getLogCmd("No results found", Warning))
		}

		m.response = msg
		m.state = DisplayingAllQuestions
		m.table.SetRows(m.response.ToRows())
		m.textarea.Blur()
		m.table.Focus()

		return m, nil

	case logMsg:
		if msg.Msg == "" {
			//TODO:Remove the overlay component here
			return m, nil
		}
		switch msg.Type {
		//TODO: Overlay component here in bottom right area
		case Info:
		case Warning:
		case Error:
		}

		go func() {
			time.Sleep(3 * time.Second)
			tui.Send(getLogCmd("", Info))
		}()

		return m, nil

	case errMsg:
		m.err = msg
		return m, nil

	}

	return m, tea.Batch(tiCmd, taCmd, vpCmd, spCmd)
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	} else if m.state == WaitingForInput {
		return m.textarea.View()
	} else if m.state == WaitingForResponse {
		return m.spinner.View()
	} else if m.state == DisplayingAllQuestions {
		return m.table.View()
	} else if m.state == DisplayingQuestion || m.state == DisplayingAllComments || m.state == DisplayingHelpScreen {
		return m.viewport.View()
	}

	return ""
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
		return logMsg{Msg: msg, Type: logType}
	}
}
