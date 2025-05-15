package httpx

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/Grizz1ya/httpx/utils"
)

type Session struct {
	client *http.Client

	headers map[string]string

	Proxy *Proxy

	redirectEnabled       bool
	customRedirectHandler func(*Response) error

	cachedCookieDomains *utils.CachedCookieDomains
}

func NewSession() *Session {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}
	return &Session{
		client:              client,
		headers:             make(map[string]string),
		cachedCookieDomains: utils.NewCachedCookieDomains(),
	}
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
	}

	s.client.Jar.SetCookies(&url.URL{
		Scheme: "https",
		Host:   domain,
	}, []*http.Cookie{cookie})

	// * Add cookie to the list of domains for cookies
	s.cachedCookieDomains.Add(domain)
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
	return request("GET", url, s.client, s.headers, s.cachedCookieDomains)
}

func (s *Session) Post(url string) *Request {
	return request("POST", url, s.client, s.headers, s.cachedCookieDomains)
}

func (s *Session) Options(url string) *Request {
	return request("OPTIONS", url, s.client, s.headers, s.cachedCookieDomains)
}

func (s *Session) Put(url string) *Request {
	return request("PUT", url, s.client, s.headers, s.cachedCookieDomains)
}

func (s *Session) rebuildTransport() error {
	// получаем независимый экземпляр транспортa
	defaultTr := http.DefaultTransport.(*http.Transport).Clone()

	// если прокси не нужен — используем defaultTr,
	// иначе клонируем default и ставим нужный Proxy-функцию
	var base *http.Transport
	if s.Proxy == nil {
		base = defaultTr
	} else {
		trFn, err := s.Proxy.TransportFunction()
		if err != nil {
			return fmt.Errorf("failed to build proxy transport: %w", err)
		}
		base = defaultTr.Clone()
		base.Proxy = trFn
	}

	// теперь base — ваш полностью независимый транспорт,
	// в который вы можете врезать обёртку для редиректов
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

func cloneCookie(c *http.Cookie) *http.Cookie {
	copy := *c
	return &copy
}

func (s *Session) Cookies(domains ...string) []*http.Cookie {
	jar, ok := s.client.Jar.(*cookiejar.Jar)
	if !ok {
		return nil
	}

	cookieSet := make(map[string]*http.Cookie)

	var filterDomains []string
	if len(domains) > 0 {
		for _, d := range domains {
			filterDomains = append(filterDomains, strings.TrimPrefix(d, "."))
		}
	} else {
		filterDomains = s.cachedCookieDomains.Domains
	}

	for _, visited := range s.cachedCookieDomains.Domains {
		for _, filter := range filterDomains {
			// Совпадение точное или по поддомену
			if visited == filter || strings.HasSuffix(visited, "."+filter) {
				u := &url.URL{Scheme: "https", Host: visited}
				for _, c := range jar.Cookies(u) {
					// Подставим домен, если он не указан
					if c.Domain == "" {
						c = cloneCookie(c) // не мутируем оригинал
						c.Domain = visited // явно указываем источник
					}
					key := c.Name + "|" + c.Domain // уникальный ключ
					cookieSet[key] = c
				}
				break
			}
		}
	}

	result := make([]*http.Cookie, 0, len(cookieSet))
	for _, c := range cookieSet {
		result = append(result, c)
	}

	return result
}

func request(method, url string, client *http.Client, headers map[string]string, cachedCookieDomains *utils.CachedCookieDomains) *Request {
	if client == nil {
		client = &http.Client{}
	}

	return &Request{
		method:              method,
		url:                 url,
		client:              client,
		staticHeaders:       headers,
		cachedCookieDomains: cachedCookieDomains,
	}
}

// * Static methods
func Get(url string) *Request {
	return request("GET", url, nil, nil, nil)
}

func Post(url string) *Request {
	return request("POST", url, nil, nil, nil)
}
