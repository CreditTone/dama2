package dama2

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

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

type Dama2Client struct {
	client         *http.Client
	key            string
	auth           string
	lastUpdateAuth time.Time
	lock           *sync.Mutex
}

func (p *Dama2Client) updateAuth(a string) {
	if len(a) > 0 {
		p.auth = a
		p.lastUpdateAuth = time.Now()
	}
}

func NewDama2Client(k string) *Dama2Client {
	return &Dama2Client{
		client:         &http.Client{},
		key:            k,
		auth:           "",
		lastUpdateAuth: time.Now(),
		lock:           &sync.Mutex{},
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

func (p *Dama2Client) Init(appid, uname, pwd string) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	var err error
	_, err = p.PreAuth()
	if err != nil {
		return err
	}
	_, err = p.Login(appid, uname, pwd, p.auth)
	if err != nil {
		return err
	}
	return nil
}

func (p *Dama2Client) Captcha(img []byte, codeType int, appid, uname, pwd string) (string, error) {
	if len(p.auth) == 0 || (time.Now().Sub(p.lastUpdateAuth).Minutes() > 5) {
		p.Init(appid, uname, pwd)
	}
	decode, err := p.Decode(img, p.auth, codeType, 0)
	if err != nil {
		return "", err
	}
	tc := time.NewTicker(time.Second * 2)
	n := 0
	for _ = range tc.C {
		n += 1
		if n > 5 {
			tc.Stop()
		}
		ret, err := p.GetResult(p.auth, decode.Id)
		if err == nil && ret.Ret == "0" {
			tc.Stop()
			return ret.Result, nil
		}
	}
	return "", nil
}
