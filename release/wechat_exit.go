package release

import (
	"DetectiveMasterServer/gocache"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Func: wx.exit handler
func WechatExit(c *gin.Context) {
	util.Info("WxExit ...")

	//var json = jsoniter.ConfigCompatibleWithStandardLibrary

	// Get Params
	var req model.WxExitReq
	err := c.Bind(&req)

	// Optional Fields List
	var optionalFields []string

	// Check Params
	if !CheckParams(req, "WxExit", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", req)
		ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	util.Info("请求参数[%+v]", req)

	// Get UnionId
	uid := req.UnionId

	// Update User Cache
	ok, err := gocache.Delete(uid)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}
	util.Info("OK:%v", ok)

	// Return
	resp := model.WxExitResp{}
	resp.Params = "success"
	c.JSON(http.StatusOK, resp)
}
