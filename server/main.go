package main

import (
	"net/http"
	"regexp"
	"vRpc"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

type UpdateValidatorsReq struct {
	RpcUrl        string   `json:"rpcUrl"`
	JustUpdate    bool     `json:"justUpdate"`
	ContractAddr  string   `json:"contractAddr"`
	CurValidators []string `json:"curValidators"`
	AccessToken   string   `json:"token"`
}

type GetBlockReq struct {
	RpcUrl string `json:"rpcUrl"`
}

//   返回：
// type ReturnInfo struct {
// 	InvokeResultCode    int         `json:"invokeResultCode"`    // 返回码
// 	InvokeResultMessage string      `json:"invokeResultMessage"` // 返回码描述
// 	Data                interface{} `json:"data"`                // 协议数据，txid
// }

func checkToken(token string) bool {
	return token == "3D3781351A3EE9E4"
}

func parseErr(err string) (string, string) {
	regCode := regexp.MustCompile(`code:(.*?),msg`)
	regMsg := regexp.MustCompile(`msg:(.+)`)
	code := regCode.FindStringSubmatch(err)[1]
	msg := regMsg.FindStringSubmatch(err)[1]
	return code, msg
}

func trustDeposit(c *gin.Context) {
	var payload UpdateValidatorsReq
	if err := c.ShouldBindJSON(&payload); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{
			"invokeResultCode":    400,
			"invokeResultMessage": err.Error(),
			"data":                "",
		})
		return
	}

	var code string
	var msg string

	var addrs []common.Address
	for _, addr := range payload.CurValidators {
		address := common.HexToAddress(addr)
		addrs = append(addrs, address)
	}
	if checkToken(payload.AccessToken) {
		tx, err := vRpc.TrustDeposit(addrs, payload.RpcUrl, payload.ContractAddr, payload.JustUpdate)
		if err != nil {
			code, msg = parseErr(err.Error())
		} else {
			code = "200"
			msg = "success"
		}

		c.JSON(http.StatusOK, gin.H{
			"invokeResultCode":    code,
			"invokeResultMessage": msg,
			"data":                tx,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"invokeResultCode":    510,
			"invokeResultMessage": "token error",
			"data":                "",
		})
	}

}

func getBlockNumber(c *gin.Context) {
	var code, msg string
	var payload GetBlockReq
	if err := c.ShouldBindJSON(&payload); err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		c.JSON(http.StatusBadRequest, gin.H{
			"invokeResultCode":    400,
			"invokeResultMessage": err.Error(),
			"data":                "",
		})
		return
	}
	number, err := vRpc.GetBlockNumber(payload.RpcUrl)
	if err != nil {
		code, msg = parseErr(err.Error())
	} else {
		code = "200"
		msg = "success"
	}

	c.JSON(http.StatusOK, gin.H{
		"invokeResultCode":    code,
		"invokeResultMessage": msg,
		"data":                number,
	})

}

func main() {

	r := gin.Default()
	r.POST("/remote/sendValidatorTx", trustDeposit)
	r.POST("/remote/getBlockNumber", getBlockNumber)
	r.Run(":2345") // listen and serve on 0.0.0.0:8080

}
