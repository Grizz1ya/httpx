package httpx

import "testing"

func TestSession(t *testing.T) {
	s := NewSession()

	if s == nil {
		t.Error("NewSession() returned nil")
	}

	s.AddHeader("key", "value")
	if s.headers["key"] != "value" {
		t.Error("AddHeader() failed")
	}

	s.RemoveHeader("key")
	if _, ok := s.headers["key"]; ok {
		t.Error("RemoveHeader() failed")
	}

	s.Proxy(nil)
	if s.client.Transport != nil {
		t.Error("Proxy() failed")
	}

	response, err := s.Get("http://example.com").Params(map[string]interface{}{"key": "value"}).Do()
	if err != nil {
		t.Error("Get() failed")
	}

	if response == nil {
		t.Error("Get() failed")
	}

	t.Logf("Response: %v", response.Text())
}