package vRpc

import (
	"fmt"
	"regexp"
	"testing"
)

func parseErr(err string) (string, string) {
	regCode := regexp.MustCompile(`code:(.*?),msg`)
	regMsg := regexp.MustCompile(`msg:(.+)`)
	code := regCode.FindStringSubmatch(err)[1]
	fmt.Println(regMsg.FindStringSubmatch(err))
	msg := regMsg.FindStringSubmatch(err)[1]
	return code, msg
}

func TestParse(t *testing.T) {
	var inputs = []RpcErr{
		{Code: 501, Msg: "json err"},
		{Code: 502, Msg: "err 502"},
		{Code: 503, Msg: "err\\503"},
		{Code: 504, Msg: "err"},
		{Code: 505, Msg: "err"},
		{Code: 506, Msg: "err"},
		{Code: 507, Msg: "err"},
		{Code: 508, Msg: "err"},
		{Code: 508, Msg: "err"},
		{Code: 510, Msg: "err"},
	}
	for _, input := range inputs {
		code, msg := parseErr(input.Error())
		fmt.Println(code, msg)
	}
}
