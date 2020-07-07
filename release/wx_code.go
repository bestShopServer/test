package release

import (
	conf "DetectiveMasterServer/config"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/util"
	"github.com/gin-gonic/gin"
)

// Func: wx.code handler
func WxCode2Session(c *gin.Context) {
	// Get Code Param
	var req model.Code2SessionReq
	err := c.ShouldBind(&req)

	// Optional Fields List
	var optionalFields []string

	// Check Param
	if !CheckParams(req, "WxCode2Session", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", req)
		ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	util.Info("请求参数[%+v]", req)

	// Get OpenId From Tencent
	config := conf.GetConfig()
	//openId, err := util.GetOpenId(config.AppId, config.AppSecret, req.Code, req.EncryptData, req.IV)
	res, err := util.WeChatLogin(config.AppId, config.AppSecret, req.Code)
	if err != nil {
		ResponseErr(model.ERR_OPENID_FAILED, c)
		return
	}
	util.Logger(util.INFO_LEVEL, "wx.code:", res)

	rtn := model.Code2SessionResp{}
	rtn.OpenId = res.OpenId
	rtn.SessionKey = res.SessionKey
	rtn.UnionId = res.UnionId
	ResponseOk(c, rtn)
}
