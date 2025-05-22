package httpx

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

var proxyPatterns = []struct {
	regex  *regexp.Regexp
	parser func(p *Proxy, groups []string)
}{
	{ // login:password:ip:port
		regexp.MustCompile(`^(\w+)://(\S+):(\S+):(\S+):(\d+)$`),
		func(p *Proxy, g []string) {
			p.Scheme = g[1]
			u, pw := g[2], g[3]
			p.Username, p.Password = &u, &pw
			p.Host, p.Port = g[4], g[5]
		},
	},
	{ // login:password@ip:port
		regexp.MustCompile(`^(\w+)://(\S+):(\S+)@(\S+):(\d+)$`),
		func(p *Proxy, g []string) {
			p.Scheme = g[1]
			u, pw := g[2], g[3]
			p.Username, p.Password = &u, &pw
			p.Host, p.Port = g[4], g[5]
		},
	},
	{ // ip:port:login:password
		regexp.MustCompile(`^(\w+)://(\S+):(\d+):(\S+):(\S+)$`),
		func(p *Proxy, g []string) {
			p.Scheme = g[1]
			p.Host, p.Port = g[2], g[3]
			u, pw := g[4], g[5]
			p.Username, p.Password = &u, &pw
		},
	},
	{ // ip:port@login:password
		regexp.MustCompile(`^(\w+)://(\S+):(\d+)@(\S+):(\S+)$`),
		func(p *Proxy, g []string) {
			p.Scheme = g[1]
			p.Host, p.Port = g[2], g[3]
			u, pw := g[4], g[5]
			p.Username, p.Password = &u, &pw
		},
	},
	{ // ip:port
		regexp.MustCompile(`^(\w+)://(\S+):(\d+)$`),
		func(p *Proxy, g []string) {
			p.Scheme = g[1]
			p.Host, p.Port = g[2], g[3]
		},
	},
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

	for _, pattern := range proxyPatterns {
		matches := pattern.regex.FindStringSubmatch(rawProxy)
		if len(matches) > 0 {
			pattern.parser(p, matches)
			return p, nil
		}
	}

	return p, fmt.Errorf("invalid proxy format: %s", rawProxy)
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
