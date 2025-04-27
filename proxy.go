package httpx

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var proxiesRegex = []*regexp.Regexp{
	regexp.MustCompile(`^\S+:\d+@\S+:\S+$`), // ip:port@login:password
	regexp.MustCompile(`^\S+:\S+@\S+:\d+$`), // login:password@ip:port
	regexp.MustCompile(`^\S+:\d+:\S+:\S+$`), // ip:port:login:password
	regexp.MustCompile(`^\S+:\S+:\S+:\d+$`), // login:password:ip:port
	regexp.MustCompile(`^\S+:\d+$`),         // ip:port
}

type Proxy struct {
	Scheme   string // e.g., "http", "https", "socks4", "socks5"
	Host     string
	Port     string
	Username *string // optional
	Password *string // optional
}

func (p *Proxy) IsAuth() bool {
	return p.Username != nil && p.Password != nil
}

func (p *Proxy) String() string {
	if p.IsAuth() {
		return fmt.Sprintf("%s://%s:%s@%s:%s", p.Scheme, *p.Username, *p.Password, p.Host, p.Port)
	}
	return fmt.Sprintf("%s://%s:%s", p.Scheme, p.Host, p.Port)
}

func (p *Proxy) TransportFunction() (func(*http.Request) (*url.URL, error), error) {
	proxyURL, err := url.Parse(p.String())
	if err != nil {
		return nil, err
	}

	return http.ProxyURL(proxyURL), nil
}

func NewProxyFromLine(rawProxy string) (*Proxy, error) {
	// Parse the proxy line and return a Proxy struct
	// Example: "http://username:password@host:port"

	p := &Proxy{}

	p.Scheme = strings.Split(rawProxy, "://")[0]
	if !isSchemeSupported(p.Scheme) {
		return nil, fmt.Errorf("unsupported proxy scheme: %s", p.Scheme)
	}

	rawProxy = strings.TrimPrefix(rawProxy, p.Scheme+"://")

	for i, regex := range proxiesRegex {
		if !regex.MatchString(rawProxy) {
			continue
		}

		parts := strings.Split(rawProxy, ":")

		switch i {
		case 0: // ip:port@login:password
			p.Host = parts[0]
			p.Port = strings.Split(parts[1], "@")[0]
			p.Username = &parts[2]
			p.Password = &parts[3]
		case 1: // login:password@ip:port
			p.Host = strings.Split(parts[1], "@")[1]
			p.Port = parts[2]
			p.Username = &parts[0]
			p.Password = &strings.Split(parts[1], "@")[0]
		case 2: // ip:port:login:password
			p.Host = parts[0]
			p.Port = parts[1]
			p.Username = &parts[2]
			p.Password = &parts[3]
		case 3: // login:password:ip:port
			p.Host = strings.Split(parts[2], "@")[1]
			p.Port = parts[3]
			p.Username = &parts[0]
			p.Password = &strings.Split(parts[2], "@")[0]
		case 4: // ip:port
			p.Host = parts[0]
			p.Port = parts[1]
		}
	}

	if p.Host == "" || p.Port == "" {
		return nil, fmt.Errorf("invalid proxy format: %s", rawProxy)
	}

	if p.Username != nil && p.Password != nil {
		if *p.Username == "" || *p.Password == "" {
			return nil, fmt.Errorf("username or password cannot be empty if provided")
		}
	}

	return p, nil
}

func NewProxy(opt *Proxy) (*Proxy, error) {
	if opt == nil {
		return nil, fmt.Errorf("proxy options cannot be nil")
	}

	if opt.Host == "" || opt.Port == "" {
		return nil, fmt.Errorf("host and port cannot be empty")
	}

	if !isSchemeSupported(opt.Scheme) {
		return nil, fmt.Errorf("unsupported proxy scheme: %s", opt.Scheme)
	}

	if opt.Username != nil && opt.Password != nil {
		if *opt.Username == "" || *opt.Password == "" {
			return nil, fmt.Errorf("username or password cannot be empty if provided")
		}
	}

	return opt, nil
}

func isSchemeSupported(scheme string) bool {
	supportedSchemes := []string{"http", "https", "socks4", "socks5"}
	for _, s := range supportedSchemes {
		if scheme == s {
			return true
		}
	}
	return false
}
