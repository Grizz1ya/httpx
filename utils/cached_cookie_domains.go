package utils

type CachedCookieDomains struct {
	Domains []string
}

func (c *CachedCookieDomains) Add(domain string) {
	for _, d := range c.Domains {
		if d == domain {
			return
		}
	}
	c.Domains = append(c.Domains, domain)
}

func NewCachedCookieDomains() *CachedCookieDomains {
	return &CachedCookieDomains{
		Domains: make([]string, 0),
	}
}
