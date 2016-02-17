/*
功能URL：http://api.dama2.com:7788/app/decode
输入参数：auth、type、len(可选)、timeout（可选）、文件数据
返回数据：ret，desc，aut，id(用于查询结果和报告结果)
*/
package dama2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"
)

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func createPart(fieldname, filename, contentType string, w *multipart.Writer, value []byte) error {
	h := make(textproto.MIMEHeader)
	if len(filename) == 0 {
		h.Set("Content-Disposition",
			fmt.Sprintf(`form-data; name="%s"`,
				escapeQuotes(fieldname)))
	} else {
		h.Set("Content-Disposition",
			fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
				escapeQuotes(fieldname), escapeQuotes(filename)))
	}
	h.Set("Content-Type", contentType)
	wt, err := w.CreatePart(h)
	if err != nil {
		return err
	}
	_, err = wt.Write(value)
	return err
}

func uploadBody(ubody []byte, link string, params map[string]string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for k, v := range params {
		createPart(k, "", "plain/text", writer, []byte(v))
	}
	err := createPart("data", "test.png", "image/png", writer, ubody)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", link, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func (p *Dama2Client) Decode(body []byte, auth string, codeType, length int) (*Result, error) {
	params := map[string]string{
		"auth": auth,
		"type": strconv.Itoa(codeType),
	}
	if length > 0 {
		params["len"] = strconv.Itoa(length)
	}
	req, err := uploadBody(body, "http://api.dama2.com:7788/app/decode", params)

	if err != nil {
		return nil, err
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	var ret Result
	err = json.Unmarshal(b, &ret)
	p.updateAuth(ret.Auth)
	return &ret, err
}

/*
功能URL：http://api.dama2.com:7788/app/getResult
输入参数：auth，id
返回数据：ret，desc，result，cookie, auth
*/
func (p *Dama2Client) GetResult(auth, id string) (*Result, error) {
	resp, err := p.client.Get("http://api.dama2.com:7788/app/getResult?auth=" + auth + "&id=" + id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	var ret Result
	err = json.Unmarshal(b, &ret)
	p.updateAuth(ret.Auth)
	return &ret, err
}
