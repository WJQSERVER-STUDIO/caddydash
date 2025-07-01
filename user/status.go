package user

import "sync"

type UserStatus struct {
	IsInitialized bool
	mu            sync.Mutex
}

func (s *UserStatus) SetInitialized(initialized bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.IsInitialized = initialized
}

func (s *UserStatus) IsUserInitialized() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.IsInitialized
}

var userStatus UserStatus
