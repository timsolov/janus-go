package janus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

type APIRequest interface {
	ActionName() string
	Endpoint() string
	Payload() map[string]interface{}
}

type BaseRequest struct {
	Action      string
	Transaction string
	Secret      string
}

func (r *BaseRequest) ActionName() string {
	return r.Action
}

func (r *BaseRequest) Endpoint() string {
	return ""
}

func (r *BaseRequest) Payload() map[string]interface{} {
	m := make(map[string]interface{})
	m["janus"] = r.Action
	m["transaction"] = r.Transaction
	m["admin_secret"] = r.Secret
	return m
}

type TokenRequest struct {
	BaseRequest
	Token   string
	Plugins []string
}

func (r *TokenRequest) Payload() map[string]interface{} {
	m := r.BaseRequest.Payload()
	m["token"] = r.Token
	m["plugins"] = r.Plugins
	return m
}

type SessionRequest struct {
	BaseRequest
	SessionID uint64
}

func (r *SessionRequest) Payload() map[string]interface{} {
	m := r.BaseRequest.Payload()
	m["session_id"] = r.SessionID
	return m
}

func (r *SessionRequest) Endpoint() string {
	return fmt.Sprintf("/%d", r.SessionID)
}

type HandleRequest struct {
	SessionRequest
	HandleID uint64
}

func (r *HandleRequest) Endpoint() string {
	return fmt.Sprintf("%s/%d", r.SessionRequest.Endpoint(), r.HandleID)
}

func (r *HandleRequest) Payload() map[string]interface{} {
	m := r.SessionRequest.Payload()
	m["handle_id"] = r.HandleID
	return m
}

type Transport interface {
	Request(APIRequest) (interface{}, error)
}

type TransportError struct {
	Code int
	Msg  string
}

func (e *TransportError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Msg)
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

	pResp, err := ParseAMResponse(r.ActionName(), body)
	if err != nil {
		return nil, err
	}

	switch pResp := pResp.(type) {
	case *ErrorAMResponse:
		return nil, pResp
	default:
		return pResp, nil
	}
}

type AdminAPI interface {
	AddToken(token string, plugins []string) (interface{}, error)
	AllowToken(token string, plugins []string) (interface{}, error)
	DisallowToken(token string, plugins []string) (interface{}, error)
	RemoveToken(token string) (interface{}, error)
	ListTokens() (interface{}, error)

	ListSessions() (interface{}, error)

	HandleInfo(sessionID, handleID uint64) (interface{}, error)
}

type AdminAPIImpl struct {
	transport Transport
	secret    string
}

func NewAdminAPI(url, secret string) (*AdminAPIImpl, error) {
	api := new(AdminAPIImpl)
	api.secret = secret

	if strings.HasPrefix(url, "http") {
		api.transport = NewHttpTransport(url)
	} else {
		return nil, fmt.Errorf("unsupported transport for %s", url)
	}

	return api, nil
}

func (api *AdminAPIImpl) AddToken(token string, plugins []string) (interface{}, error) {
	return api.transport.Request(api.makeTokenRequest("add_token", token, plugins))
}

func (api *AdminAPIImpl) AllowToken(token string, plugins []string) (interface{}, error) {
	return api.transport.Request(api.makeTokenRequest("allow_token", token, plugins))
}

func (api *AdminAPIImpl) DisallowToken(token string, plugins []string) (interface{}, error) {
	return api.transport.Request(api.makeTokenRequest("disallow_token", token, plugins))
}

func (api *AdminAPIImpl) RemoveToken(token string) (interface{}, error) {
	return api.transport.Request(api.makeTokenRequest("remove_token", token, nil))
}

func (api *AdminAPIImpl) ListTokens() (interface{}, error) {
	return api.transport.Request(api.makeBaseRequest("list_tokens"))
}

func (api *AdminAPIImpl) ListSessions() (interface{}, error) {
	return api.transport.Request(api.makeBaseRequest("list_sessions"))
}

func (api *AdminAPIImpl) HandleInfo(sessionID, handleID uint64) (interface{}, error) {
	return api.transport.Request(api.makeHandleRequest("handle_info", sessionID, handleID))
}

func (api *AdminAPIImpl) makeBaseRequest(action string) *BaseRequest {
	return &BaseRequest{
		Action:      action,
		Transaction: RandString(12),
		Secret:      api.secret,
	}
}

func (api *AdminAPIImpl) makeTokenRequest(action, token string, plugins []string) *TokenRequest {
	return &TokenRequest{
		BaseRequest: *api.makeBaseRequest(action),
		Token:       token,
		Plugins:     plugins,
	}
}

func (api *AdminAPIImpl) makeSessionRequest(action string, sessionID uint64) *SessionRequest {
	return &SessionRequest{
		BaseRequest: *api.makeBaseRequest(action),
		SessionID:   sessionID,
	}
}

func (api *AdminAPIImpl) makeHandleRequest(action string, sessionID, handleID uint64) *HandleRequest {
	return &HandleRequest{
		SessionRequest: *api.makeSessionRequest(action, sessionID),
		HandleID:       handleID,
	}
}
