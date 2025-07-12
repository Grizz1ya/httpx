package utils

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

// CookieOriginMap хранит: домен → (имя cookie → *http.Cookie)
type CookieOriginMap struct {
	data map[string]map[string]*http.Cookie
	mu   sync.RWMutex
}

func NewCookieOriginMap() *CookieOriginMap {
	return &CookieOriginMap{
		data: make(map[string]map[string]*http.Cookie),
	}
}

// Add добавляет/заменяет куки для домена.
func (m *CookieOriginMap) Add(domain string, cookies ...*http.Cookie) {
	now := time.Now()
	domain = strings.TrimPrefix(domain, ".")
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.data[domain]; !ok {
		m.data[domain] = make(map[string]*http.Cookie, len(cookies))
	}
	for _, c := range cookies {
		if c == nil || c.Value == "" || c.MaxAge <= 0 ||
			(!c.Expires.IsZero() && c.Expires.Before(now)) {
			continue // пропускаем “пустышки”
		}
		m.data[domain][c.Name] = c // заменяем по имени
	}
}

// Remove удаляет конкретные куки по имени.
func (m *CookieOriginMap) Remove(domain string, cookies ...*http.Cookie) {
	domain = strings.TrimPrefix(domain, ".")
	m.mu.Lock()
	defer m.mu.Unlock()

	dict, ok := m.data[domain]
	if !ok {
		return
	}
	for _, c := range cookies {
		delete(dict, c.Name)
	}
	if len(dict) == 0 {
		delete(m.data, domain)
	}
}

// Get возвращает **копию** []*http.Cookie для домена.
func (m *CookieOriginMap) Get(domain string) []*http.Cookie {
	domain = strings.TrimPrefix(domain, ".")
	m.mu.RLock()
	defer m.mu.RUnlock()

	dict, ok := m.data[domain]
	if !ok {
		return nil
	}
	out := make([]*http.Cookie, 0, len(dict))
	for _, c := range dict {
		out = append(out, c)
	}
	return out
}

// All возвращает копию map[домен][]*http.Cookie
func (m *CookieOriginMap) All() map[string][]*http.Cookie {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make(map[string][]*http.Cookie, len(m.data))
	for domain, dict := range m.data {
		clist := make([]*http.Cookie, 0, len(dict))
		for _, c := range dict {
			clist = append(clist, c)
		}
		out[domain] = clist
	}
	return out
}
