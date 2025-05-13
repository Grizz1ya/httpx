package httpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Request struct {
	client        *http.Client
	headers       map[string]string
	staticHeaders map[string]string

	method string
	url    string

	params map[string]string
	body   *bytes.Buffer
}

func (r *Request) Do() (*Response, error) {
	// * if body is nil, we should create a new buffer
	if r.body == nil {
		r.body = bytes.NewBuffer([]byte{})
	}

	rq, err := http.NewRequest(r.method, r.url, r.body)
	if err != nil {
		return nil, err
	}

	for k, v := range r.headers {
		rq.Header.Add(k, v)
	}

	for k, v := range r.staticHeaders {
		rq.Header.Add(k, v)
	}

	q := rq.URL.Query()
	for k, v := range r.params {
		q.Add(k, v)
	}

	rq.URL.RawQuery = q.Encode()

	// * if client is nil, we should use the default client
	var client *http.Client
	if r.client != nil {
		client = r.client
	} else {
		client = &http.Client{}
	}

	_response, err := client.Do(rq)
	if err != nil {
		return nil, err
	}

	response := &Response{
		response:   _response,
		StatusCode: _response.StatusCode,
		URL:        _response.Request.URL.String(),
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

func (r *Request) Data(data map[string]interface{}) *Request {
	formData := bytes.NewBufferString("")
	for key, value := range data {
		if formData.Len() > 0 {
			formData.WriteString("&")
		}
		formData.WriteString(key + "=" + fmt.Sprintf("%v", value))
	}
	r.body = formData

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

func (r *Request) Body(body []byte) *Request {
	r.body = bytes.NewBuffer(body)
	return r
}

func (r *Request) Headers(headers map[string]string) *Request {
	r.headers = headers
	return r
}
