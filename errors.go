package vRpc

import "fmt"

type RpcErr struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (err RpcErr) Error() string {
	return fmt.Sprintf("code:%d,msg:%s", err.Code, err.Msg)
}
