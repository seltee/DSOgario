package tokens

import (
	"sync"
	"time"
)

type PlayerInfo struct {
	Name       string
	ColorIndex uint8
	CreatedAt  time.Time
}

type Manager struct {
	mu            sync.RWMutex
	tokens        map[string]PlayerInfo // token -> info
	cleanupTicker *time.Ticker
	done          chan struct{} // for graceful shutdown
}

func New() *Manager {
	m := &Manager{
		tokens: make(map[string]PlayerInfo),
		done:   make(chan struct{}),
	}

	m.startCleanupRoutine()
	return m
}

func (m *Manager) AddNewUser(name string, colorIndex uint8) string {
	for {
		token := randToken(12) // generate outside the lock!

		if !m.isTokenUsed(token) {
			m.mu.Lock()
			if _, exists := m.tokens[token]; !exists {
				m.tokens[token] = PlayerInfo{
					Name:       name,
					ColorIndex: colorIndex,
					CreatedAt:  time.Now(),
				}
				m.mu.Unlock()
				return token
			}
			m.mu.Unlock()
		}
	}
}

func (m *Manager) isTokenUsed(token string) bool {
	m.mu.RLock()
	_, ok := m.tokens[token]
	m.mu.RUnlock()
	return ok
}

func (m *Manager) Remove(token string) {
	m.mu.Lock()
	delete(m.tokens, token)
	m.mu.Unlock()
}

func (m *Manager) Validate(token string) (PlayerInfo, bool) {
	m.mu.RLock()
	info, ok := m.tokens[token]
	m.mu.RUnlock()
	return info, ok
}

func (m *Manager) GetAll() map[string]PlayerInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	copy := make(map[string]PlayerInfo, len(m.tokens))
	for k, v := range m.tokens {
		copy[k] = v
	}
	return copy
}

// Graceful shutdown
func (m *Manager) Stop() {
	close(m.done)
}
