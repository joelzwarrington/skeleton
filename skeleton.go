package skeleton

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

// Skeleton is a helper for rendering the Skeleton of the terminal.
type Skeleton struct {
	// termReady is control terminal is ready or not, it responsible for the terminal size
	termReady bool

	// termSizeNotEnoughToHandleHeaders is control terminal size is enough to handle headers
	termSizeNotEnoughToHandleHeaders bool

	// termSizeNotEnoughToHandleWidgets is control terminal size is enough to handle widgets
	termSizeNotEnoughToHandleWidgets bool

	// lockTabs is control the tabs (headers) are locked or not
	lockTabs bool

	// currentTab is hold the current tab index
	currentTab int

	// viewport is hold the viewport, it responsible for the terminal size
	viewport *viewport.Model

	// header is hold the header
	header *header

	// widget is hold the widget
	widget *widget

	// KeyMap responsible for the key bindings
	KeyMap *keyMap

	// pages are hold the pages
	pages []tea.Model

	// properties are hold the properties of the Skeleton
	properties *skeletonProperties

	updater *Updater
}

// NewSkeleton returns a new Skeleton.
func NewSkeleton() *Skeleton {
	return &Skeleton{
		properties: defaultSkeletonProperties(),
		viewport:   newTerminalViewport(),
		header:     newHeader(),
		widget:     newWidget(),
		KeyMap:     newKeyMap(),
		updater:    NewUpdater(),
	}
}

// skeletonProperties are hold the properties of the Skeleton.
type skeletonProperties struct {
	borderColor  string
	pagePosition lipgloss.Position
}

// defaultSkeletonProperties returns the default properties of the Skeleton.
func defaultSkeletonProperties() *skeletonProperties {
	return &skeletonProperties{
		borderColor:  "39",
		pagePosition: lipgloss.Center,
	}
}

// SetBorderColor sets the border color of the Skeleton.
func (s *Skeleton) SetBorderColor(color string) *Skeleton {
	s.header.SetBorderColor(color)
	s.widget.SetBorderColor(color)
	s.properties.borderColor = color
	s.updater.Update()
	return s
}

// GetBorderColor returns the border color of the Skeleton.
func (s *Skeleton) GetBorderColor() string {
	return s.properties.borderColor
}

// GetWidgetBorderColor returns the border color of the Widget.
func (s *Skeleton) GetWidgetBorderColor() string {
	return s.widget.GetBorderColor()
}

// SetPagePosition sets the position of the page.
func (s *Skeleton) SetPagePosition(position lipgloss.Position) *Skeleton {
	s.properties.pagePosition = position
	s.updater.Update()
	return s
}

// GetPagePosition returns the position of the page.
func (s *Skeleton) GetPagePosition() lipgloss.Position {
	return s.properties.pagePosition
}

// SetInactiveTabTextColor sets the idle tab color of the Skeleton.
func (s *Skeleton) SetInactiveTabTextColor(color string) *Skeleton {
	s.header.SetInactiveTabTextColor(color)
	s.updater.Update()
	return s
}

// SetInactiveTabBorderColor sets the idle tab border color of the Skeleton.
func (s *Skeleton) SetInactiveTabBorderColor(color string) *Skeleton {
	s.header.SetInactiveTabBorderColor(color)
	s.updater.Update()
	return s
}

// SetActiveTabTextColor sets the active tab color of the Skeleton.
func (s *Skeleton) SetActiveTabTextColor(color string) *Skeleton {
	s.header.SetActiveTabTextColor(color)
	s.updater.Update()
	return s
}

// SetActiveTabBorderColor sets the active tab border color of the Skeleton.
func (s *Skeleton) SetActiveTabBorderColor(color string) *Skeleton {
	s.header.SetActiveTabBorderColor(color)
	s.updater.Update()
	return s
}

// SetWidgetBorderColor sets the border color of the Widget.
func (s *Skeleton) SetWidgetBorderColor(color string) *Skeleton {
	s.widget.SetWidgetBorderColor(color)
	s.updater.Update()
	return s
}

// SetTabLeftPadding sets the left padding of the Skeleton.
func (s *Skeleton) SetTabLeftPadding(padding int) *Skeleton {
	s.header.SetLeftPadding(padding)
	s.updater.Update()
	return s
}

// SetTabRightPadding sets the right padding of the Skeleton.
func (s *Skeleton) SetTabRightPadding(padding int) *Skeleton {
	s.header.SetRightPadding(padding)
	s.updater.Update()
	return s
}

// SetWidgetLeftPadding sets the left padding of the Skeleton.
func (s *Skeleton) SetWidgetLeftPadding(padding int) *Skeleton {
	s.widget.SetLeftPadding(padding)
	s.updater.Update()
	return s
}

// SetWidgetRightPadding sets the right padding of the Skeleton.
func (s *Skeleton) SetWidgetRightPadding(padding int) *Skeleton {
	s.widget.SetRightPadding(padding)
	s.updater.Update()
	return s
}

// LockTabs locks the tabs (headers). It prevents switching tabs. It is useful when you want to prevent switching tabs.
func (s *Skeleton) LockTabs() *Skeleton {
	s.header.SetLockTabs(true)
	s.lockTabs = true
	s.updater.Update()
	return s
}

// UnlockTabs unlocks the tabs (headers). It allows switching tabs. It is useful when you want to allow switching tabs.
func (s *Skeleton) UnlockTabs() *Skeleton {
	s.header.SetLockTabs(false)
	s.lockTabs = false
	s.updater.Update()
	return s
}

// IsTabsLocked returns the tabs (headers) are locked or not.
func (s *Skeleton) IsTabsLocked() bool {
	return s.lockTabs
}

// AddPage adds a new page to the Skeleton.
func (s *Skeleton) AddPage(key string, title string, page tea.Model) *Skeleton {
	// do not add if key already exists
	for _, hdr := range s.header.headers {
		if hdr.key == key {
			return s
		}
	}

	s.header.AddCommonHeader(key, title)
	s.pages = append(s.pages, page)

	s.updater.Update()
	return s
}

// UpdatePageTitle updates the title of the page by the given key.
func (s *Skeleton) UpdatePageTitle(key string, title string) *Skeleton {
	s.header.UpdateCommonHeader(key, title)
	s.updater.Update()
	return s
}

// DeletePage deletes the page by the given key.
func (s *Skeleton) DeletePage(key string) *Skeleton {
	s.deletePage(key)
	s.updater.Update()
	return s
}

// deletePage deletes the page by the given key.
func (s *Skeleton) deletePage(key string) {
	if len(s.pages) == 1 {
		// skeleton should have at least one page
		return
	}

	// if active tab is about deleting tab, switch to the first tab
	if s.GetActivePage() == key {
		s.currentTab = 0
		s.header.SetCurrentTab(0)
	}

	var pages []tea.Model
	for i := range s.pages {
		if s.header.headers[i].key != key {
			pages = append(pages, s.pages[i])
		}
	}
	s.header.DeleteCommonHeader(key)
	s.pages = pages
}

// AddWidget adds a new widget to the Skeleton.
func (s *Skeleton) AddWidget(key string, value string) *Skeleton {
	s.widget.addNewWidget(key, value)
	s.updater.Update()
	return s
}

// UpdateWidgetValue updates the Value content by the given key.
// Adds the widget if it doesn't exist.
func (s *Skeleton) UpdateWidgetValue(key string, value string) *Skeleton {
	go func() {
		// if widget not exists, add it
		if s.widget.GetWidget(key) == nil {
			s.AddWidget(key, value)
		}

		s.widget.updateWidgetContent(key, value)

		s.updater.Update()
	}()

	return s
}

// DeleteWidget deletes the Value by the given key.
func (s *Skeleton) DeleteWidget(key string) *Skeleton {
	s.widget.deleteWidget(key)
	s.updater.Update()
	return s
}

// DeleteAllWidgets deletes all the widgets.
func (s *Skeleton) DeleteAllWidgets() *Skeleton {
	s.widget.DeleteAllWidgets()
	s.updater.Update()
	return s
}

// SetActivePage sets the active page by the given key.
func (s *Skeleton) SetActivePage(key string) *Skeleton {
	for i, header := range s.header.headers {
		if header.key == key {
			s.currentTab = i
			s.header.SetCurrentTab(i)
			s.updater.Update()
			break
		}
	}
	return s
}

// GetActivePage returns the active page key.
func (s *Skeleton) GetActivePage() string {
	return s.header.headers[s.currentTab].key
}

// IAMActivePage is a message to trigger the update of the active page.
type IAMActivePage struct{}

// IAMActivePageCmd returns the IAMActivePage command.
func (s *Skeleton) IAMActivePageCmd() tea.Cmd {
	return func() tea.Msg {
		return IAMActivePage{}
	}
}

func (s *Skeleton) switchPage(cmds []tea.Cmd, position string) []tea.Cmd {
	switch position {
	case "left":
		if !s.IsTabsLocked() {
			s.currentTab = max(s.currentTab-1, 0)
			cmds = append(cmds, s.IAMActivePageCmd())
		}
	case "right":
		if !s.IsTabsLocked() {
			s.currentTab = min(s.currentTab+1, len(s.pages)-1)
			cmds = append(cmds, s.IAMActivePageCmd())
		}
	}

	return cmds
}

func (s *Skeleton) updateSkeleton(msg tea.Msg, cmd tea.Cmd, cmds []tea.Cmd) []tea.Cmd {
	s.header, cmd = s.header.Update(msg)
	cmds = append(cmds, cmd)

	s.widget, cmd = s.widget.Update(msg)
	cmds = append(cmds, cmd)

	s.pages[s.currentTab], cmd = s.pages[s.currentTab].Update(msg)
	cmds = append(cmds, cmd)

	return cmds
}

func (s *Skeleton) Init() tea.Cmd {
	if len(s.pages) == 0 {
		panic("skeleton: no pages added, please add at least one page")
	}

	return tea.Batch(s.updater.Listen(), s.header.Init(), s.widget.Init())
}

func (s *Skeleton) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	s.currentTab = s.header.GetCurrentTab()

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !s.termReady {
			if msg.Width > 0 && msg.Height > 0 {
				s.termReady = true
			}
		}
		s.viewport.Width = msg.Width
		s.viewport.Height = msg.Height

		cmds = s.updateSkeleton(msg, cmd, cmds)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.KeyMap.Quit):
			return s, tea.Quit
		case key.Matches(msg, s.KeyMap.SwitchTabLeft):
			cmds = s.switchPage(cmds, "left")
		case key.Matches(msg, s.KeyMap.SwitchTabRight):
			cmds = s.switchPage(cmds, "right")
		}
		cmds = s.updateSkeleton(msg, cmd, cmds)
	case UpdateMsg:
		// do nothing, just to trigger the update
		cmds = append(cmds, s.updater.Listen()) // listen to the update channel
	case HeaderSizeMsg:
		s.termSizeNotEnoughToHandleHeaders = msg.NotEnoughToHandleHeaders
	case WidgetSizeMsg:
		s.termSizeNotEnoughToHandleWidgets = msg.NotEnoughToHandleWidgets

	default:
		cmds = s.updateSkeleton(msg, cmd, cmds)
	}

	return s, tea.Batch(cmds...)
}

func (s *Skeleton) View() string {
	if !s.termReady {
		return "setting up terminal..."
	}
	if !s.termSizeNotEnoughToHandleHeaders {
		return "terminal size is not enough to show headers"
	}
	if !s.termSizeNotEnoughToHandleWidgets {
		return "terminal size is not enough to show widgets"
	}

	base := lipgloss.NewStyle().
		BorderForeground(lipgloss.Color(s.properties.borderColor)).
		Align(s.properties.pagePosition).
		Border(lipgloss.RoundedBorder()).
		BorderTop(false).BorderBottom(false).
		Width(s.viewport.Width - 2)

	body := s.pages[s.currentTab].View()

	bodyHeight := s.viewport.Height - 5 // for header height and Value height
	if len(s.widget.widgets) > 0 {
		bodyHeight -= 1
	}
	if lipgloss.Height(body) < bodyHeight {
		body += strings.Repeat("\n", bodyHeight-lipgloss.Height(body))
	}

	return lipgloss.JoinVertical(lipgloss.Top, s.header.View(), base.Render(body), s.widget.View())
}
