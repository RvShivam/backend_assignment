package main

import (
	"sync"
	"time"
)

const (
	rateLimit = 5
	window    = time.Minute
)

type UserData struct {
	Timestamps []time.Time
	Rejected   int
}

type UserStats struct {
	AcceptedInWindow int `json:"accepted_in_window"`
	RejectedTotal    int `json:"rejected_total"`
}

type Store struct {
	mu    sync.Mutex
	users map[string]*UserData
}

func NewStore() *Store {
	return &Store{
		users: make(map[string]*UserData),
	}
}

func (s *Store) Accept(userID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	data, exists := s.users[userID]
	if !exists {
		data = &UserData{}
		s.users[userID] = data
	}

	s.pruneOld(now, data)

	if len(data.Timestamps) >= rateLimit {
		data.Rejected++
		return false
	}

	data.Timestamps = append(data.Timestamps, now)
	return true
}

func (s *Store) Stats() map[string]UserStats {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	result := make(map[string]UserStats, len(s.users))

	for id, data := range s.users {
		s.pruneOld(now, data)
		result[id] = UserStats{
			AcceptedInWindow: len(data.Timestamps),
			RejectedTotal:    data.Rejected,
		}
	}

	return result
}

func (s *Store) pruneOld(now time.Time, data *UserData) {
	cutoff := now.Add(-window)
	i := 0
	for i < len(data.Timestamps) && !data.Timestamps[i].After(cutoff) {
		i++
	}
	data.Timestamps = data.Timestamps[i:]
}