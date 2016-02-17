package dama2

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
)

type Dama2Client struct {
	client *http.Client
	key    string
}

type Result struct {
	Id     string `json:"id"`
	Ret    string `json:"ret"`
	Desc   string `json:"desc"`
	Auth   string `json:"auth"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	QQ     string `json:"qq"`
	Result string `json:"result"`
}

func (p *Result) String() string {
	b, _ := json.Marshal(p)
	return string(b)
}

func NewDama2Client(k string) *Dama2Client {
	return &Dama2Client{
		client: &http.Client{},
		key:    k,
	}
}

func (p *Dama2Client) getKey8() []byte {
	key16, _ := hex.DecodeString(p.key)
	key8 := make([]byte, 8, 8)
	for i := 0; i < 8; i++ {
		key8[i] = (byte)((key16[i] ^ key16[i+8]) & 0xff)
	}
	return key8
}
