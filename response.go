package httpx

import (
	"io"
	"net/http"

	"github.com/Anderson-Lu/gofasion/gofasion"
	"golang.org/x/net/html/charset"
)

type Response struct {
	response *http.Response
}

func (r *Response) Json() *gofasion.Fasion {
	return gofasion.NewFasion(r.Text())
}

func (r *Response) Text() string {
	// defer r.response.Body.Close()

	// content, err := io.ReadAll(r.response.Body)
	// if err != nil {
	// 	panic(err)
	// }

	// return string(content)

	defer r.response.Body.Close()

	utf8Reader, err := charset.NewReader(r.response.Body, r.response.Header.Get("Content-Type"))
	if err != nil {
		panic(err)
	}

	content, err := io.ReadAll(utf8Reader)
	if err != nil {
		panic(err)
	}

	return string(content)
}

func (r *Response) Cookies() []*http.Cookie {
	return r.response.Cookies()
}
