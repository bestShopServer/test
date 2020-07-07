package pay

import (
	"DetectiveMasterServer/release"
	"DetectiveMasterServer/util"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
)

/*
https://mp.tuilidashi.xin/tuilidashitemp/index.php/Home/Dist/syncOrder 同步订单，请求之后在公众号就算支付过了
https://mp.tuilidashi.xin/tuilidashitemp/index.php/Home/Dist/getOrderStatus 查询用户支付状态
post 参数一样，scriptName 中文不编码，unionId
*/

type OrderPublicResp struct {
	Err int    `json:"err"`
	Msg string `json:"msg"`
}

//同步公众号订单
func OrderSyncPublic(unionId, scriptName string) (err error) {
	util.Info("OrderSyncPublic...")
	var resp OrderPublicResp
	res, err := http.PostForm("https://mp.tuilidashi.xin/tuilidashitemp/index.php/Home/Dist/syncOrder",
		url.Values{
			"unionId":    {unionId},
			"scriptName": {scriptName},
		})
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return err
	}

	if res.StatusCode != 200 {
		util.Error("Unexpected status code:%v", res.StatusCode)
		return err
	}

	// Read the token out of the response body
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, res.Body)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return err
	}

	err = json.Unmarshal(buf.Bytes(), &resp)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return err
	}
	util.Debug("buf:%+v", resp)
	if resp.Err != 0 {
		return errors.New(resp.Msg)
	}

	return err
}

//测试同步订单信息到公众号服务
func TestOrderSyncPublic(c *gin.Context) {

	unionId := "o9XY71VQa96OWPktAFgdSgxQvzn"
	scriptName := "第十二夜"
	err := OrderSyncPublic(unionId, scriptName)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		release.ResponseErrMsg(c, err.Error())
		return
	}

	release.ResponseSuccess(c)
}
