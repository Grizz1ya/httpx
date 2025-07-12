package httpx

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/Grizz1ya/httpx/utils"
	utls "github.com/refraction-networking/utls"
)

type Session struct {
	client *http.Client

	headers map[string]string

	Proxy *Proxy

	redirectEnabled       bool
	customRedirectHandler func(*Response) error
	maxRedirects          int

	cookieOrigins *utils.CookieOriginMap

	tlsFingerprint utls.ClientHelloID
}

func NewSession() *Session {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
	}
	return &Session{
		client:          client,
		headers:         make(map[string]string),
		cookieOrigins:   utils.NewCookieOriginMap(),
		redirectEnabled: true,
		maxRedirects:    5,
		tlsFingerprint:  utls.HelloChrome_131, // по умолчанию Chrome
	}
}

func (s *Session) SetTLSFingerprint(fingerprint utls.ClientHelloID) error {
	s.tlsFingerprint = fingerprint
	return s.rebuildTransport()
}

func (s *Session) SetTimeout(timeout time.Duration) {
	s.client.Timeout = timeout
}

func (s *Session) Redirect(enable bool, handler func(*Response) error) error {
	s.redirectEnabled = enable
	s.customRedirectHandler = handler
	return s.rebuildTransport()
}

func (s *Session) SetProxy(proxy *Proxy) error {
	s.Proxy = proxy
	return s.rebuildTransport()
}

func (s *Session) AddCookie(domain, name, value string) {
	// * Add cookie to the session
	cookie := &http.Cookie{
		Name:   name,
		Value:  value,
		Domain: domain,
		MaxAge: 10000000,
	}

	s.client.Jar.SetCookies(&url.URL{
		Scheme: "https",
		Host:   domain,
	}, []*http.Cookie{cookie})

	// * Add cookie to the list
	s.cookieOrigins.Add(domain, cookie)
}

func (s *Session) RemoveCookie(domain, name string) {
	// * Remove cookie from the session
	cookie := &http.Cookie{
		Name:   name,
		MaxAge: -1,
	}

	s.client.Jar.SetCookies(&url.URL{
		Scheme: "https",
		Host:   domain,
	}, []*http.Cookie{cookie})

	// * Remove cookie from the list
	s.cookieOrigins.Remove(domain, cookie)
}

func (s *Session) AddStaticHeader(key, value string) {
	// * Add static headers to the session
	s.headers[key] = value
}

func (s *Session) GetStaticHeader(key string) (string, bool) {
	// * Get static headers from the session
	value, ok := s.headers[key]
	return value, ok
}

func (s *Session) RemoveStaticHeader(key string) {
	// * Remove static headers from the session
	delete(s.headers, key)
}

func (s *Session) Get(url string) *Request {
	return request("GET", url, s.client, s.headers, s.cookieOrigins)
}

func (s *Session) Post(url string) *Request {
	return request("POST", url, s.client, s.headers, s.cookieOrigins)
}

func (s *Session) Options(url string) *Request {
	return request("OPTIONS", url, s.client, s.headers, s.cookieOrigins)
}

func (s *Session) Put(url string) *Request {
	return request("PUT", url, s.client, s.headers, s.cookieOrigins)
}

func (s *Session) rebuildTransport() error {
	defaultTr := http.DefaultTransport.(*http.Transport).Clone()

	var base *http.Transport
	if s.Proxy == nil {
		base = defaultTr.Clone()
	} else {
		trFn, err := s.Proxy.TransportFunction()
		if err != nil {
			return fmt.Errorf("failed to build proxy transport: %w", err)
		}
		base = defaultTr.Clone()
		base.Proxy = trFn
	}

	// 🌐 Подменяем TLS через uTLS
	base.DialTLSContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		dialer := &net.Dialer{}
		tcpConn, err := dialer.DialContext(ctx, network, addr)
		if err != nil {
			return nil, err
		}

		host, _, _ := net.SplitHostPort(addr)
		config := &utls.Config{
			ServerName: host,
		}

		utlsConn := utls.UClient(tcpConn, config, s.tlsFingerprint)
		if err := utlsConn.Handshake(); err != nil {
			tcpConn.Close()
			return nil, err
		}

		return utlsConn, nil
	}

	// 🔄 Обработка редиректов (как раньше)
	if !s.redirectEnabled {
		var lastResp *http.Response
		wrapped := utils.NewRedirectTransport(
			base,
			func(resp *http.Response) {
				lastResp = resp
			},
		)
		s.client.Transport = wrapped
		s.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if len(via) >= s.maxRedirects {
				return fmt.Errorf("stopped after %d redirects", s.maxRedirects)
			}
			if s.customRedirectHandler != nil {
				return s.customRedirectHandler(&Response{
					response:   lastResp,
					StatusCode: lastResp.StatusCode,
					URL:        lastResp.Request.URL.String(),
				})
			}
			return http.ErrUseLastResponse
		}
	} else {
		s.client.Transport = base
		s.client.CheckRedirect = nil
	}

	return nil
}

func (s *Session) Cookies(domains ...string) []*http.Cookie {
	if len(domains) == 0 {
		all := s.cookieOrigins.All()
		var result []*http.Cookie
		for _, cookies := range all {
			result = append(result, cookies...)
		}
		return result
	}

	var result []*http.Cookie
	for _, domain := range domains {
		result = append(result, s.cookieOrigins.Get(domain)...)
	}
	return result
}

func (s *Session) Destroy() {
	// Очистить куки: создаём новый jar
	jar, _ := cookiejar.New(nil)
	s.client.Jar = jar

	// Очистить заголовки
	for k := range s.headers {
		delete(s.headers, k)
	}

	// Обнулить redirect-обработку
	s.client.CheckRedirect = nil
	s.redirectEnabled = false
	s.customRedirectHandler = nil

	// Обнулить proxy и транспорт (это закроет соединения и освободит память)
	s.Proxy = nil
	s.client.Transport = nil

	// Сброс timeout и прочего — опционально
	s.client.Timeout = 0

	// Сброс TLS fingerprint — если хочешь
	// s.tlsFingerprint = utls.HelloChrome_131
}
