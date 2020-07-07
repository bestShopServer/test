package release

import (
	"DetectiveMasterServer/config"
	"DetectiveMasterServer/global"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/util"
	"gitee.com/sunki/gutils/encryption"
	"github.com/gin-gonic/gin"
)

func SyncUserBase(param model.GetWxUserPhoneReq) (code int) {
	util.Info("Sync User Base ...")

	taskRequest := make(map[string]interface{})
	taskRequest["UnionId"] = param.UnionId
	taskRequest["Phone"] = param.Phone
	taskRequest["NickName"] = param.NickName
	taskRequest["OpenId"] = param.OpenId
	taskRequest["Gender"] = param.Gender
	taskRequest["Country"] = param.Country
	taskRequest["Province"] = param.Province
	taskRequest["City"] = param.City

	util.Info("User:%+v", taskRequest)
	dbResult, err := global.Task.TaskJson(global.NewDBRequest("db.wx.UserSync", taskRequest))
	if err != nil {
		util.Error("db.ScriptQuestionGet ERROR[%v]", err.Error())
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

// Func: wx.phone handler
func WxUserPhone(c *gin.Context) {
	util.Debug("wx get weixin user phone ...")

	//获取请求数据
	var jparm model.GetWxUserPhoneReq
	err := c.ShouldBindJSON(&jparm)
	if err != nil {
		util.Error("ERR_WRONG_FORMAT:%+v", jparm)
		ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	util.Info("请求参数:%+v", jparm)
	optionalFields := []string{"Phone", "City", "Province", "Country", "Gender",
		"EncryptedData", "Iv", "SessionKey"}

	if !CheckParams(jparm, "WxUserPhone", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", jparm)
		ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}

	// 2020-01-17 去掉手机号
	//str, err := util.DncryptWx(jparm.EncryptedData, jparm.SessionKey, jparm.Iv)
	//if err != nil {
	//	ResponseErr(model.ERR_WRONG_FORMAT, c)
	//	return
	//}
	//util.Info("解析后字符串[%v]", str)

	var resp model.WxUserPhoneResp
	//err = json.Unmarshal([]byte(str), &resp)
	//if err != nil {
	//	ResponseErrMsg(c, "解析JSON串出错")
	//	return
	//}

	//同步用户数据
	jparm.Phone = resp.PurePhoneNumber
	code := SyncUserBase(jparm)
	if code != model.ERR_OK {
		util.Logger(util.ERROR_LEVEL, "SyncUserBase", "Record Score Err")
		ResponseErr(model.ERR_USER_BASE_INFO_SYNC, c)
		return
	}
	util.Info("处理结束...")

	//生成Token
	token, err := encryption.GenerateToken(jparm.UnionId, resp.PurePhoneNumber,
		jparm.OpenId, config.GetConfig().TokenTimeout)
	if err != nil {
		ResponseErr(model.ERR_DEFAULT, c)
		return
	}

	//准备响应信息
	rsp := model.GetWxUserPhoneResp{}
	rsp.Err = model.ERR_OK
	rsp.Msg = model.ErrMap[model.ERR_OK]
	//rsp.Notice = notice
	//rsp.Num = num
	rsp.Res = resp
	rsp.Res.Token = token

	//util.Info("响应信息:%+v", rsp)
	ResponseOk(c, rsp)
}
