package httpx

import (
	"io"
	"net/http"

	"github.com/Anderson-Lu/gofasion/gofasion"
)

type Response struct {
	response *http.Response

	StatusCode int
	URL        string
}

func (r *Response) Json() (*gofasion.Fasion, error) {
	text, err := r.Text()
	if err != nil {
		return nil, err
	}
	return gofasion.NewFasion(text), nil
}

func (r *Response) Text() (string, error) {
	defer r.response.Body.Close()

	content, err := io.ReadAll(r.response.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (r *Response) Cookies() []*http.Cookie {
	return r.response.Cookies()
}

func (r *Response) Headers() http.Header {
	return r.response.Header
}
