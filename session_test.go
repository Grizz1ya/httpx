package httpx

import "testing"

func TestSession(t *testing.T) {
	s := NewSession()

	if s == nil {
		t.Error("NewSession() returned nil")
	}

	s.AddStaticHeader("key", "value")
	if s.headers["key"] != "value" {
		t.Error("AddStaticHeader() failed")
	}

	s.RemoveStaticHeader("key")
	if _, ok := s.headers["key"]; ok {
		t.Error("RemoveStaticHeader() failed")
	}

	s.Proxy(nil)
	if s.client.Transport != nil {
		t.Error("Proxy() failed")
	}

	response, err := s.Get("http://example.com").Params(map[string]interface{}{"param1": "value"}).Do()
	if err != nil {
		t.Error("Get() failed")
	}

	if response == nil {
		t.Error("Get() failed")
	}

	t.Logf("Response: %v", response.Text())

	response, err = Get("http://example.com").Headers(map[string]string{"header1": "value"}).Do()
	if err != nil {
		t.Error("Get() failed")
	}

	if response == nil {
		t.Error("Get() failed")
	}

	t.Logf("Response: %v", response.Text())

}