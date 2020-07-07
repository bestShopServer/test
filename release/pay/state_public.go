package pay

import (
	"DetectiveMasterServer/release"
	"DetectiveMasterServer/util"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

/*
https://mp.tuilidashi.xin/tuilidashitemp/index.php/Home/Dist/syncOrder 同步订单，请求之后在公众号就算支付过了
https://mp.tuilidashi.xin/tuilidashitemp/index.php/Home/Dist/getOrderStatus 查询用户支付状态
post 参数一样，scriptName 中文不编码，unionId
*/

//同步公众号订单
func OrderStatePublic(unionId, scriptName string) (err error) {
	util.Info("OrderStatePublic ...")
	var response OrderPublicResp
	url := "https://mp.tuilidashi.xin/tuilidashitemp/index.php/Home/Dist/getOrderStatus"
	//res, err := http.PostForm(url,
	//	url.Values{
	//		"unionId":    {unionId},
	//		"scriptName": {scriptName},
	//	})
	//if err != nil {
	//	util.Error("ERROR:%v", err.Error())
	//	return err
	//}

	//// Read the token out of the response body
	//buf := new(bytes.Buffer)
	//_, err = io.Copy(buf, res.Body)
	//if err != nil {
	//	util.Error("ERROR:%v", err.Error())
	//	return err
	//}

	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	w.WriteField("unionId", unionId)
	w.WriteField("scriptName", scriptName)
	w.Close()
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, _ := http.DefaultClient.Do(req)
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	util.Debug("code:%v", resp.StatusCode)
	util.Debug("%s", data)

	if resp.StatusCode != 200 {
		util.Error("Unexpected status code:%v", resp.StatusCode)
		return err
	}

	err = json.Unmarshal(data, &response)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return err
	}
	util.Debug("buf:%+v", response)
	if response.Err != 0 {
		return errors.New(response.Msg)
	}

	return err
}

//测试同步订单信息到公众号服务
func TestOrderStatePublic(c *gin.Context) {

	unionId := "o9XY71VQa96OWPktAFgdSgxQvzn"
	scriptName := "第十二夜"
	err := OrderStatePublic(unionId, scriptName)
	if err != nil {
		util.Error("ERROR:%v", err)
		release.ResponseErrMsg(c, err.Error())
		//release.ResponseErr(model.ERR_SCRIPT_ORDER_SYNC, c)
		return
	}

	release.ResponseSuccess(c)
}
