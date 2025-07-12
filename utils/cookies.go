package utils

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

type CookieOriginMap struct {
	data map[string][]*http.Cookie
	mu   sync.Mutex
}

func NewCookieOriginMap() *CookieOriginMap {
	return &CookieOriginMap{
		data: make(map[string][]*http.Cookie),
	}
}

func (m *CookieOriginMap) Add(domain string, cookies ...*http.Cookie) {
	m.mu.Lock()
	defer m.mu.Unlock()

	domain = strings.TrimPrefix(domain, ".")
	existing := m.data[domain]
	existingMap := make(map[string]struct{}, len(existing))
	for _, c := range existing {
		existingMap[c.Name] = struct{}{}
	}

	for _, c := range cookies {
		if c == nil || c.Value == "" || c.MaxAge <= 0 || (!c.Expires.IsZero() && c.Expires.Before(time.Now())) {
			continue // удалённая / пустая
		}
		if _, ok := existingMap[c.Name]; ok {
			continue // уже есть
		}
		existing = append(existing, c)
	}

	m.data[domain] = existing
}

func (m *CookieOriginMap) Remove(domain string, cookies ...*http.Cookie) {
	m.mu.Lock()
	defer m.mu.Unlock()

	domain = strings.TrimPrefix(domain, ".")
	existing := m.data[domain]
	newList := make([]*http.Cookie, 0, len(existing))

	toDelete := make(map[string]bool)
	for _, c := range cookies {
		toDelete[c.Name] = true
	}

	for _, c := range existing {
		if !toDelete[c.Name] {
			newList = append(newList, c)
		}
	}

	m.data[domain] = newList
}

func (m *CookieOriginMap) Get(domain string) []*http.Cookie {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.data[strings.TrimPrefix(domain, ".")]
}

func (m *CookieOriginMap) All() map[string][]*http.Cookie {
	m.mu.Lock()
	defer m.mu.Unlock()
	copy := make(map[string][]*http.Cookie)
	for k, v := range m.data {
		copy[k] = append([]*http.Cookie{}, v...)
	}
	return copy
}
