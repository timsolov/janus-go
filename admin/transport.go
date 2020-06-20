package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type TransportError struct {
	Code int
	Msg  string
}

func (e *TransportError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Msg)
}

type Transport interface {
	Request(APIRequest) (interface{}, error)
	Close() error
}

type HttpTransport struct {
	client *http.Client
	url    string
}

func NewHttpTransport(url string) *HttpTransport {
	c := new(HttpTransport)
	c.url = url
	c.client = &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 5 * time.Second,
		},

		Timeout: 10 * time.Second,
	}
	return c
}

func (t *HttpTransport) Request(r APIRequest) (interface{}, error) {
	b, err := json.Marshal(r.Payload())
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, t.url+r.Endpoint(), bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, &TransportError{Code: resp.StatusCode, Msg: resp.Status}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	pResp, err := ParseAMResponse(r, body)
	if err != nil {
		return nil, err
	}

	switch pResp := pResp.(type) {
	case error:
		return nil, pResp
	default:
		return pResp, nil
	}
}
func (t *HttpTransport) Close() error {
	return nil
}
