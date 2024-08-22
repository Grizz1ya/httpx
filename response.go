package httpx

import (
	"io"
	"net/http"

	"github.com/Anderson-Lu/gofasion/gofasion"
)

type Response struct {
	response *http.Response
}

func (r *Response) Json() *gofasion.Fasion {
	defer r.response.Body.Close()

	content, err := io.ReadAll(r.response.Body)
	if err != nil {
		panic(err)
	}

	json := gofasion.NewFasion(string(content))

	return json
}