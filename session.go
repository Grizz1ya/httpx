package httpx

import (
	"net/http"
	"net/url"
)

type Session struct {
	client *http.Client

	headers map[string]string
}


func NewSession() *Session {
	return &Session{
		client: http.DefaultClient,
		headers: make(map[string]string),
	}
}

func (s *Session) Proxy(proxy func(*http.Request) (*url.URL, error)) {
	if proxy == nil {
		s.client.Transport = nil
	} else {
		s.client.Transport = &http.Transport{
			Proxy: proxy,
		}
	}
}


func (s *Session) AddStaticHeader(key, value string) {
	// * Add static headers to the session
	s.headers[key] = value
}

func (s *Session) RemoveStaticHeader(key string) {
	// * Remove static headers from the session
	delete(s.headers, key)
}

func (s *Session) Get(url string) *Request {
	return request("GET", url, s.client, s.headers)
}

func (s *Session) Post(url string) *Request {
	return request("POST", url, s.client, s.headers)
}

func request(method, url string, client *http.Client, headers map[string]string) *Request {
	if client == nil {
		client = http.DefaultClient
	}

	return &Request{
		method: method,
		url: url,
		client: client,
		headers: headers,
	}
}


// * Static methods

func Get(url string) *Request {
	return request("GET", url, nil, nil)
}

func Post(url string) *Request {
	return request("POST", url, nil, nil)
}