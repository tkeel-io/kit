package auth

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/tkeel-io/kit"
	"github.com/tkeel-io/kit/log"
	transportHTTP "github.com/tkeel-io/kit/transport/http"
)

const (
	AuthTokenURLTestRemote string = "http://192.168.123.9:30707/apis/security/v1/oauth/authenticate"
	AuthTokenURLInvoke     string = "http://localhost:3500/v1.0/invoke/keel/method/apis/security/v1/oauth/authenticate"

	_Authorization string = "Authorization"
)

var _auth = AuthTokenURLInvoke

type User struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
	Token    string `json:"token"`
}

func Authenticate(token interface{}, urls ...string) (*User, error) {
	if len(urls) > 0 {
		_auth = urls[0]
	}

	tokenStr := ""
	switch t := token.(type) {
	case string:
		tokenStr = t
	case context.Context:
		val, ok := transportHTTP.HeaderFromContext(t)[_Authorization]
		if !ok || len(tokenStr) == 0 {
			return nil, errors.New("invalid Authenticate")
		}
		tokenStr = val[0]
	default:
		return nil, errors.New("invalid token type")
	}
	if tokenStr == "" {
		return nil, errors.New("token is empty")
	}

	req, err := http.NewRequest("GET", _auth, nil)
	if nil != err {
		return nil, err
	}
	req.Header.Add(_Authorization, tokenStr)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("error ", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := getBody(resp)
	if nil != err {
		log.Error("error parse token, ", err)
		return nil, err
	}

	var response kit.Response
	if err = json.Unmarshal(body, &response); nil != err {
		log.Error("resp Unmarshal error, ", err)
		return nil, err
	}

	if response.Code != "io.tkeel.SUCCESS" {
		return nil, errors.New(response.Msg)
	}

	respData, ok := response.Data.(map[string]interface{})
	if !ok {
		log.Error("resp data is not map[string]interface{}")
		return nil, errors.New("resp data is not map[string]interface{}")
	}

	id, ok := respData["user_id"].(string)
	if !ok {
		return nil, errors.New("parse token user_id data error")
	}

	tenantId, ok := respData["tenant_id"].(string)
	if !ok {
		return nil, errors.New("parse token tenant_id data error")
	}

	return &User{
		ID:       id,
		TenantID: tenantId,
		Token:    tokenStr,
	}, nil
}

func getBody(resp *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("error ReadAll", err)
		return body, err
	}

	log.Debug("receive resp, ", string(body))
	if resp.StatusCode != 200 {
		log.Error("bad status ", resp.StatusCode)
		return body, errors.New(resp.Status)
	}
	return body, nil
}
