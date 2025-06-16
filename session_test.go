package httpx

import (
	"testing"

	utls "github.com/refraction-networking/utls"
)

func TestSession(t *testing.T) {
	s := NewSession()

	if s == nil {
		t.Error("NewSession() returned nil")
	}

	s.SetTLSFingerprint(utls.HelloEdge_106)

	// s.AddStaticHeader("key", "value")
	// if s.headers["key"] != "value" {
	// 	t.Error("AddStaticHeader() failed")
	// }

	// s.RemoveStaticHeader("key")
	// if _, ok := s.headers["key"]; ok {
	// 	t.Error("RemoveStaticHeader() failed")
	// }

	// t.Logf("[1] Proxy: %v", s.Proxy)

	// s.SetProxy(nil)
	// t.Logf("[2] Proxy: %v", s.Proxy)

	proxy, err := NewProxyFromLine("http://127.0.0.1:8080")
	if err != nil {
		t.Error("NewProxyFromLine() failed")
	}
	s.SetProxy(proxy)
	t.Logf("[3] Proxy: %v", s.Proxy)

	s.AddStaticHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.2 Safari/605.1.15")
	s.AddStaticHeader("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")

	response, err := s.Get("https://tls.peet.ws/api/all").Do()
	if err != nil {
		t.Error("Get() failed")
	}

	t.Logf("Response: %v", response)

	// s.Get("https://adsmanager.facebook.com/").Do()
	// s.Get("https://business.facebook.com/").Do()

	// cookies := s.Cookies()
	// for _, cookie := range cookies {
	// 	var expires string
	// 	if cookie.Expires.IsZero() {
	// 		expires = strconv.FormatInt(time.Now().Add(10*365*24*time.Hour).Unix(), 10)
	// 	} else {
	// 		expires = strconv.FormatInt(cookie.Expires.Unix(), 10)
	// 	}
	// 	t.Logf("Domain: %s, Name: %s, Value: %s, Expires: %s",
	// 		cookie.Domain,
	// 		cookie.Name,
	// 		cookie.Value,
	// 		expires,
	// 	)
	// }
	text, err := response.Text()
	if err != nil {
		t.Error("Response.Text() failed")
	}

	t.Logf("Response text: %v", text)

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
