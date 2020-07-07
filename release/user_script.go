package release

import (
	"DetectiveMasterServer/global"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/util"
	"gitee.com/sunki/gutils/trans"
)

func CheckUserScriptCost(req model.ScriptUserCostReq) ([]model.UserScript, int) {
	util.Info("Check User Script Cost...")
	var userScripts []model.UserScript

	taskRequest := make(map[string]interface{})
	taskRequest["ScriptId"] = req.ScriptId
	taskRequest["UnionIds"] = req.UnionIds
	taskRequest["OpenId"] = req.OpenId
	util.Info("UnionIds[%v]", req.UnionIds)

	util.Info("db.ScriptUserCost ...")
	dbResult, err := global.Task.TaskJson(global.NewDBRequest("db.ScriptUserCost", taskRequest))
	if err != nil {
		util.Error("db.ScriptUserCost ERROR[%v]", err.Error())
		return userScripts, model.ERR_TASK_JSON
	}
	util.Info("获取数据成功...")

	dbcode, dbparams := global.UnwrapArrayPackage(dbResult)

	switch dbcode {
	case global.ERR_DB_OK:
		if len(dbparams) == 0 {
			util.Info("数据数为零!")
			return userScripts, global.ERR_DB_NOTFOUND_DATA
		}
		for _, p := range dbparams {
			var tmp model.UserScript
			pm := p.(map[string]interface{})
			if pm["ScriptId"] != nil {
				tmp.ScriptId = int(pm["ScriptId"].(float64))
			}
			if pm["UserId"] != nil {
				tmp.UserId = int(pm["UserId"].(float64))
			}
			if pm["Status"] != nil {
				util.Info("付款状态:%v", pm["Status"])
				tmp.Status = int(pm["Status"].(float64))
				if tmp.Status == 2 {
					tmp.CostFlag = 2 //付款标识 1 未付款 2 已付款
					util.Info("标识已付款:%v", tmp.CostFlag)
				}
			}

			if pm["EffTime"] != nil {
				tmp.EffTime = trans.StringToDate(pm["EffTime"].(string))
			}
			if pm["InvTime"] != nil {
				tmp.InvTime = trans.StringToDate(pm["InvTime"].(string))
			}
			if pm["UnionId"] != nil {
				tmp.UnionId = pm["UnionId"].(string)
			}
			userScripts = append(userScripts, tmp)
		}
		return userScripts, model.ERR_OK
	default:
		return userScripts, model.ERR_DEFAULT
	}
}
