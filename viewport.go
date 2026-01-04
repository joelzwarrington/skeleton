package skeleton

import (
	"sync"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

// --------------------------------------------

var (
	onceViewport sync.Once
	vp           *viewport.Model
)

func newTerminalViewport() *viewport.Model {
	onceViewport.Do(func() {
		vp = &viewport.Model{Width: 80, Height: 24} // Question: Is it best to use 80x24 as default?
	})
	return vp
}

// --------------------------------------------

// GetTerminalViewport returns the viewport.
func (s *Skeleton) GetTerminalViewport() *viewport.Model {
	return vp
}

// SetTerminalViewportWidth sets the width of the viewport.
func (s *Skeleton) SetTerminalViewportWidth(width int) {
	vp.Width = width
}

// SetTerminalViewportHeight sets the height of the viewport.
func (s *Skeleton) SetTerminalViewportHeight(height int) {
	vp.Height = height
}

// GetTerminalWidth returns the width of the terminal.
func (s *Skeleton) GetTerminalWidth() int {
	return vp.Width
}

// GetTerminalHeight returns the height of the terminal.
func (s *Skeleton) GetTerminalHeight() int {
	return vp.Height
}

// GetContentWidth returns the available width for content (terminal width minus borders).
func (s *Skeleton) GetContentWidth() int {
	return vp.Width - 2
}

// GetContentHeight returns the available height for content (terminal height minus header and widgets).
func (s *Skeleton) GetContentHeight() int {
	headerHeight := lipgloss.Height(s.header.View())
	footerHeight := lipgloss.Height(s.widget.View())
	return vp.Height - headerHeight - footerHeight
}
