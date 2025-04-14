package httpx

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Session struct {
	client *http.Client

	headers map[string]string
}

func NewSession() *Session {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}
	return &Session{
		client:  client,
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

func (s *Session) AddCookie(domain, name, value string) {
	// * Add cookie to the session
	cookie := &http.Cookie{
		Name:   name,
		Value:  value,
		Domain: domain,
	}

	s.client.Jar.SetCookies(&url.URL{
		Scheme: "https",
		Host:   domain,
	}, []*http.Cookie{cookie})
}

func (s *Session) RemoveCookie(domain, name string) {
	// * Remove cookie from the session
	cookie := &http.Cookie{
		Name:   name,
		MaxAge: -1,
	}

	s.client.Jar.SetCookies(&url.URL{
		Scheme: "http",
		Host:   domain,
	}, []*http.Cookie{cookie})
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
		client = &http.Client{}
	}

	return &Request{
		method:        method,
		url:           url,
		client:        client,
		staticHeaders: headers,
	}
}

// * Static methods

func Get(url string) *Request {
	return request("GET", url, nil, nil)
}

func Post(url string) *Request {
	return request("POST", url, nil, nil)
}
