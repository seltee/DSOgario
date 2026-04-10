package tokens

import (
	"time"
)

func (m *Manager) startCleanupRoutine() {
	m.cleanupTicker = time.NewTicker(60 * time.Second) // every minute

	go func() {
		for {
			select {
			case <-m.cleanupTicker.C:
				m.cleanupExpiredTokens()
			case <-m.done:
				m.cleanupTicker.Stop()
				return
			}
		}
	}()
}

func (m *Manager) cleanupExpiredTokens() {
	const maxAge = 1 * time.Minute

	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for token, info := range m.tokens {
		if now.Sub(info.CreatedAt) > maxAge {
			delete(m.tokens, token)
		}
	}
}
