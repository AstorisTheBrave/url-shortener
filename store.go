package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"sync"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	slugLen  = 6
)

type Store struct {
	mu   sync.RWMutex
	urls map[string]string // slug -> original url
	path string
}

func NewStore(path string) *Store {
	s := &Store{urls: make(map[string]string), path: path}
	s.load()
	return s
}

func (s *Store) Set(url string) (string, error) {
	slug, err := randomSlug()
	if err != nil {
		return "", err
	}
	s.mu.Lock()
	s.urls[slug] = url
	s.mu.Unlock()
	s.save()
	return slug, nil
}

func (s *Store) Get(slug string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, ok := s.urls[slug]
	return url, ok
}

func randomSlug() (string, error) {
	b := make([]byte, slugLen)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", fmt.Errorf("rand: %w", err)
		}
		b[i] = alphabet[n.Int64()]
	}
	return string(b), nil
}

func (s *Store) save() {
	if s.path == "" {
		return
	}
	s.mu.RLock()
	data, _ := json.Marshal(s.urls)
	s.mu.RUnlock()
	_ = os.WriteFile(s.path, data, 0600)
}

func (s *Store) load() {
	if s.path == "" {
		return
	}
	data, err := os.ReadFile(s.path)
	if err != nil {
		return
	}
	_ = json.Unmarshal(data, &s.urls)
}
