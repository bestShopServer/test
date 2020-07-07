package game

import (
	"DetectiveMasterServer/config"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/release"
	"DetectiveMasterServer/util"
	"github.com/gin-gonic/gin"
)

//判断版本号大于服务器设置的版本号则展示微信审核页入口
func Audit(c *gin.Context) {
	util.Info("Audit ...")
	var jparm model.AuditReq
	var resp model.AuditResp

	err := c.Bind(&jparm)
	util.Info("请求参数[%+v]", jparm)

	// Optional Fields List
	optionalFields := []string{}
	// Check Param
	if !release.CheckParams(jparm, "Audit", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", jparm)
		release.ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	version := config.GetConfig().Version
	util.Info("服务器版本:%v 小程序版本:%v", version, jparm.Version)
	if jparm.Version > version {
		resp.Params.IsShow = 1
	} else {
		resp.Params.IsShow = 2
	}
	util.Info("处理结束...")

	release.ResponseOk(c, &resp)
}
