package httpx

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Request struct {
	method string
	url string

	params map[string]string

	body *bytes.Buffer
}

func (r *Request) Do() (*Response, error) {
	rq, err := http.NewRequest(r.method, r.url, r.body)
	if err != nil {
		return nil, err
	}

	q := rq.URL.Query()
	for k, v := range r.params {
		q.Add(k, v)
	}

	rq.URL.RawQuery = q.Encode()

	_response, err := http.DefaultClient.Do(rq)
	if err != nil {
		return nil, err
	}

	response := &Response{
		response: _response,
	}

	return response, nil
}

func (r *Request) Params(params map[string]interface{}) *Request {
	r.params = make(map[string]string)

	for k, v := range params {
		r.params[k] = v.(string)
	}

	return r
}


func (r *Request) Json(_json map[string]interface{}) *Request {
	marshalled, err := json.Marshal(_json)
	if err != nil {
		panic(err)
	}

	r.body = bytes.NewBuffer(marshalled)

	return r
}