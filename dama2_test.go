package dama2

import (
	"crypto/md5"
	"fmt"
	"testing"
)

func TestEncode(t *testing.T) {
	uname := "user"
	pwd := "test"
	preauth := "1234567890abcdef1234567890abcdef"
	p := NewDama2Client("9503ce045ad14d83ea876ab578bd3184")
	src := fmt.Sprintf("%s\n%s\n%s", preauth, uname, fmt.Sprintf("%x", md5.Sum([]byte(pwd))))
	enc, _ := DesEncrypt([]byte(src), p.getKey8())
	if fmt.Sprintf("%x", enc) != "a733506fda6e182300d34a1bcca568d8a733506fda6e182300d34a1bcca568d8f6368183d9f688b7a43c783d0184850751f1d013925667735ddd457beaa8da8010b74e6f64b15c4d" {
		t.Error()
		return
	}
}
