package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/termkit/skeleton"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// -----------------------------------------------------------------------------
// CPU Model
type cpuModel struct {
	skeleton    *skeleton.Skeleton
	usage       float64
	lastUpdate  time.Time
	color       string // tab color
	borderColor string // border color
}

func newCPUModel(s *skeleton.Skeleton) *cpuModel {
	return &cpuModel{
		skeleton:    s,
		usage:       0,
		lastUpdate:  time.Now(),
		color:       "39", // bright blue
		borderColor: "27", // darker blue
	}
}

func (m *cpuModel) Init() tea.Cmd {
	return nil
}

func (m *cpuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case skeleton.UpdateMsg:
		if time.Since(m.lastUpdate) >= time.Second {
			if percent, err := cpu.Percent(0, false); err == nil && len(percent) > 0 {
				m.usage = percent[0]
			}
			m.lastUpdate = time.Now()
		}
	case skeleton.IAMActivePage:
		m.skeleton.SetActiveTabBorderColor(m.color)
		m.skeleton.SetBorderColor(m.borderColor)
	}
	return m, nil
}

func (m *cpuModel) View() string {
	mainStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Padding(1)

	bar, valueColor := createProgressBar(progressBarOptions{
		value:      m.usage,
		width:      50,
		barColor:   "205",
		warningAt:  60,
		criticalAt: 80,
	})

	// CPU usage value
	usageValue := lipgloss.NewStyle().
		Foreground(lipgloss.Color(valueColor)).
		Bold(true).
		Render(fmt.Sprintf("%.1f%%", m.usage))

	info := []string{
		"CPU Usage",
		"",
		bar,
		"",
		fmt.Sprintf("Current Usage: %s", usageValue),
		"",
		"Press Ctrl+Left/Right to switch tabs",
	}

	return mainStyle.Render(strings.Join(info, "\n"))
}

// -----------------------------------------------------------------------------
// Memory Model
type memoryModel struct {
	skeleton    *skeleton.Skeleton
	used        uint64
	total       uint64
	lastUpdate  time.Time
	color       string // tab color
	borderColor string // border color
}

func newMemoryModel(s *skeleton.Skeleton) *memoryModel {
	return &memoryModel{
		skeleton:    s,
		lastUpdate:  time.Now(),
		color:       "162", // bright purple
		borderColor: "126", // darker purple
	}
}

func (m *memoryModel) Init() tea.Cmd {
	return nil
}

func (m *memoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case skeleton.UpdateMsg:
		if time.Since(m.lastUpdate) >= time.Second {
			if v, err := mem.VirtualMemory(); err == nil {
				m.used = v.Used
				m.total = v.Total
			}
			m.lastUpdate = time.Now()
		}
	case skeleton.IAMActivePage:
		m.skeleton.SetActiveTabBorderColor(m.color)
		m.skeleton.SetBorderColor(m.borderColor)
	}
	return m, nil
}

func (m *memoryModel) View() string {
	mainStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Padding(1)

	// Calculate memory values
	usedGB := float64(m.used) / 1024 / 1024 / 1024
	totalGB := float64(m.total) / 1024 / 1024 / 1024
	usagePercent := float64(0)

	if totalGB > 0 {
		usagePercent = (usedGB / totalGB) * 100
	}

	bar, valueColor := createProgressBar(progressBarOptions{
		value:      usagePercent,
		width:      50,
		barColor:   "205",
		warningAt:  60,
		criticalAt: 80,
	})

	// Memory values
	memoryStats := lipgloss.NewStyle().
		Foreground(lipgloss.Color(valueColor)).
		Bold(true).
		Render(fmt.Sprintf("%.1f GB / %.1f GB (%.1f%%)", usedGB, totalGB, usagePercent))

	info := []string{
		"Memory Usage",
		"",
		bar,
		"",
		fmt.Sprintf("Used Memory: %s", memoryStats),
		"",
		"Press Ctrl+Left/Right to switch tabs",
	}

	return mainStyle.Render(strings.Join(info, "\n"))
}

// -----------------------------------------------------------------------------
// Disk Model
type diskModel struct {
	skeleton    *skeleton.Skeleton
	used        uint64
	total       uint64
	lastUpdate  time.Time
	color       string // tab color
	borderColor string // border color
}

func newDiskModel(s *skeleton.Skeleton) *diskModel {
	return &diskModel{
		skeleton:    s,
		lastUpdate:  time.Now(),
		color:       "136", // bright gold
		borderColor: "94",  // darker gold
	}
}

func (m *diskModel) Init() tea.Cmd {
	return nil
}

func (m *diskModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case skeleton.UpdateMsg:
		if time.Since(m.lastUpdate) >= time.Second {
			if usage, err := disk.Usage("/"); err == nil {
				m.used = usage.Used
				m.total = usage.Total
			}
			m.lastUpdate = time.Now()
		}
	case skeleton.IAMActivePage:
		m.skeleton.SetActiveTabBorderColor(m.color)
		m.skeleton.SetBorderColor(m.borderColor)
	}
	return m, nil
}

func (m *diskModel) View() string {
	mainStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Padding(1)

	// Calculate disk values
	usedGB := float64(m.used) / 1024 / 1024 / 1024
	totalGB := float64(m.total) / 1024 / 1024 / 1024
	usagePercent := float64(0)

	if totalGB > 0 {
		usagePercent = (usedGB / totalGB) * 100
	}

	bar, valueColor := createProgressBar(progressBarOptions{
		value:      usagePercent,
		width:      50,
		barColor:   "205",
		warningAt:  70,
		criticalAt: 90,
	})

	// Disk values
	diskStats := lipgloss.NewStyle().
		Foreground(lipgloss.Color(valueColor)).
		Bold(true).
		Render(fmt.Sprintf("%.1f GB / %.1f GB (%.1f%%)", usedGB, totalGB, usagePercent))

	info := []string{
		"Disk Usage",
		"",
		bar,
		"",
		fmt.Sprintf("Used Space: %s", diskStats),
		fmt.Sprintf("Free Space: %.1f GB", totalGB-usedGB),
		"",
		"Press Ctrl+Left/Right to switch tabs",
	}

	return mainStyle.Render(strings.Join(info, "\n"))
}

// -----------------------------------------------------------------------------
// Helper Functions

// progressBar creates a styled progress bar with given percentage
type progressBarOptions struct {
	value      float64 // percentage value (0-100)
	width      int     // total width of the bar
	barColor   string  // color of the progress bar
	warningAt  float64 // warning threshold (orange)
	criticalAt float64 // critical threshold (red)
}

func createProgressBar(opts progressBarOptions) (string, string) {
	// Ensure value is within bounds
	if opts.value < 0 {
		opts.value = 0
	} else if opts.value > 100 {
		opts.value = 100
	}

	// Calculate used width
	usedWidth := int(float64(opts.width) * opts.value / 100)
	if usedWidth < 0 {
		usedWidth = 0
	} else if usedWidth > opts.width {
		usedWidth = opts.width
	}

	// Create progress bar
	barStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(opts.barColor))
	bar := "█" + strings.Repeat("█", usedWidth) +
		strings.Repeat("░", opts.width-usedWidth)

	// Select color based on thresholds
	valueColor := "39" // blue for normal
	if opts.value > opts.criticalAt {
		valueColor = "196" // red for critical
	} else if opts.value > opts.warningAt {
		valueColor = "208" // orange for warning
	}

	return barStyle.Render(bar), valueColor
}

// -----------------------------------------------------------------------------
// Main Program
func main() {
	s := skeleton.NewSkeleton()

	// Add tabs (pages)
	s.AddPage("cpu", "CPU", newCPUModel(s))
	s.AddPage("memory", "Memory", newMemoryModel(s))
	s.AddPage("disk", "Disk", newDiskModel(s))

	s.SetActiveTabBorderColor("39") // dark blue

	// Add widgets
	s.AddWidget("app", "System Monitor")
	s.AddWidget("time", time.Now().Format("15:04:05"))

	// Update system stats every second
	go func() {
		for {
			time.Sleep(time.Second)
			s.TriggerUpdate()
			s.UpdateWidgetValue("time", time.Now().Format("15:04:05"))
		}
	}()

	p := tea.NewProgram(s)
	if err := p.Start(); err != nil {
		panic(err)
	}
}
