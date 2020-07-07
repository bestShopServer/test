package pay

import (
	"DetectiveMasterServer/config"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/util"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Func: Get OpenId From Tencent
// 参考资料 https://pay.weixin.qq.com/wiki/doc/api/wxa/wxa_api.php?chapter=9_1
func WeChatPayParam(appId, order string, amt int) (res model.MarketPayParams, params model.PayParams, err error) {
	WxPayUrl := "https://api.mch.weixin.qq.com/pay/unifiedorder"
	util.Debug("URL[%v]", WxPayUrl)
	var tmpReq model.WxUnifiedPayXml

	tmpReq.NotifyUrl = config.GetConfig().NotifyUrl
	tmpReq.AppId = config.GetConfig().AppId
	tmpReq.OpenId = appId
	tmpReq.MchId = config.GetConfig().MchId
	tmpReq.Body = "推理大师"
	tmpReq.NonceStr = util.GetRandomString(30)
	//tmpReq.OutTradeNo	= 	GenOrder(info.ShopId)
	tmpReq.OutTradeNo = order
	tmpReq.TradeType = "JSAPI"
	tmpReq.SpbillCreateIp = config.GetConfig().PayIp
	//金额判断，后面修改判断与订单金额一致才行
	if amt < 1 {
		util.Error("金额有误[%v]", amt)
		return res, params, errors.New("金额有误")
	}
	//strAmt := strconv.FormatFloat(amt, 'f', 2, 64)
	//strs := strings.ReplaceAll(strAmt, ".", "")
	//tmpReq.TotalFee, _ = strconv.Atoi(strs)
	tmpReq.TotalFee = amt
	//sign
	stringA := fmt.Sprintf("appid=%v&body=%v&mch_id=%v&nonce_str=%v&notify_url=%v"+
		"&openid=%v&out_trade_no=%v&spbill_create_ip=%v&total_fee=%v&trade_type=%v",
		tmpReq.AppId, tmpReq.Body, tmpReq.MchId,
		tmpReq.NonceStr, tmpReq.NotifyUrl, tmpReq.OpenId, tmpReq.OutTradeNo,
		tmpReq.SpbillCreateIp, tmpReq.TotalFee, tmpReq.TradeType,
	)
	util.Debug("stringA[%v]", stringA)

	key := config.GetConfig().MchKey
	stringSignTemp := stringA + "&key=" + key
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(stringSignTemp))
	cipherStr := md5Ctx.Sum(nil)
	tmpReq.Sign = strings.ToUpper(hex.EncodeToString(cipherStr))
	util.Info("sign[%v]", tmpReq.Sign)

	output, err := xml.MarshalIndent(tmpReq, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	util.Info("string[%v]", string(output))

	client := &http.Client{}
	req, err := http.NewRequest("POST",
		WxPayUrl, strings.NewReader(string(output)))
	if err != nil {
		// handle error
		util.Error("ERROR[%v]", err.Error())
	}

	req.Header.Set("Content-Type", "text/xml")
	resp, err := client.Do(req)

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		// handle error
		util.Error("ERROR[%v]", err.Error())
		return res, params, err
	}
	util.Info(string(body))

	//准备解析返回数据
	var tmpRtn model.WxUnifiedReturn
	err = xml.Unmarshal(body, &tmpRtn)
	if err != nil {
		// handle error
		util.Error("ERROR[%v]", err.Error())
		return res, params, err
	}

	util.Debug("响应信息[%v]", tmpRtn)
	if tmpRtn.ReturnCode != "SUCCESS" {
		util.Error("获取参数失败[%v]", tmpRtn.ReturnMsg)
		return res, params, errors.New(tmpRtn.ReturnMsg)
	}
	if tmpRtn.ResultCode != "SUCCESS" {
		util.Error("获取参数失败[%v]", tmpRtn.ErrCodeDes)
		return res, params, errors.New(tmpRtn.ErrCodeDes)
	}

	//准备返回给小程序端的参数
	res.AppId = tmpReq.AppId
	res.Package = "prepay_id=" + tmpRtn.PrepayId
	res.NonceStr = util.GetRandomString(30)
	res.TimeStamp = fmt.Sprintf("%v", time.Now().Unix())
	res.SignType = "MD5"

	stringB := fmt.Sprintf("appId=%v&nonceStr=%v&package=%v&signType=%v&timeStamp=%v",
		res.AppId, res.NonceStr, res.Package, res.SignType, res.TimeStamp,
	)
	util.Debug("stringB[%v]", stringB)

	stringResSignTemp := stringB + "&key=" + key
	md5CtxRes := md5.New()
	md5CtxRes.Write([]byte(stringResSignTemp))
	cipherResStr := md5CtxRes.Sum(nil)
	res.PaySign = strings.ToUpper(hex.EncodeToString(cipherResStr))
	util.Info("sign[%v]", res.PaySign)

	//准备登记数据库参数
	params.WxUnifiedReturn = tmpRtn
	params.OutTradeNo = tmpReq.OutTradeNo
	params.PayAmt = float64(amt / 100)
	params.TotalFee = tmpReq.TotalFee
	params.Content = stringA
	params.PayContent = stringB

	return res, params, err
}
