package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mmcdole/gofeed"
	"github.com/termkit/skeleton"
)

// -----------------------------------------------------------------------------
// Types and Interfaces

// NewsItem represents a single news article with its metadata
type NewsItem struct {
	Title     string    // Title of the article
	Summary   string    // Full content of the article
	Link      string    // URL to the original article
	Published time.Time // Publication date
	IsRead    bool      // Read status
}

// -----------------------------------------------------------------------------
// News List Model

// categoryModel implements tea.Model and handles the news list view
type categoryModel struct {
	skeleton   *skeleton.Skeleton // Reference to main skeleton
	color      string             // Color theme for this tab
	feed       string             // RSS feed URL
	items      []NewsItem         // List of news items
	selected   int                // Currently selected item index
	lastUpdate time.Time          // Last feed update time
	lastError  error              // Last fetch error if any
	isLoading  bool               // Loading state indicator
}

// fetchMsg represents a completed fetch operation
type fetchMsg struct {
	items []NewsItem
	err   error
}

// -----------------------------------------------------------------------------
// News Detail Model

// newsDetailModel shows the full content of a news item
type newsDetailModel struct {
	skeleton *skeleton.Skeleton
	news     NewsItem
	color    string
}

// -----------------------------------------------------------------------------
// Model Implementations

// Category Model
func (m *categoryModel) Init() tea.Cmd {
	// Start fetching immediately when the model is initialized
	m.isLoading = true
	return m.fetchCmd()
}

func (m *categoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case fetchMsg:
		m.isLoading = false
		if msg.err != nil {
			m.lastError = msg.err
		} else {
			m.lastError = nil
			m.items = msg.items
			m.updateWidgets()
		}
		return m, nil

	case skeleton.IAMActivePage:
		m.skeleton.SetActiveTabBorderColor(m.color)

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.items)-1 {
				m.selected++
			}
		case "enter":
			if m.selected < len(m.items) {
				news := m.items[m.selected]
				m.items[m.selected].IsRead = true
				m.updateWidgets()

				// Create a unique key for the detail tab
				detailKey := fmt.Sprintf("detail-%d", time.Now().UnixNano())
				// Create a shorter title for the tab
				shortTitle := news.Title
				if len(shortTitle) > 20 {
					shortTitle = shortTitle[:17] + "..."
				}

				// Add new detail tab
				m.skeleton.AddPage(detailKey, shortTitle, newNewsDetailModel(m.skeleton, news))
				return m, nil
			}
		}
	}
	return m, nil
}

func (m *categoryModel) View() string {
	style := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Padding(1)

	var content []string

	if m.isLoading {
		content = append(content, lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Render("Loading news..."))
	} else if m.lastError != nil {
		content = append(content, lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Render(fmt.Sprintf("Error: %v", m.lastError)))
	}

	for i, item := range m.items {
		title := item.Title
		if item.IsRead {
			title = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Render("✓ " + title)
		} else {
			title = lipgloss.NewStyle().
				Bold(true).
				Render("• " + title)
		}

		if i == m.selected {
			title = lipgloss.NewStyle().
				Background(lipgloss.Color("214")).
				Foreground(lipgloss.Color("235")).
				Render("> " + title)
		} else {
			title = "  " + title
		}

		content = append(content, title)
	}

	if len(content) == 0 && !m.isLoading && m.lastError == nil {
		content = append(content, "No news available")
	}

	// Add key bindings help
	help := []string{
		"",
		lipgloss.NewStyle().Faint(true).Render("j/k or ↑/↓: Navigate"),
		lipgloss.NewStyle().Faint(true).Render("enter: Open article"),
		lipgloss.NewStyle().Faint(true).Render("ctrl+w: Close tab"),
		lipgloss.NewStyle().Faint(true).Render("ctrl+c: Quit"),
	}
	content = append(content, help...)

	return style.Render(lipgloss.JoinVertical(lipgloss.Left, content...))
}

// News Detail Model
func (m *newsDetailModel) Init() tea.Cmd {
	return nil
}

func (m *newsDetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case skeleton.IAMActivePage:
		m.skeleton.SetActiveTabBorderColor(m.color)

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlW {
			if key := m.skeleton.GetActivePage(); key != "" {
				m.skeleton.DeletePage(key)
			}
		}
	}
	return m, nil
}

func (m *newsDetailModel) View() string {
	style := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Padding(1)

	content := []string{
		lipgloss.NewStyle().Bold(true).Render(m.news.Title),
		"",
		lipgloss.NewStyle().Faint(true).Render(m.news.Published.Format("2006-01-02 15:04")),
		"",
		m.news.Summary,
		"",
		lipgloss.NewStyle().Underline(true).Render(m.news.Link),
	}

	return style.Render(lipgloss.JoinVertical(lipgloss.Left, content...))
}

// -----------------------------------------------------------------------------
// Helper Methods

func newCategoryModel(s *skeleton.Skeleton, feed string, color string) *categoryModel {
	return &categoryModel{
		skeleton:   s,
		color:      color,
		feed:       feed,
		items:      make([]NewsItem, 0),
		selected:   0,
		lastUpdate: time.Now(),
	}
}

func (m *categoryModel) fetchCmd() tea.Cmd {
	return func() tea.Msg {
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(m.feed)
		if err != nil {
			return fetchMsg{err: err}
		}

		items := make([]NewsItem, 0)
		for _, item := range feed.Items {
			pubDate := time.Now()
			if item.PublishedParsed != nil {
				pubDate = *item.PublishedParsed
			}

			content := item.Content
			if content == "" {
				content = item.Description // fallback to description if content is empty
			}

			items = append(items, NewsItem{
				Title:     item.Title,
				Summary:   content, // use content instead of description
				Link:      item.Link,
				Published: pubDate,
				IsRead:    false,
			})
		}

		return fetchMsg{items: items}
	}
}

func (m *categoryModel) updateWidgets() {
	m.skeleton.UpdateWidgetValue("count", fmt.Sprintf("News: %d | Unread: %d", len(m.items), m.countUnread()))
}

func (m *categoryModel) countUnread() int {
	unread := 0
	for _, item := range m.items {
		if !item.IsRead {
			unread++
		}
	}
	return unread
}

func newNewsDetailModel(s *skeleton.Skeleton, news NewsItem) *newsDetailModel {
	return &newsDetailModel{
		skeleton: s,
		news:     news,
		color:    "205", // Pink color for detail tabs
	}
}

// -----------------------------------------------------------------------------
// Main Program

func main() {
	s := skeleton.NewSkeleton()

	// Main tab with gruvbox theme
	s.AddPage("news", "News", newCategoryModel(s, "https://dev.to/feed", "142")) // Gruvbox green

	s.AddWidget("app", "News Reader")
	s.AddWidget("count", "Loading...")
	s.AddWidget("time", time.Now().Format("15:04:05"))

	s.SetActiveTabBorderColor("142") // Gruvbox green
	s.SetWidgetBorderColor("142")    // Gruvbox green
	s.SetBorderColor("214")          // Gruvbox orange

	// Update time every second
	go func() {
		for {
			time.Sleep(time.Second)
			s.UpdateWidgetValue("time", time.Now().Format("15:04:05"))
		}
	}()

	p := tea.NewProgram(s)
	if err := p.Start(); err != nil {
		panic(err)
	}
}
