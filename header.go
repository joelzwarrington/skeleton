package skeleton

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// header is a helper for rendering the header of the terminal.
type header struct {
	// termReady is control terminal is ready or not, it responsible for the terminal size
	termReady bool

	// currentTab is hold the current tab index
	currentTab int

	// viewport is hold the viewport, it is responsible for the terminal size
	viewport *viewport.Model

	// keyMap responsible for the key bindings
	keyMap *keyMap

	// headers are hold the headers of the terminal
	headers []commonHeader

	// properties are hold the properties of the header
	properties *headerProperties

	// titleLength is hold the length of the title
	titleLength int

	updater *Updater

	// lockedTabs holds the keys of individually locked tabs
	lockedTabs map[string]bool
}

// newHeader returns a new header.
func newHeader() *header {
	return &header{
		properties: defaultHeaderProperties(),
		viewport:   newTerminalViewport(),
		currentTab: 0,
		keyMap:     newKeyMap(),
		updater:    NewUpdater(),
		lockedTabs: make(map[string]bool),
	}
}

// headerProperties are hold the properties of the header.
type headerProperties struct {
	borderColor        string
	leftTabPadding     int
	rightTabPadding    int
	titleStyleActive   lipgloss.Style
	titleStyleInactive lipgloss.Style
	titleStyleDisabled lipgloss.Style
}

// defaultHeaderProperties returns the default properties of the header.
func defaultHeaderProperties() *headerProperties {
	borderColor := "39"
	leftPadding := 2
	rightPadding := 2
	return &headerProperties{
		borderColor:     borderColor,
		leftTabPadding:  leftPadding,
		rightTabPadding: rightPadding,
		titleStyleActive: func() lipgloss.Style {
			b := lipgloss.DoubleBorder()
			b.Right = "├"
			b.Left = "┤"
			return lipgloss.NewStyle().BorderStyle(b).
				PaddingLeft(leftPadding).PaddingRight(rightPadding).
				BorderForeground(lipgloss.Color("205"))
		}(),
		titleStyleInactive: func() lipgloss.Style {
			b := lipgloss.RoundedBorder()
			b.Right = "├"
			b.Left = "┤"
			return lipgloss.NewStyle().BorderStyle(b).
				PaddingLeft(leftPadding).PaddingRight(rightPadding).
				BorderForeground(lipgloss.Color("255"))
		}(),
		titleStyleDisabled: func() lipgloss.Style {
			b := lipgloss.RoundedBorder()
			b.Right = "├"
			b.Left = "┤"
			return lipgloss.NewStyle().BorderStyle(b).
				PaddingLeft(leftPadding).PaddingRight(rightPadding).
				BorderForeground(lipgloss.Color("240")).Foreground(lipgloss.Color("240"))
		}(),
	}
}

// commonHeader is hold the header required fields.
type commonHeader struct {
	key   string
	title string
}

func (h *header) Init() tea.Cmd {
	return nil
}

func (h *header) Update(msg tea.Msg) (*header, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !h.termReady {
			if msg.Width > 0 && msg.Height > 0 {
				h.termReady = true
			}
		}
		h.viewport.Width = msg.Width
		h.viewport.Height = msg.Height

		h.calculateTitleLength()

		cmds = append(cmds, h.calculateTitleLength())
	}

	return h, tea.Batch(cmds...)
}

type HeaderSizeMsg struct {
	NotEnoughToHandleHeaders bool
}

// calculateTitleLength calculates the length of the title.
func (h *header) calculateTitleLength() tea.Cmd {
	var titleLen int
	for _, hdr := range h.headers {
		titleLen += len([]rune(hdr.title))
		titleLen += h.properties.leftTabPadding + h.properties.rightTabPadding
		titleLen += 2 // for the border between titles
	}

	requiredLineCountForLine := h.viewport.Width - (titleLen + 2)

	h.titleLength = titleLen
	if requiredLineCountForLine < 0 {
		return func() tea.Msg {
			return HeaderSizeMsg{NotEnoughToHandleHeaders: false}
		}
	} else {
		return func() tea.Msg {
			return HeaderSizeMsg{NotEnoughToHandleHeaders: true}
		}
	}
}

// View renders the header.
func (h *header) View() string {
	if !h.termReady {
		return "setting up terminal..."
	}

	requiredLineCount := h.viewport.Width - (h.titleLength + 2)

	if requiredLineCount < 0 {
		return ""
	}

	line := strings.Repeat("─", requiredLineCount)
	line = lipgloss.NewStyle().Foreground(lipgloss.Color(h.properties.borderColor)).Render(line)

	var renderedTitles []string
	renderedTitles = append(renderedTitles, "")
	for i, hdr := range h.headers {
		if i == h.currentTab {
			renderedTitles = append(renderedTitles, h.properties.titleStyleActive.Render(hdr.title))
		} else {
			if h.GetLockTabs() || h.IsTabLocked(hdr.key) {
				renderedTitles = append(renderedTitles, h.properties.titleStyleDisabled.Render(hdr.title))
			} else {
				renderedTitles = append(renderedTitles, h.properties.titleStyleInactive.Render(hdr.title))
			}
		}
	}

	leftCorner := lipgloss.JoinVertical(lipgloss.Top, "╭", "│")
	rightCorner := lipgloss.JoinVertical(lipgloss.Top, "╮", "│")
	leftCorner = lipgloss.NewStyle().Foreground(lipgloss.Color(h.properties.borderColor)).Render(leftCorner)
	rightCorner = lipgloss.NewStyle().Foreground(lipgloss.Color(h.properties.borderColor)).Render(rightCorner)

	return lipgloss.JoinHorizontal(lipgloss.Bottom, leftCorner, lipgloss.JoinHorizontal(lipgloss.Center, append(renderedTitles, line)...), rightCorner)
}

// SetLeftPadding sets the left padding of the header.
func (h *header) SetLeftPadding(padding int) {
	h.properties.leftTabPadding = padding
	h.properties.titleStyleActive = h.properties.titleStyleActive.PaddingLeft(padding)
	h.properties.titleStyleInactive = h.properties.titleStyleInactive.PaddingLeft(padding)
	h.properties.titleStyleDisabled = h.properties.titleStyleDisabled.PaddingLeft(padding)

	h.calculateTitleLength()
}

// SetRightPadding sets the right padding of the header.
func (h *header) SetRightPadding(padding int) {
	h.properties.rightTabPadding = padding
	h.properties.titleStyleActive = h.properties.titleStyleActive.PaddingRight(padding)
	h.properties.titleStyleInactive = h.properties.titleStyleInactive.PaddingRight(padding)
	h.properties.titleStyleDisabled = h.properties.titleStyleDisabled.PaddingRight(padding)

	h.calculateTitleLength()
}

// SetInactiveTabTextColor sets the idle tab color of the header.
func (h *header) SetInactiveTabTextColor(color string) {
	h.properties.titleStyleInactive = h.properties.titleStyleInactive.Foreground(lipgloss.Color(color))
}

// SetInactiveTabBorderColor sets the idle tab border color of the header.
func (h *header) SetInactiveTabBorderColor(color string) {
	h.properties.titleStyleInactive = h.properties.titleStyleInactive.BorderForeground(lipgloss.Color(color))
}

// SetActiveTabTextColor sets the active tab color of the header.
func (h *header) SetActiveTabTextColor(color string) {
	h.properties.titleStyleActive = h.properties.titleStyleActive.Foreground(lipgloss.Color(color))
}

// SetActiveTabBorderColor sets the active tab border color of the header.
func (h *header) SetActiveTabBorderColor(color string) {
	h.properties.titleStyleActive = h.properties.titleStyleActive.BorderForeground(lipgloss.Color(color))
}

// SetBorderColor sets the border color of the header.
func (h *header) SetBorderColor(color string) {
	h.properties.borderColor = color
}

// SetCurrentTab sets the current tab index.
func (h *header) SetCurrentTab(tab int) {
	h.currentTab = tab
}

// SetLockTabs sets the lock tabs status.
func (h *header) SetLockTabs(lock bool) {
	if lock {
		for _, header := range h.headers {
			h.LockTab(header.key)
		}
	} else {
		h.lockedTabs = make(map[string]bool)
	}
	h.updater.Update()
}

// GetLockTabs returns the lock tabs status.
func (h *header) GetLockTabs() bool {
	for _, header := range h.headers {
		if !h.IsTabLocked(header.key) {
			return false
		}
	}
	return true
}

// GetCurrentTab returns the current tab index.
func (h *header) GetCurrentTab() int {
	return h.currentTab
}

// AddCommonHeader adds a new header to the header.
func (h *header) AddCommonHeader(key string, title string) {
	h.headers = append(h.headers, commonHeader{
		key:   key,
		title: title,
	})
	h.calculateTitleLength()
	h.updater.Update()
}

// UpdateCommonHeader updates the header by the given key.
func (h *header) UpdateCommonHeader(key string, title string) {
	for i, header := range h.headers {
		if header.key == key {
			h.headers[i].title = title
		}
	}
	h.calculateTitleLength()
	h.updater.Update()
}

// DeleteCommonHeader deletes the header by the given key.
func (h *header) DeleteCommonHeader(key string) {
	for i, header := range h.headers {
		if header.key == key {
			h.headers = append(h.headers[:i], h.headers[i+1:]...)
		}
	}
	h.calculateTitleLength()
	h.updater.Update()
}

// IsTabLocked checks if a specific tab is locked
func (h *header) IsTabLocked(key string) bool {
	return h.lockedTabs[key]
}

// LockTab locks a specific tab by its key
func (h *header) LockTab(key string) {
	h.lockedTabs[key] = true
	h.updater.Update()
}

// UnlockTab unlocks a specific tab by its key
func (h *header) UnlockTab(key string) {
	delete(h.lockedTabs, key)
	h.updater.Update()
}
