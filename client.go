package dama2

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

var Instance *Dama2Client

type Dama2Client struct {
	client         *http.Client
	key            string
	auth           string
	lastUpdateAuth time.Time
	lock           *sync.Mutex
	uname          string
	pwd            string
	appid          string
}

func (p *Dama2Client) updateAuth(a string) {
	if len(a) > 0 {
		p.auth = a
		p.lastUpdateAuth = time.Now()
	}
}

func NewDama2Client(k, uname, pwd, appid string) *Dama2Client {
	return &Dama2Client{
		client:         &http.Client{},
		key:            k,
		auth:           "",
		lastUpdateAuth: time.Now(),
		lock:           &sync.Mutex{},
		uname:          uname,
		pwd:            pwd,
		appid:          appid,
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

func (p *Dama2Client) Init() error {
	p.lock.Lock()
	defer p.lock.Unlock()
	if len(p.auth) == 0 || time.Now().Sub(p.lastUpdateAuth).Minutes() > 5 {
		var err error
		_, err = p.PreAuth()
		if err != nil {
			return err
		}
		_, err = p.Login(p.appid, p.uname, p.pwd, p.auth)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Dama2Client) Query(img []byte, codeType, format string) (string, error) {
	if len(p.auth) == 0 || (time.Now().Sub(p.lastUpdateAuth).Minutes() > 5) {
		p.Init()
	}
	codeTypeInt, _ := strconv.Atoi(codeType)
	fmt.Println(p.auth)
	decode, err := p.Decode(img, format, p.auth, codeTypeInt, 0)
	if err != nil {
		return "", err
	}
	tc := time.NewTicker(time.Second * 2)
	n := 0
	for _ = range tc.C {
		n += 1
		if n > 5 {
			tc.Stop()
			break
		}
		ret, err := p.GetResult(p.auth, decode.Id)
		if err == nil && ret.Ret == "0" {
			tc.Stop()
			return ret.Result, nil
		}
	}
	return "", nil
}
