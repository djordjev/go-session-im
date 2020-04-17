package imsession

import (
	"errors"
	"time"
)

type sessionInfo struct {
	timestamp time.Time
	payload   interface{}
}

/*
SessionManager type that holds active sessions in the map
*/
type SessionManager struct {
	sessions      map[string]*sessionInfo
	signalStop    chan bool
	maxInactivity time.Duration
}

// Initialize Creates a new session manager
func Initialize(inactivityDuration time.Duration) *SessionManager {
	manager := SessionManager{
		sessions:      make(map[string]*sessionInfo),
		signalStop:    make(chan bool),
		maxInactivity: inactivityDuration,
	}

	go garbageCollector(&manager)

	return &manager
}

/*
Create Creates a new session.
Returns token as a string
*/
func (manager *SessionManager) Create(payload interface{}) string {
	info := sessionInfo{
		payload:   payload,
		timestamp: time.Now(),
	}
	token := generateRandomString()

	manager.sessions[token] = &info

	return token
}

/*
Get Takes session from session manager and returns it's associated payload.
Returns error if no session is associated to token. Updates last access time
so the session won't be cleared
*/
func (manager *SessionManager) Get(token string) (interface{}, error) {
	info, exists := manager.sessions[token]

	if !exists {
		return nil, errors.New("Session with token does not exist")
	}

	info.timestamp = time.Now()

	return info.payload, nil
}

/*
Update updates payload of selected session
*/
func (manager *SessionManager) Update(token string, payload interface{}) error {
	info, exists := manager.sessions[token]

	if !exists {
		return errors.New("Session with token does not exist")
	}

	info.payload = payload
	return nil
}

/*
Kill Removes session form manager
*/
func (manager *SessionManager) Kill(token string) error {
	_, present := manager.sessions[token]

	if !present {
		return errors.New("No session associated with token")
	}

	delete(manager.sessions, token)
	return nil
}

/*
Stop Function should be called when SessionManger will not be used anymore
like gracefull exit of whole program. It stops garbage collection goroutine.
*/
func (manager *SessionManager) Stop() {
	manager.signalStop <- true
}

func garbageCollector(manager *SessionManager) {
	for {

		select {
		case <-manager.signalStop:
			return
		case <-time.After(time.Minute * 2):
			break
		}

		for k, v := range manager.sessions {
			diff := time.Since(v.timestamp)

			if diff > manager.maxInactivity {
				delete(manager.sessions, k)
			}
		}
	}
}
