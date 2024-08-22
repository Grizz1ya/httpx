package httpx

import "net/http"

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