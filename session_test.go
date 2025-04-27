package httpx

import (
	"testing"
)

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

	t.Logf("[1] Proxy: %v", s.Proxy)

	s.SetProxy(nil)
	t.Logf("[2] Proxy: %v", s.Proxy)

	proxy, err := NewProxyFromLine("http://127.0.0.1:8080")
	if err != nil {
		t.Error("NewProxyFromLine() failed")
	}
	s.SetProxy(proxy)
	t.Logf("[3] Proxy: %v", s.Proxy)

	s.AddStaticHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	s.AddCookie("domain.com", "auth_token", "qwerty123")

	response, err := s.Get("https://x.com/").Do()
	if err != nil {
		t.Error("Get() failed")
	}

	if response == nil {
		t.Error("Get() failed")
	}

	t.Logf("Response: %v", response)

	t.Logf("Response text: %v", response.Text())

	// response, err := s.Get("http://example.com").Params(map[string]interface{}{"param1": "value"}).Do()
	// if err != nil {
	// 	t.Error("Get() failed")
	// }

	// if response == nil {
	// 	t.Error("Get() failed")
	// }

	// t.Logf("Response: %v", response.Text())

	// response, err = Get("http://example.com").Headers(map[string]string{"header1": "value"}).Do()
	// if err != nil {
	// 	t.Error("Get() failed")
	// }

	// if response == nil {
	// 	t.Error("Get() failed")
	// }

	// t.Logf("Response: %v", response.Text())

}
