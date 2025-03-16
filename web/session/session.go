package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/Ozoniuss/hangcounts/domain/model"
	"github.com/Ozoniuss/hangcounts/domain/storage"
)

func generateSecureSessionId() (string, error) {
	id := make([]byte, 32)

	_, err := io.ReadFull(rand.Reader, id)
	if err != nil {
		return "", fmt.Errorf("could not generate random id: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(id), nil
}

const (
	SESSION_COOKIE_NAME = "session_id"
)

type Session struct {
	createdAt      time.Time
	lastActivityAt time.Time
	id             string
	userID         model.IndividualId
	logger         *slog.Logger
}

func newSession() (Session, error) {
	id, err := generateSecureSessionId()
	if err != nil {
		return Session{}, fmt.Errorf("failed to generate a session id: %w", err)
	}
	return Session{
		id:             id,
		createdAt:      time.Now(),
		lastActivityAt: time.Now(),
	}, nil
}

type SessionManager struct {
	storage            storage.AppStorage
	idleExpiration     time.Duration
	absoluteExpiration time.Duration
	cookieName         string
}

func NewSessionManager(
	store storage.AppStorage,
	idleExpiration,
	absoluteExpiration time.Duration,
	cookieName string,
) *SessionManager {
	return &SessionManager{
		storage:            store,
		idleExpiration:     idleExpiration,
		absoluteExpiration: absoluteExpiration,
		cookieName:         cookieName,
	}
}

func (m *SessionManager) isExpired(session Session) bool {
	return time.Since(session.createdAt) > m.absoluteExpiration ||
		time.Since(session.lastActivityAt) > m.idleExpiration
}
