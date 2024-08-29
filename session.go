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

func (s *Session) request(method, url string) *Request {
	return &Request{
		method: method,
		url: url,
	}
}

func (s *Session) Get(url string) *Request {
	return s.request("GET", url)
}

func (s *Session) Post(url string) *Request {
	return s.request("POST", url)
}