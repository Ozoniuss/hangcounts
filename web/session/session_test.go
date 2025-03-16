package session

import (
	"encoding/base64"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSessionId_ReturnsABase64EncodedSession(t *testing.T) {
	s, err := generateSecureSessionId()
	assert.Nil(t, err, "expected no error when generating a session id")

	b, err := base64.RawURLEncoding.DecodeString(s)
	assert.Nil(t, err, "expected session to be base64 encoded")
	assert.Equal(t, 32, len(b), "expected decoded session id to have 32")
}

func TestGenerateSessionId_Returns100DifferentSessions_WhenCallingTheGenerator100Times(t *testing.T) {
	existing := make([]string, 100)
	for i := range 100 {
		s, err := generateSecureSessionId()
		assert.Nil(t, err, "expected no error when generating a session id")
		assert.NotContains(t, existing, s)
		existing[i] = s
	}
	t.Log(existing)
}

func TestSessionManager_detectsNewSessionAsNotExpired(t *testing.T) {
	idleExp := 10 * time.Second
	absExp := 30 * time.Second
	m := NewSessionManager(nil, idleExp, absExp, "")

	s := Session{
		createdAt:      time.Now(),
		lastActivityAt: time.Now(),
	}

	assert.False(t, m.isExpired(s), "expected new session to not be expired")
}

func TestSessionManager_detectsSessionAsNotExpired_WhenIdleIsCloseToExpiration(t *testing.T) {
	synctest.Run(func() {
		idleExp := 10 * time.Second
		absExp := 30 * time.Second
		m := NewSessionManager(nil, idleExp, absExp, "")

		s := Session{
			createdAt:      time.Now(),
			lastActivityAt: time.Now(),
		}
		time.Sleep(9 * time.Second)

		assert.False(t, m.isExpired(s), "expected session to not be expired")
	})
}

func TestSessionManager_detectsSessionAsNotExpired_WhenCreatedForLongerThanIdle_AndLastActivityIsRefreshed(t *testing.T) {
	synctest.Run(func() {
		idleExp := 10 * time.Second
		absExp := 30 * time.Second
		m := NewSessionManager(nil, idleExp, absExp, "")

		s := Session{
			createdAt:      time.Now(),
			lastActivityAt: time.Now(),
		}
		time.Sleep(9 * time.Second)

		s.lastActivityAt = time.Now()
		time.Sleep(9 * time.Second)

		assert.False(t, m.isExpired(s), "expected session to not be expired")
	})
}

func TestSessionManager_detectsSessionAsExpired_WhenIdleExpiredIt(t *testing.T) {
	synctest.Run(func() {
		idleExp := 10 * time.Second
		absExp := 30 * time.Second
		m := NewSessionManager(nil, idleExp, absExp, "")

		s := Session{
			createdAt:      time.Now(),
			lastActivityAt: time.Now(),
		}
		time.Sleep(20 * time.Second)
		assert.True(t, m.isExpired(s), "expected session to be expired")
	})
}

func TestSessionManager_detectsSessionAsExpired_WhenAbsoluteExpiredIt_AndIdleDidnt(t *testing.T) {
	synctest.Run(func() {
		idleExp := 10 * time.Second
		absExp := 15 * time.Second
		m := NewSessionManager(nil, idleExp, absExp, "")

		s := Session{
			createdAt:      time.Now(),
			lastActivityAt: time.Now(),
		}
		time.Sleep(9 * time.Second)
		s.lastActivityAt = time.Now()
		time.Sleep(9 * time.Second)
		assert.True(t, m.isExpired(s), "expected session to be expired")
	})
}
