package core

import (
	"context"
	"gochat/internal/infrastructure/metrics"
	"sync"

	"go.uber.org/zap"
)

// Manager 扮演会话池的角色，管理本地连接
type Manager struct {
	mu       sync.RWMutex
	sessions map[string]map[string]Session // map[UserID]map[SessionID]Session 以支持多端登录
}

func NewManager() *Manager {
	return &Manager{
		sessions: make(map[string]map[string]Session, 10000), // 可通过配置传入
	}
}

// Register 将 Session 注册到管理器
func (m *Manager) Register(s Session) {
	userID := s.UserID()
	sessionID := s.ID()

	m.mu.Lock()
	if m.sessions[userID] == nil {
		m.sessions[userID] = make(map[string]Session)
	}
	m.sessions[userID][sessionID] = s
	m.mu.Unlock()

	metrics.WSConnectionDelta(context.Background(), 1)
}

// Unregister 移除 Session
func (m *Manager) Unregister(s Session) {
	if s == nil || s.UserID() == "" || s.ID() == "" {
		return
	}

	userID := s.UserID()
	sessionID := s.ID()
	found := false

	m.mu.Lock()
	if userSessions, ok := m.sessions[userID]; ok {
		if _, exists := userSessions[sessionID]; exists {
			delete(userSessions, sessionID)
			found = true

			// 如果该用户所有设备的连接都断开了，清理 map 释放内存
			if len(userSessions) == 0 {
				delete(m.sessions, userID)
			}
		}
	}
	m.mu.Unlock()

	if found {
		metrics.WSConnectionDelta(context.Background(), -1)
		metrics.WSDisconnect(context.Background(), "unregister_by_pointer")
		zap.L().Debug("gateway manager: unregistered session", zap.String("userID", userID), zap.String("sessionID", sessionID))
	}
}

// GetByUserID 根据目标 UserID 获取当前的 Session（单体版调度）
func (m *Manager) GetByUserID(userID string) ([]Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	userSessions, ok := m.sessions[userID]
	if !ok || len(userSessions) == 0 {
		return nil, false
	}

	result := make([]Session, 0, len(userSessions))
	for _, s := range userSessions {
		result = append(result, s)
	}
	return result, true
}
