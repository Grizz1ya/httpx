package httpx

import (
	"io"
	"net/http"

	"github.com/Anderson-Lu/gofasion/gofasion"
)

type Response struct {
	response *http.Response

	StatusCode int
}

func (r *Response) Json() *gofasion.Fasion {
	return gofasion.NewFasion(r.Text())
}

func (r *Response) Text() string {
	defer r.response.Body.Close()

	content, err := io.ReadAll(r.response.Body)
	if err != nil {
		panic(err)
	}

	return string(content)
}

func (r *Response) Cookies() []*http.Cookie {
	return r.response.Cookies()
}

func (r *Response) Headers() http.Header {
	return r.response.Header
}
