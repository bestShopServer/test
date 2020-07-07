package pay

import (
	"DetectiveMasterServer/global"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/release"
	"DetectiveMasterServer/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

//接收微信支付回调
func PostWxNotifyRrl(c *gin.Context) {
	util.Info("接收微信支付回调...")

	var unionId, scriptName string
	//准备响应信息
	rsp := model.XmlResponse{}
	rsp.ReturnCode = "SUCCESS"

	var jparm model.WxPayNotifyRrlReq
	err := c.ShouldBindXML(&jparm)
	if err != nil {
		util.Error("解析回调失败[%v],参数:%+v", err.Error(), jparm)
		rsp.ReturnCode = "FAIL"
		rsp.ReturnMsg = "参数格式校验错误"
	}
	util.Info("商行订单号[%v]", jparm.OutTradeNo)

	taskRequest := make(map[string]interface{})
	taskRequest["OutTradeNo"] = jparm.OutTradeNo
	taskRequest["ReturnCode"] = jparm.ReturnCode
	taskRequest["ReturnMsg"] = jparm.ReturnMsg
	taskRequest["ResultCode"] = jparm.ResultCode
	taskRequest["ErrCode"] = jparm.ErrCode
	taskRequest["ErrCodeDes"] = jparm.ErrCodeDes
	taskRequest["TransactionId"] = jparm.TransactionId
	taskRequest["TimeEnd"] = jparm.TimeEnd

	util.Debug("record pay notify url taskRequest:", taskRequest)
	dbResult, err := global.Task.TaskJson(global.NewDBRequest("db.PayNotifyUrl", taskRequest))
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		release.ResponseErr(model.ERR_TASK_JSON, c)
		rsp.ReturnCode = "FAIL"
		rsp.ReturnMsg = "参数格式校验错误"
		util.Info("处理微信支付回调失败")
		c.XML(http.StatusOK, rsp)
		return
	}

	dbcode, dbparams := global.UnwrapObjectPackage(dbResult)
	util.Debug("UnwrapObjectPackage:", dbcode, dbparams)
	switch dbcode {
	case global.ERR_DB_OK:
		rsp.ReturnCode = "SUCCESS"
		rsp.ReturnMsg = "处理成功"
		util.Info("处理微信支付回调成功")

		//同步到公众号，临时需求，后期再删除
		dbResult, err := global.Task.TaskJson(global.NewDBRequest("db.OrderInfoBase", taskRequest))
		if err != nil {
			util.Error("ERROR:%v", err.Error())
		} else {
			dbcode, dbparams := global.UnwrapObjectPackage(dbResult)
			util.Debug("UnwrapObjectPackage:", dbcode, dbparams)
			switch dbcode {
			case global.ERR_DB_OK:
				if dbparams["UnionId"] != nil {
					unionId = dbparams["UnionId"].(string)
				}
				if dbparams["ScriptName"] != nil {
					scriptName = dbparams["ScriptName"].(string)
				}
				util.Debug("剧本:%v 用户:%v", scriptName, unionId)
				if len(unionId) > 0 && len(scriptName) > 0 {
					err = OrderStatePublic(unionId, scriptName)
					if err != nil {
						util.Error("ERROR:%v", err)
					}
				}
			default:
				util.Error("未查询到数据:%+v", dbparams)
			}
		}

	default:
		//测试都返回错误
		rsp.ReturnCode = "FAIL"
		rsp.ReturnMsg = "参数格式校验错误"
		util.Info("处理微信支付回调失败")
	}
	util.Info("响应信息:%+v", rsp)

	c.XML(http.StatusOK, rsp)
}
