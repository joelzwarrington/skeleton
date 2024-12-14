package skeleton

import (
	tea "github.com/charmbracelet/bubbletea"
	"sync"
)

type Updater struct {
	rcv chan struct{}
}

var (
	updaterInstance *Updater
	onceUpdater     sync.Once
)

func NewUpdater() *Updater {
	onceUpdater.Do(func() {
		updaterInstance = &Updater{
			rcv: make(chan struct{}, 1),
		}
	})

	return updaterInstance
}

type UpdateMsg struct{}

func (u *Updater) Listen() tea.Cmd {
	return func() tea.Msg {
		<-u.rcv
		return UpdateMsg{}
	}
}

func (u *Updater) Update() {
	// Add non-blocking send
	select {
	case u.rcv <- struct{}{}:
		// Successfully sent
	default:
		// Channel is full, skip update
	}
}
