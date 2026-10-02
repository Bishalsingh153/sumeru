package orm

import "sync"

var (
	mailThreadMu     sync.RWMutex
	mailThreadModels = map[string]bool{}
)

// SetModelMailThread marks a model as chatter-capable (auto-subscribe on create).
func SetModelMailThread(modelName string, enabled bool) {
	mailThreadMu.Lock()
	defer mailThreadMu.Unlock()
	if enabled {
		mailThreadModels[modelName] = true
		return
	}
	delete(mailThreadModels, modelName)
}

// ModelHasMailThread reports whether model uses mail thread hooks.
func ModelHasMailThread(modelName string) bool {
	mailThreadMu.RLock()
	defer mailThreadMu.RUnlock()
	return mailThreadModels[modelName]
}
