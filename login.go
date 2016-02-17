package dama2

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

/*
功能URL：http://api.dama2.com:7788/app/login
输入参数：appID、sname(可选参数)、encinfo
返回数据：ret，desc，auth
*/
func (p *Dama2Client) Login(appid, uname, pwd, preauth string) (*Result, error) {
	src := fmt.Sprintf("%s\n%s\n%s", preauth, uname, fmt.Sprintf("%x", md5.Sum([]byte(pwd))))
	enc, _ := DesEncrypt([]byte(src), p.getKey8())
	params := url.Values{}
	params.Add("appID", appid)
	params.Add("appid", appid)
	params.Add("encinfo", fmt.Sprintf("%x", enc))
	link := "http://api.dama2.com:7788/app/login?" + params.Encode()
	resp, err := p.client.Get(link)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	var ret Result
	err = json.Unmarshal(b, &ret)
	return &ret, nil
}

/*
功能URL：http://api.dama2.com:7788/app/readInfo
输入参数：auth
返回数据：ret，desc，auth, name, qq、email、tel
*/
func (p *Dama2Client) ReadInfo(auth string) (*Result, error) {
	resp, err := p.client.Get("http://api.dama2.com:7788/app/readInfo?auth=" + auth)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	var ret Result
	err = json.Unmarshal(b, &ret)
	return &ret, nil
}
