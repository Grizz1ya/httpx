package httpx

import (
	"net/http"
	"net/url"
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

	s.Proxy(nil)
	if s.client.Transport != nil {
		t.Error("Proxy() failed")
	}

	proxyRaw, err := url.Parse("http://127.0.0.1:8080")
	if err != nil {
		t.Error("url.Parse() failed")
	}
	s.Proxy(http.ProxyURL(proxyRaw))
	t.Logf("Proxy: %v", s.client.Transport)

	s.AddStaticHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	s.AddCookie("x.com", "auth_token", "ae773328f828f1d9f7f4d324363b9068fc1e73bbdae0f5a2f590bce46e829f4dfd11aea4805ef8f2708b97408ee748199d298503bce7b6b7cbd4b7226671ef99564c2a15257d5a4454f3642afb1f0659")

	response, err := s.Get("https://x.com/").Do()
	if err != nil {
		t.Error("Get() failed")
	}

	if response == nil {
		t.Error("Get() failed")
	}

	t.Logf("Response: %v", response.Text())

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
