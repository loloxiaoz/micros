package registry

import (
	"errors"
)

type mockWatcher struct {
	exit chan bool
}

func (m *mockWatcher) Next() (*Result, error) {
	// not implement so we just block until exit
	select {
	case <-m.exit:
		return nil, errors.New("watcher stopped")
	}
}

func (m *mockWatcher) Stop() {
	select {
	case <-m.exit:
		return
	default:
		close(m.exit)
	}
}
