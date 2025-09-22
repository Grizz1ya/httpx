package httpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/Grizz1ya/httpx/utils"
)

type Request struct {
	client        *http.Client
	headers       map[string]string
	staticHeaders map[string]string

	method string
	url    string

	params map[string]string
	body   *bytes.Buffer

	cookieOrigins *utils.CookieOriginMap
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

	if r.cookieOrigins != nil {
		for _, cookie := range _response.Cookies() {
			if cookie.Domain == "" {
				cookie.Domain = rq.URL.Host
			}

			// Handle server-side cookie removals so Session.Cookies stays in sync.
			if cookie.MaxAge <= 0 || cookie.Value == "" ||
				(!cookie.Expires.IsZero() && cookie.Expires.Before(time.Now())) {
				r.cookieOrigins.Remove(_response.Request.URL.Host, cookie)
				if cookie.Domain != "" {
					r.cookieOrigins.Remove(cookie.Domain, cookie)
				}
				continue
			}

			r.cookieOrigins.Add(_response.Request.URL.Host, cookie)
		}
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

		var valueStr string
		switch v := value.(type) {
		case map[string]interface{}:
			// Serialize nested maps as JSON
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				valueStr = fmt.Sprintf("%v", value)
			} else {
				valueStr = string(jsonBytes)
			}
		case []interface{}, []int, []float64, []string:
			// Serialize arrays/slices as JSON
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				valueStr = fmt.Sprintf("%v", value)
			} else {
				valueStr = string(jsonBytes)
			}
		default:
			valueStr = fmt.Sprintf("%v", value)
		}

		formData.WriteString(key + "=" + valueStr)
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

func request(method, url string, client *http.Client, headers map[string]string, cookieOrigins *utils.CookieOriginMap) *Request {
	if client == nil {
		jar, _ := cookiejar.New(nil)
		client = &http.Client{
			Jar:     jar,
			Timeout: 10 * time.Second,
		}
	}

	return &Request{
		method:        method,
		url:           url,
		client:        client,
		staticHeaders: headers,
		cookieOrigins: cookieOrigins,
	}
}

// * Static methods
func Get(url string) *Request {
	return request("GET", url, nil, nil, nil)
}

func Post(url string) *Request {
	return request("POST", url, nil, nil, nil)
}
