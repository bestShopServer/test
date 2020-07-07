package order

import (
	"DetectiveMasterServer/global"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/release"
	"DetectiveMasterServer/util"
	"github.com/gin-gonic/gin"
)

func queryOrderState(param model.OrderStateReq) (code int) {
	util.Info("Record Script Order Sync ...")

	taskRequest := make(map[string]interface{})
	taskRequest["ScriptName"] = param.ScriptName
	taskRequest["UnionId"] = param.UnionId

	util.Info("UnionId[%+v]", taskRequest)
	dbResult, err := global.Task.TaskJson(global.NewDBRequest("db.wx.UserScriptOrder", taskRequest))
	if err != nil {
		util.Error("db.wx.PublicScriptOrder ERROR[%v]", err.Error())
		return model.ERR_TASK_JSON
	}
	util.Info("获取数据成功[%+v]", dbResult)
	dbcode, _ := global.UnwrapArrayPackage(dbResult)

	switch dbcode {
	case global.ERR_DB_OK:
		return model.ERR_OK
	default:
		return model.ERR_DEFAULT
	}
}

//查询用户订单状态
func PostOrderState(c *gin.Context) {
	util.Info("PublicOrderSync ...")

	var req model.OrderStateReq

	err := c.Bind(&req)
	util.Info("请求参数[%+v]", req)

	// Optional Fields List
	optionalFields := []string{}
	// Check Param
	if !release.CheckParams(req, "PublicOrderSync", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", req)
		release.ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}

	code := queryOrderState(req)
	if code != model.ERR_OK {
		util.Error("code:%v", code)
		release.ResponseErr(model.ERR_SCRIPT_ORDER_SYNC, c)
		return
	}
	util.Info("处理结束...")

	Resp := model.ErrResp{}
	release.ResponseOk(c, &Resp)
}
