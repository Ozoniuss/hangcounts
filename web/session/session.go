package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/Ozoniuss/hangcounts/domain/model"
	"github.com/Ozoniuss/hangcounts/domain/storage"
)

func GenerateSecureSessionId() (string, error) {
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

type SessionStorage interface {
	StoreSession(ctx context.Context, session Session) error
	GetSession(ctx context.Context, cookie string) (Session, error)
	UpdateLastAccessed(ctx context.Context, cookie string, lastAccessed string) error
}

var ErrNotFound = errors.New("session not found in database")
var ErrUnknown = errors.New("unknown error")
var ErrUserNotFound = errors.New("user not found for session")
var ErrUserDeleted = errors.New("attempted to create a session for a deleted user")
var ErrCookieInvalidLength = errors.New("encoded cookie should have 44 characters")

type Session struct {
	CreatedAt    time.Time
	LastAccessed time.Time
	CookieValue  string
	UserID       model.IndividualId
}

func NewSessionForUser(userId model.IndividualId) (Session, error) {
	id, err := GenerateSecureSessionId()
	if err != nil {
		return Session{}, fmt.Errorf("failed to generate a session id: %w", err)
	}
	return Session{
		CookieValue:  id,
		CreatedAt:    time.Now(),
		LastAccessed: time.Now(),
		UserID:       userId,
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
	return time.Since(session.CreatedAt) > m.absoluteExpiration ||
		time.Since(session.LastAccessed) > m.idleExpiration
}
