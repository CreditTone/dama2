package dama2

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestEncode(t *testing.T) {
	uname := "user"
	pwd := "password"
	appid := "123"
	preauth := "1234567890abcdef1234567890abcdef"
	p := NewDama2Client("2ed9f7b2466caf9ff3f8d3109ca16530", uname, pwd, appid)
	src := fmt.Sprintf("%s\n%s\n%s", preauth, uname, fmt.Sprintf("%x", md5.Sum([]byte(pwd))))
	enc, _ := DesEncrypt([]byte(src), p.getKey8())
	if fmt.Sprintf("%x", enc) != "55cdb9bd553232859bb7b0077dc51e3255cdb9bd553232859bb7b0077dc51e3221d4157f116dff5d30d57ca592186a776474a14edaf4250db7b2a6696eba226980ac6fa145ea789b9bebbda4cdc16b8f" {
		t.Error(fmt.Sprintf("%x", enc))
		return
	}
	img, err := ioutil.ReadFile("createVerifyCode.jpeg")
	if err != nil {
		t.Error(err)
	}
	format := "jpeg"
	codeType := "90"
	ret, err := p.Query(img, codeType, format)
	if err != nil {
		t.Error(err)
	}
	t.Log(ret)
}
