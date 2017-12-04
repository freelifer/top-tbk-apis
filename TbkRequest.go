package tbk

import (
	"net/url"
	"strconv"
)

type ITbkRequest interface {
	GetBody() string
	GetReqUrl() string
}

type TbkRequest struct {
	reqSysMap     *TbkMap
	reqAppMap     *TbkMap
	reqUrl        string
	reqAppKey     string
	reqAppSecret  string
	reqSignMethod string
	values        url.Values
}

func (t *TbkRequest) InitDefaultValue(url, app_key, app_secret string) {
	t.reqUrl = url
	t.reqAppKey = app_key
	t.reqAppSecret = app_secret
	t.reqSignMethod = SIGN_METHOD_HMAC
}

func (t *TbkRequest) AddReqParam(key, value string) {
	if t.reqAppMap == nil {
		t.reqAppMap = NewTbkMap()
	}
	t.reqAppMap.Put(key, value)
}

func (t *TbkRequest) GetReqParma(key string) string {
	return t.reqAppMap.Get(key)
}

func (t *TbkRequest) AddSysReqParam(key, value string) {
	if t.reqSysMap == nil {
		t.reqSysMap = NewTbkMap()
	}
	t.reqSysMap.Put(key, value)
}

func (t *TbkRequest) GetSysReqParma(key string) string {
	return t.reqSysMap.Get(key)
}

func (t *TbkRequest) GetBody() string {
	return getBody(t.reqAppMap)
}

func (t *TbkRequest) GetReqUrl() string {
	return createUrl(t.reqUrl, t.reqSysMap)
}

func (t *TbkRequest) GetResponse(method string, resp interface{}) ([]byte, error) {
	// 1. protocalMustParams
	t.AddSysReqParam("method", method)
	t.AddSysReqParam("v", "2.0")
	t.AddSysReqParam("app_key", t.reqAppKey)
	t.AddSysReqParam("timestamp", currentTime())
	t.AddSysReqParam("sign_method", t.reqSignMethod)
	// sign

	// 2. protocalOptParams
	t.AddSysReqParam("session", "")
	t.AddSysReqParam("target_app_key", "")
	t.AddSysReqParam("partner_id", "")
	t.AddSysReqParam("format", "json")
	t.AddSysReqParam("simplify", "")

	// 添加签名参数
	t.AddSysReqParam("sign", signTopRequest(t.reqSysMap, t.reqAppMap, t.reqAppSecret, t.reqSignMethod))
	// 3. appParams

	return executeRequest(t, resp)
}

type TbkMap struct {
	Values map[string]string
}

func NewTbkMap() *TbkMap {
	return &TbkMap{Values: make(map[string]string)}
}

func (t *TbkMap) Put(key, value string) {
	if areNotEmpty(key, value) {
		t.Values[key] = value
	}
}

func (t *TbkMap) Get(key string) string {
	return t.Values[key]
}

func areNotEmpty(key, value string) bool {
	if len(key) == 0 {
		return false
	}

	if len(value) == 0 {
		return false
	}
	return true
}

type ErrResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type TbkErrResponse struct {
	Error            string       `json:"error"`
	ErrorDescription string       `json:"error_description"`
	ErrResponse      *ErrResponse `json:"error_response"`
}

func (t *TbkErrResponse) GetErr() string {
	if t.Error != "" {
		return "[" + t.Error + "]" + t.ErrorDescription
	}
	if t.ErrResponse != nil {
		if t.ErrResponse.Code > 0 {
			return "[" + strconv.Itoa(t.ErrResponse.Code) + "]" + t.ErrResponse.Msg
		}
	}
	return ""
}
