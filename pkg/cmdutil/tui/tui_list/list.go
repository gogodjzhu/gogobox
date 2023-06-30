package tui_list

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	list_fancy "gogobox/pkg/cmdutil/tui/tui_list_fancy"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Option struct {
	title, desc string
}

func NewOption(title, desc string) Option {
	return Option{
		title: title,
		desc:  desc,
	}
}

type CallbackFunc struct {
	Keys             []string
	ShortDescription string
	FullDescription  string
	Callback         func(option Option) []Option
}

func (i Option) Title() string       { return i.title }
func (i Option) Description() string { return i.desc }
func (i Option) FilterValue() string { return i.title }

type model struct {
	list             list.Model
	key2callbackFunc map[string]func(option Option) []Option
}

func NewModel(title string, options []Option, callback []CallbackFunc) tea.Model {
	var items []list.Item
	for _, option := range options {
		items = append(items, option)
	}
	l := list.New(items, list_fancy.NewDefaultDelegate(), 0, 0)
	l.Title = title

	k2cMap := make(map[string]func(option Option) []Option)
	shortHelpKeys := make([]key.Binding, 0)
	fullHelpKeys := make([]key.Binding, 0)
	for _, cb := range callback {
		k2cMap[cb.Keys[0]] = cb.Callback
		if cb.ShortDescription != "" {
			shortHelpKeys = append(shortHelpKeys, key.NewBinding(
				key.WithKeys(cb.Keys[0]),
				key.WithHelp(cb.Keys[0], cb.ShortDescription),
			))
		}
		if cb.FullDescription != "" {
			fullHelpKeys = append(fullHelpKeys, key.NewBinding(
				key.WithKeys(cb.Keys...),
				key.WithHelp(cb.Keys[0], cb.FullDescription),
			))
		}
	}
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return shortHelpKeys
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return fullHelpKeys
	}
	return model{
		list:             l,
		key2callbackFunc: k2cMap,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if callbackFunc, ok := m.key2callbackFunc[msg.String()]; ok && m.list.SelectedItem() != nil {
			if options := callbackFunc(m.list.SelectedItem().(Option)); options != nil {
				var items []list.Item
				for _, option := range options {
					items = append(items, option)
				}
				m.list.SetItems(items)
			}
		}
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}
