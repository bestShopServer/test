package user

import (
	"DetectiveMasterServer/gocache"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/release"
	"DetectiveMasterServer/util"
	"github.com/gin-gonic/gin"
)

// 查询用户所在房间号
func WxUserRoomBase(c *gin.Context) {
	util.Info("WxUserRoomBase ...")
	//var json = jsoniter.ConfigCompatibleWithStandardLibrary
	// Get Params
	req := model.WxUserRoomReq{}
	res := model.MemberBaseInfo{}
	err := c.BindJSON(&req)

	// Optional Fields List
	var optionalFields []string

	// Check Params
	if !release.CheckParams(req, "WxExit", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", req)
		release.ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	util.Info("请求参数[%+v]", req)

	// Get OpenId
	uid := req.UnionId
	//获取用户缓存房间号
	roomId, err := gocache.GetUserRoom(uid)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		//release.ResponseErr(model.ERR_NOT_IN_ROOM, c)
		//return
	}

	//判断上一房间号状态
	roomInfo := model.RoomInfo{}
	err = gocache.GetRoomInfo(roomId, &roomInfo)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		//release.ResponseErr(model.ERR_ROOM_LINK, c)
		//return
	}
	if roomInfo.Status == 2 {
		util.Info("上局房间%v已结束", roomId)
		roomId = "" // 房间号赋空
		//release.ResponseErr(model.ERR_ROOM_ALREADY_OVER, c)
		//return
	}

	//查看用户是否在白名单
	isWhite := false
	whiteUsers, err := gocache.GetWhitelist()
	if err != nil {
		util.Error("获取白名单失败, ERROR:%v", err.Error())
	}
	util.Info("白名单数据:%+v", whiteUsers)
	for _, tmp := range whiteUsers {
		if tmp == req.UnionId {
			isWhite = true
		}
	}
	util.Info("是否在白名单:%v", isWhite)
	if isWhite {
		res.Member = 1
		res.InvTime = "2099-12-30 00:00:00"
		util.Debug("会员级别:%v", res.Member)
	} else {
		res, err = release.GetUserMemberBase(req.UnionId)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			//release.ResponseErr(model.ERR_DEFAULT, c)
			//return
		}
		util.Debug("会员信息:%+v", res)
	}

	// Return
	resp := model.UserRoomResp{}
	resp.Params.RoomId = roomId
	resp.Params.MemberBaseInfo = res
	release.ResponseOk(c, resp)
}
