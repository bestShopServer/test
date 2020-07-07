package pay

import (
	"DetectiveMasterServer/global"
	"DetectiveMasterServer/gocache"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/release"
	"DetectiveMasterServer/util"
	"github.com/gin-gonic/gin"
	"strconv"
)

//
//func CheckUserScriptCost(req model.ScriptUserCostReq) ([]model.UserScript, int) {
//	util.Info("Check User Script Cost...")
//	var userScripts []model.UserScript
//
//	taskRequest := make(map[string]interface{})
//	taskRequest["ScriptId"] = req.ScriptId
//	taskRequest["UnionIds"] = req.UnionIds
//	taskRequest["OpenId"] = req.OpenId
//	util.Info("UnionIds[%v]", req.UnionIds)
//
//	util.Info("db.ScriptUserCost ...")
//	dbResult, err := global.Task.TaskJson(global.NewDBRequest("db.ScriptUserCost", taskRequest))
//	if err != nil {
//		util.Error("db.ScriptUserCost ERROR[%v]", err.Error())
//		return userScripts, model.ERR_TASK_JSON
//	}
//	util.Info("获取数据成功...")
//
//	//dbcode, dbparam := global.UnwrapObjectPackage(dbResult)
//	//
//	//switch dbcode {
//	//case global.ERR_DB_OK:
//	//	util.Debug("%v", dbparam)
//	//	total := int(dbparam["Total"].(float64))
//	//	if total == 0 {
//	//		userScript.CostFlag = 1 //付款标识 1 未付款 2 已付款
//	//	} else {
//	//		userScript.ScriptId = int(dbparam["ScriptId"].(float64))
//	//		userScript.UserId = int(dbparam["UserId"].(float64))
//	//		userScript.Status = int(dbparam["Status"].(float64))
//	//		userScript.EffTime = trans.StringToDate(dbparam["EffTime"].(string))
//	//		userScript.InvTime = trans.StringToDate(dbparam["InvTime"].(string))
//	//		if userScript.Status == 2 {
//	//			userScript.CostFlag = 2 //付款标识 1 未付款 2 已付款
//	//		}
//	//	}
//	//
//	//	return userScript, model.ERR_OK
//	//default:
//	//	return userScript, model.ERR_DEFAULT
//	//}
//
//	dbcode, dbparams := global.UnwrapArrayPackage(dbResult)
//
//	switch dbcode {
//	case global.ERR_DB_OK:
//		if len(dbparams) == 0 {
//			util.Info("数据数为零!")
//			return userScripts, global.ERR_DB_NOTFOUND_DATA
//		}
//		for _, p := range dbparams {
//			var tmp model.UserScript
//			pm := p.(map[string]interface{})
//			if pm["ScriptId"] != nil {
//				tmp.ScriptId = int(pm["ScriptId"].(float64))
//			}
//			if pm["UserId"] != nil {
//				tmp.UserId = int(pm["UserId"].(float64))
//			}
//			if pm["Status"] != nil {
//				tmp.Status = int(pm["Status"].(float64))
//				if tmp.Status == 2 {
//					tmp.CostFlag = 2 //付款标识 1 未付款 2 已付款
//				}
//			}
//
//			if pm["EffTime"] != nil {
//				tmp.EffTime = trans.StringToDate(pm["EffTime"].(string))
//			}
//			if pm["InvTime"] != nil {
//				tmp.InvTime = trans.StringToDate(pm["InvTime"].(string))
//			}
//			userScripts = append(userScripts, tmp)
//		}
//		return userScripts, model.ERR_OK
//	default:
//		return userScripts, model.ERR_DEFAULT
//	}
//}

//获取剧本基本信息
func GetScriptBase(sid int) (res model.Script, code int) {
	util.Info("Check User Script Cost...")

	taskRequest := make(map[string]interface{})
	taskRequest["ScriptId"] = sid

	dbResult, err := global.Task.TaskJson(global.NewDBRequest("db.ScriptBaseGet", taskRequest))
	if err != nil {
		return res, model.ERR_TASK_JSON
	}

	dbcode, dbparam := global.UnwrapObjectPackage(dbResult)

	switch dbcode {
	case global.ERR_DB_OK:
		util.Debug("%+v", dbparam)
		res.Num = int(dbparam["Num"].(float64))
		res.Price = int(dbparam["Price"].(float64))
		res.Name = dbparam["Name"].(string)

		return res, model.ERR_OK
	default:
		return res, model.ERR_DEFAULT
	}
}

// Func: wx.pay_param handler获取支付参数
func WxPayParams(c *gin.Context) {
	util.Info("Get Wechat User Pay Params ...")

	//获取请求数据
	var price int
	var jparm model.GetWxPayParamsReq
	var req model.ScriptUserCostReq
	var cost model.UserScript
	var iNum int //付款用户数
	var unionIds []string
	var rsp model.MarketPayParamsResp

	err := c.ShouldBindJSON(&jparm)
	if err != nil {
		util.Error("解析报文出错[%v] param:%+v", err.Error(), jparm)
		release.ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	util.Info("请求参数[%+v]", jparm)
	// Optional Fields List
	optionalFields := []string{"Price"}

	// Check Param
	if !release.CheckParams(jparm, "WxPayParams", err, optionalFields) {
		util.Error("参数有误")
		release.ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	util.Info("剧本ID[%v]", jparm.SId)
	//获取剧本基本信息
	scriptBase, code := GetScriptBase(jparm.SId)
	if code != model.ERR_OK {
		util.Logger(util.ERROR_LEVEL, "CheckUserScriptCost", "Get User ScriptId Err")
		release.ResponseErr(model.ERR_CHECK_SCRIPT_COST, c)
		return
	}

	orderNo := util.GenOrderNo()
	//1请大家2AA付款
	if jparm.Flag == 2 { //会员不再进行付款
		util.Info("AA付款...")
		//查询用户是否为会员
		res, err := release.GetUserMemberBase(jparm.UnionId)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			release.ResponseErr(model.ERR_DEFAULT, c)
			return
		}
		if res.Member > 0 { //用户是会员
			util.Debug("%v会员等级:%v", jparm.UnionId, res.Member)
			cost.CostFlag = 2 //设置为已付款调过付款环节

		} else {
			price = scriptBase.Price // / float64(scriptBase.Num)
			//AA付款时，查看是否需要支付，若已经购买过则无需支付
			req.ScriptId = jparm.SId
			req.OpenId = jparm.OpenId
			req.UnionIds = append(req.UnionIds, jparm.UnionId)
			costs2, code := release.CheckUserScriptCost(req)
			if code != model.ERR_OK && code != global.ERR_DB_NOTFOUND_DATA {
				util.Error("Get User ScriptId Err [%v]", code)
				release.ResponseErr(model.ERR_CHECK_SCRIPT_COST, c)
				return
			}
			if code == global.ERR_DB_NOTFOUND_DATA {
				cost.CostFlag = 1
			} else {
				cost.Status = costs2[0].Status
				cost.EffTime = costs2[0].EffTime
				cost.InvTime = costs2[0].InvTime
				cost.CostFlag = costs2[0].CostFlag
			}
			util.Debug("%+v", cost)
			if cost.CostFlag == 2 { //已付款
				util.Info("已付款，无需再次支付...")
			}
		}
		unionIds = append(unionIds, jparm.UnionId)
	} else { //包场
		util.Info("请大家玩付款...")
		req.ScriptId = jparm.SId
		req.OpenId = jparm.OpenId
		roomId := strconv.Itoa(jparm.RoomId)
		util.Info("房间号[%s]", roomId)

		// Find Room In Room Cache
		//_, ok := global.RoomCache[roomId]
		//if !ok {
		//	util.Error("房间不存在[%v]", jparm.RoomId)
		//	release.ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		//	return
		//}
		ok, err := gocache.CheckRoomExists(roomId)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			release.ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
			return
		}
		if !ok {
			util.Error("ERR_ROOM_NOT_EXIST")
			release.ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
			return
		}

		// Find KV in Room Bucket
		//b := boltdb.View([]byte(roomId), "RoomBucket")
		//if b == nil {
		//	util.Error("ERR_ROOM_NOT_EXIST")
		//	release.ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		//	return
		//}

		// Decode Room Info
		roomInfo := model.RoomInfo{}
		//de := json.Unmarshal(b, &roomInfo)
		//if de != nil {
		//	//util.Logger(util.ERROR_LEVEL, "WxJoin", "Decoding Room Info Err:"+de.Error())
		//	util.Error("WxJoin Decoding Room Info Err:" + de.Error())
		//}
		//读取房间信息
		err = gocache.GetRoomInfo(roomId, &roomInfo)
		if err != nil {
			util.Error("ERROR[%v]", err.Error())
			release.ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
			return
		}
		//util.Debug("房间信息[%+v]", roomInfo)

		req.UnionIds = roomInfo.UnionIdSlice
		costs2, code := release.CheckUserScriptCost(req)
		if code != model.ERR_OK && code != global.ERR_DB_NOTFOUND_DATA {
			util.Error("CheckUserScriptCost Get User ScriptId Err:%v", code)
			release.ResponseErr(model.ERR_CHECK_SCRIPT_COST, c)
			return
		}
		util.Debug("costs2:%+v", costs2)

		if code == global.ERR_DB_NOTFOUND_DATA {
			//iNum = 0
			for _, unionId := range roomInfo.UnionIdSlice {
				//查询用户是否为会员
				res, err := release.GetUserMemberBase(unionId)
				if err != nil {
					util.Error("ERROR:%v", err.Error())
					release.ResponseErr(model.ERR_DEFAULT, c)
					return
				}
				if res.Member > 0 {
					util.Info("用户:%v 会员等级:%v 需要付费...", unionId, res.Member)
					iNum++
				}
			}

		} else {
			for _, tmp := range costs2 {
				util.Debug("校验用户:%v", tmp.UnionId)
				//查询用户是否为会员
				res, err := release.GetUserMemberBase(tmp.UnionId)
				if err != nil {
					util.Error("ERROR:%v", err.Error())
					release.ResponseErr(model.ERR_DEFAULT, c)
					return
				}
				if tmp.CostFlag == 2 || res.Member > 0 {
					util.Info("用户:%v 会员等级:%v 需要付费...", tmp.UnionId, res.Member)
					iNum++
				}
			}
		}

		util.Info("已付款的用户有[%v]个", iNum)
		//计算已付款的，不在进行付款
		//price = scriptBase.Price * float64(scriptBase.Num-iNum)
		//if price < 0.01 {
		price = scriptBase.Price * (scriptBase.Num - iNum)
		if price < 1 {
			util.Info("价格低于1分钱，无需付款")
			cost.CostFlag = 2
		}
		//房间用户全部登记
		unionIds = append(unionIds, roomInfo.UnionIdSlice...)
	}
	util.Info("付费标识:%v 付费价格[%v],用户切片:%+v", cost.CostFlag, price, unionIds)
	if cost.CostFlag != 2 { //已付费
		//生成支付参数
		res, params, err := WeChatPayParam(jparm.OpenId, orderNo, price)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			release.ResponseErr(model.ERR_OPENID_FAILED, c)
			return
		}
		util.Logger(util.INFO_LEVEL, "wx.pay_param:", res)

		//var resp model.WxUserPhoneResp
		//err = json.Unmarshal([]byte(str), &resp)
		//if err != nil {
		//	ResponseErrMsg(c, "解析JSON串出错")
		//	return
		//}

		taskRequest := make(map[string]interface{})
		taskRequest["OutTradeNo"] = params.OutTradeNo
		taskRequest["PayAmt"] = params.PayAmt
		taskRequest["PrepayId"] = params.PrepayId
		taskRequest["CodeUrl"] = params.CodeUrl
		taskRequest["TotalFee"] = params.TotalFee
		taskRequest["ReturnCode"] = params.ReturnCode
		taskRequest["ReturnMsg"] = params.ReturnMsg
		taskRequest["ResultCode"] = params.ResultCode
		taskRequest["ErrCode"] = params.ErrCode
		taskRequest["ErrCodeDes"] = params.ErrCodeDes
		taskRequest["Content"] = params.Content
		taskRequest["PayContent"] = params.PayContent
		taskRequest["UnionIds"] = unionIds
		taskRequest["ScriptId"] = jparm.SId

		util.Debug("record pay params taskRequest:", taskRequest)
		dbResult, err := global.Task.TaskJson(global.NewDBRequest("db.PayParams", taskRequest))
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			release.ResponseErr(model.ERR_TASK_JSON, c)
			return
		}
		dbcode, dbparams := global.UnwrapObjectPackage(dbResult)

		util.Debug("UnwrapObjectPackage:", dbcode, dbparams)
		rsp.Res = res
	}

	//准备响应信息
	rsp.Err = model.ERR_OK
	rsp.Msg = model.ErrMap[model.ERR_OK]
	rsp.Res.UserScript = cost
	rsp.Res.ScriptName = scriptBase.Name
	rsp.Res.UnpaidNum = scriptBase.Num - iNum
	util.Info("响应信息:%+v", rsp)

	release.ResponseOk(c, rsp)
}
