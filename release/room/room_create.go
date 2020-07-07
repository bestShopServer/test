package room

import (
	"DetectiveMasterServer/config"
	"DetectiveMasterServer/global"
	"DetectiveMasterServer/gocache"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/release"
	"DetectiveMasterServer/release/point"
	"DetectiveMasterServer/util"
	"DetectiveMasterServer/websocket"
	"fmt"
	"github.com/gin-gonic/gin"
)

// Func: create room handler
func RoomCreate(c *gin.Context) {
	fmt.Println("RoomCreate ...")

	// Get Code Param
	req := model.WxRoomCreateReq{}
	UserInfoSlice := []model.UserInfo{}
	unionIdSlice := []string{}
	user := model.UserInfo{}
	cost := model.ScriptUserCostReq{}
	err := c.BindJSON(&req)

	// Optional Fields List
	//var optionalFields []string
	optionalFields := []string{"Flag"}

	// Check Param
	if !release.CheckParams(req, "RoomCreate", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", req)
		release.ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	util.Info("请求参数[%+v]", req)

	//uid := req.OpenId
	uid := req.UnionId
	sid := req.ScriptId

	//无线校验能否创建新房间，直接创建新房间
	err = ExitLastRoom(uid)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		release.ResponseErr(model.ERR_DEFAULT, c)
		return
	}

	// Get Script Id By Name
	script, no := GetScriptById(sid)
	if no > 0 {
		util.Error("剧本不存在！")
		release.ResponseErr(no, c)
		return
	}
	//util.Info("新建剧本基本信息:%+v", script)

	// Create Room
	rid := global.CreateRoomId(8)
	util.Info("WxCreate rid:", rid)

	user.UnionId = uid
	user.Name = "房主"
	cost.ScriptId = script.Id
	cost.UnionIds = append(cost.UnionIds, req.UnionId)

	//校验用户是否付款
	util.Info("剧本价格:%v", script.Price)
	if script.Price > 0 {
		us, code := release.CheckUserScriptCost(cost)
		util.Debug("查询用户剧本付费情况,code:%v", code)
		if code != model.ERR_OK && code != global.ERR_DB_NOTFOUND_DATA {
			util.Error("Get User ScriptId Err [%v]", code)
			release.ResponseErr(model.ERR_CHECK_SCRIPT_COST, c)
			return
		}
		if code == global.ERR_DB_NOTFOUND_DATA { //未付款
			user.IsPay = 1
		} else {
			user.IsPay = us[0].CostFlag
		}
	}

	//查看用户是否在白名单
	isWhite := false
	whiteUsers, err := gocache.GetWhitelist()
	if err != nil {
		util.Error("获取白名单失败, ERROR:%v", err.Error())
	}
	util.Info("白名单数据:%+v", whiteUsers)
	for _, tmp := range whiteUsers {
		if tmp == uid {
			isWhite = true
		}
	}
	util.Info("是否在白名单:%v", isWhite)
	if isWhite {
		user.Member = 1
		util.Debug("会员级别:%v", user.Member)
	} else {
		//查询用户是否为会员
		res, err := release.GetUserMemberBase(req.UnionId)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			release.ResponseErr(model.ERR_DEFAULT, c)
			return
		}
		user.Member = res.Member
		util.Debug("会员级别:%v", user.Member)
	}

	unionIdSlice = append(unionIdSlice, uid)
	UserInfoSlice = append(UserInfoSlice, user)

	util.Info("WxCreate unionIdSlice:", unionIdSlice)
	playerSlice, code := release.GetPlayerInfoByScriptId(script.Id)
	if code != model.ERR_OK {
		//util.Logger(util.ERROR_LEVEL, "WxCreate", "Get Role Info By Script Name Err")
		util.Error("WxCreate Get Role Info By Script Name Err")
		release.ResponseErr(model.ERR_GET_ROLE_INFO, c)
		return
	}
	//util.Info("WxCreate playerSlice:%+v \n code:%v", playerSlice, code)

	rif := model.RoomInfo{
		ScriptId:    script.Id,
		Num:         script.Num,
		Owner:       uid,
		Price:       float64(script.Price) / 100.00,
		Flag:        req.Flag,
		VoteFlag:    script.VoteFlag,
		TopicFlag:   script.TopicFlag,
		ExploreFlag: script.ExploreFlag,
		//OpenIdSlice: openIdSlice,
		UnionIdSlice:  unionIdSlice,
		PlayerSlice:   playerSlice,
		UserInfoSlice: UserInfoSlice,
		Round:         script.Round,
	}
	//util.Info("房间信息:%+v", rif)

	err = gocache.SetRoomInfo(rid, rif) //保存房间信息
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		release.ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	err = gocache.SetUserRoom(uid, rid) //保存用户房间号
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		release.ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}
	err = gocache.SetRoomUser(rid, uid) //保存房间用户列表
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		release.ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	// Send RoomEnter Message
	go websocket.SendRoomEnterMessage(unionIdSlice, user)

	//记录房间数据
	param := model.RoomRecordBase{}
	param.ScriptId = sid
	param.RoomId = rid
	param.Owner = uid
	go point.RecordRoomMainData(param)

	// Return roomId & playerSlice
	resp := model.RoomCreateResp{}
	resp.Params.RoomId = rid
	resp.Params.RoomName = script.Name
	resp.Params.Num = script.Num
	resp.Params.Price = float64(script.Price / 100)

	release.ResponseOk(c, resp)
}

// Func: Get Script Id By Name
func GetScriptById(scriptId int) (res model.Script, errno int) {
	req := model.ScriptGetReq{}
	req.ScriptId = scriptId
	util.Info("剧本ID:%v", scriptId)

	conf := config.GetConfig()
	util.Debug("GetScriptIdByName conf:%+v", conf)
	release.FixParams(&req, conf)
	scriptList, code := release.GetScripts(req)
	//util.Info("scriptList:%+v code:%v", scriptList, code)

	if code != model.ERR_OK {
		util.Error("Get Script Error")
		return res, model.ERR_GET_SCRIPTS
	}

	if len(scriptList) > 0 {
		res = scriptList[0]
		return res, 0
	} else {
		return res, model.ERR_GET_SCRIPTS
	}
}
