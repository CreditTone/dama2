/*
功能URL：http://api.dama2.com:7788/app/preauth
输入参数：无
返回数据：ret，desc，auth(预授信息，仅用于计算加密信息)
*/
package dama2

import (
	"encoding/json"
	"io/ioutil"
)

func (p *Dama2Client) PreAuth() (*Result, error) {
	var ret Result
	resp, err := p.client.Get("http://api.dama2.com:7788/app/preauth")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &ret)
	p.updateAuth(ret.Auth)
	return &ret, err
}
