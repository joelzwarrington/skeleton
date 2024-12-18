package skeleton

import (
	tea "github.com/charmbracelet/bubbletea"
	"sync"
)

type Updater struct {
	rcv       chan any
	listening bool
	mu        sync.Mutex
}

var (
	updaterInstance *Updater
	onceUpdater     sync.Once
)

func NewUpdater() *Updater {
	onceUpdater.Do(func() {
		updaterInstance = &Updater{
			rcv: make(chan any, 256), // 256 is a reasonable buffer size for most cases, but it depends on your application's needs.
		}
	})

	return updaterInstance
}

type UpdateMsg struct{}

var UpdateMsgInstance UpdateMsg

func (u *Updater) Listen() tea.Cmd {
	u.mu.Lock()
	defer u.mu.Unlock()

	// Ensure only one listener is active
	if u.listening {
		return nil
	}

	u.listening = true

	return func() tea.Msg {
		// This function will block until a message is received
		msg := <-u.rcv
		u.mu.Lock()
		u.listening = false
		u.mu.Unlock()
		return msg
	}
}

func (u *Updater) Update() {
	// Add non-blocking send
	select {
	case u.rcv <- UpdateMsgInstance:
		// Successfully sent
	default:
		// Channel is full, skip update
	}
}

func (u *Updater) UpdateWithMsg(msg any) {
	// Add non-blocking send
	select {
	case u.rcv <- msg:
		// Successfully sent
	default:
		// Channel is full, skip update
	}
}
