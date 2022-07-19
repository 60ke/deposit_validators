package vRpc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	rpcErr RpcErr
	ret    Results
)

type Results struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result"`
}

func GetAbi() (abi.ABI, error) {
	vABI, err := abi.JSON(strings.NewReader(ABI))
	if err != nil {

		rpcErr = RpcErr{
			Code: 501,
			Msg:  err.Error(),
		}
		Error(rpcErr)
		return vABI, rpcErr
	}
	return vABI, nil
}

func post(url string, payload *strings.Reader) ([]byte, error) {

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		rpcErr = RpcErr{
			Code: 502,
			Msg:  err.Error(),
		}
		Error(rpcErr)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		rpcErr = RpcErr{
			Code: 503,
			Msg:  err.Error(),
		}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		rpcErr = RpcErr{
			Code: 504,
			Msg:  err.Error(),
		}
		Error(rpcErr)
		return nil, rpcErr
	}
	return body, nil
}

func GetCoin(url string) (string, error) {

	payload := strings.NewReader(`{"jsonrpc":"2.0","method":"eth_coinbase", "id":1}`)
	b, err := post(url, payload)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(b, &ret)
	if err != nil {
		rpcErr = RpcErr{
			Code: 506,
			Msg:  err.Error(),
		}
		Error(rpcErr)
		return "", rpcErr
	}
	return ret.Result, nil
}

func TrustDeposit(addrs []common.Address, url string, addr string, justUpdate bool) (string, error) {
	Info("start deposit")
	vABI, err := GetAbi()
	if err != nil {
		Error(err)
		return "", err
	}
	from, err := GetCoin(url)
	if err != nil {
		Error(err)
		return "", err
	}
	fmt.Println("from:", from)
	trustData, _ := vABI.Pack("trustdeposit", addrs, justUpdate)
	data := `{"jsonrpc":"2.0","method": "eth_sendTransaction", "params": [{"from": "%s", "to": "%s", "data": "%s"}], "id": 1}`

	payload := fmt.Sprintf(data, from, addr, (hexutil.Bytes)(trustData).String())

	b, err := post(url, strings.NewReader(payload))
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(b, &ret)
	if err != nil {
		rpcErr = RpcErr{
			Code: 506,
			Msg:  err.Error(),
		}
		Error(rpcErr)
		return "", rpcErr
	}
	Info(ret)
	return ret.Result, nil
}

// func GetValidators(url string, addr string) string {
// 	vABI := GetAbi()
// 	from := GetCoin(url)
// 	trustData, _ := vABI.Pack("getValidators")
// 	data := `{"jsonrpc":"2.0","method": "eth_sendTransaction", "params": [{"from": "%s", "to": "%s", "data": "%s"}], "id": 1}`
// 	data = fmt.Sprintf(from, addr, (hexutil.Bytes)(trustData))
// 	return post(url, strings.NewReader(data))
// }

func GetBlockNumber(url string) (string, error) {
	payload := strings.NewReader(`{"jsonrpc":"2.0","method":"eth_blockNumber", "id":1}`)
	b, err := post(url, payload)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(b, &ret)
	if err != nil {
		rpcErr = RpcErr{
			Code: 506,
			Msg:  err.Error(),
		}
		Error(rpcErr)
		return "", rpcErr
	}
	return ret.Result, nil
}
