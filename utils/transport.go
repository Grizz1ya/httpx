package utils

import "net/http"

type RedirectTransport struct {
	base  http.RoundTripper
	store func(*http.Response)
}

func (t *RedirectTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		t.store(resp)
	}
	return resp, nil
}

func NewRedirectTransport(base http.RoundTripper, store func(*http.Response)) *RedirectTransport {
	return &RedirectTransport{
		base:  base,
		store: store,
	}
}
