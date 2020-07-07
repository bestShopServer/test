package release

import (
	"DetectiveMasterServer/config"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/util"
	"github.com/gin-gonic/gin"
)

// Func: wx.login handler
func WxGenSig(c *gin.Context) {
	util.Info("WxGenSig...")

	// Get Code Param
	var req model.GenSigReq
	err := c.ShouldBind(&req)

	// Optional Fields List
	var optionalFields []string

	// Check Param
	if !CheckParams(req, "WxGenSig", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", req)
		ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	util.Info("请求参数[%+v]", req)

	// Get OpenId From Tencent
	conf := config.GetConfig()
	//sig, err := util.GenSig(conf.ImSdkAppId, conf.ImKey, conf.ImIdent, 86400*180)
	sig, err := util.GenSig(conf.ImSdkAppId, conf.ImKey, req.OpenId, 86400*180)
	if err != nil {
		util.Logger(util.ERROR_LEVEL, "wx.genSin:", sig)
	}
	util.Logger(util.INFO_LEVEL, "wx.genSin:", sig)

	rtn := model.ImGenSigResp{}
	rtn.UserSig = sig
	ResponseOk(c, rtn)
}
