package room

import (
	"DetectiveMasterServer/global"
	"DetectiveMasterServer/gocache"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/release"
	"DetectiveMasterServer/release/point"
	"DetectiveMasterServer/util"
	"DetectiveMasterServer/websocket"
	"github.com/gin-gonic/gin"
	"strconv"
)

//用户强制加入新房间
func WxRoomJoinV2(c *gin.Context) {
	util.Info("WxJoin V2...")
	var user model.UserInfo
	var req model.WxJoinReq
	roomRoleNum := 0 //房间可被选人数

	err := c.Bind(&req)

	// Optional Fields List
	//var optionalFields []string
	optionalFields := []string{"UserName"}

	// Check Param
	if !release.CheckParams(req, "WxJoin", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", req)
		release.ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	util.Info("请求参数[%+v]", req)

	//uid := req.OpenId
	uid := req.UnionId
	rid := req.RoomId
	user.UnionId = req.UnionId
	user.Name = req.UserName
	if len(user.Name) == 0 {
		user.Name = "正在进入"
	}

	//获取用户所在房间号
	roomCode, err := gocache.GetUserRoom(uid)
	if err != nil {
		util.Error("ERROR[%v]", err.Error())
		release.ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}
	util.Info("用户:%v在房间:%v", uid, roomCode)

	// Decode Room Info
	rif := model.RoomInfo{}
	//读取房间信息
	err = gocache.GetRoomInfo(rid, &rif)
	if err != nil {
		util.Error("ERROR[%v]", err.Error())
		release.ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	//房间游戏已结束，则无法再加入房间
	if rif.Status == 2 { //标识房间已结束
		util.Error("房间:%v 状态:%v", rid, rif.Status)
		release.ResponseErr(model.ERR_ROOM_ALREADY_OVER, c)
		return
	}

	//房间存在切不等于上送的房间号
	if roomCode != "" && roomCode != rid {
		//无线校验能否创建新房间，直接创建新房间
		err = ExitLastRoom(uid)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			release.ResponseErr(model.ERR_DEFAULT, c)
			return
		}
	}

	//判断用户是否在房间
	bl, err := gocache.CheckUserInRoom(rid, uid)
	if err != nil {
		util.Error("ERROR[%v]", err.Error())
		release.ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}
	if bl {
		util.Info("用户:%v已在房间:%v", uid, rid)
	} else {
		//重新计算房间可被选人
		for _, tmp := range rif.PlayerSlice {
			if tmp.Role.Choice {
				roomRoleNum += 1
			}
		}
		util.Info("剧本可被选角色数:%d, 当前房间人数:%v", roomRoleNum, len(rif.UnionIdSlice))

		// If openIdSlice length eq playerSlice
		//if len(rif.UnionIdSlice) == len(rif.PlayerSlice) {
		if len(rif.UnionIdSlice) >= roomRoleNum {
			util.Error("ERR_ROOM_IS_FULL")
			release.ResponseErr(model.ERR_ROOM_IS_FULL, c)
			return
		}

		//用户超过房间人数也需要报错
		if len(rif.UnionIdSlice) >= rif.Num {
			util.Error("ERR_ROOM_IS_FULL")
			release.ResponseErr(model.ERR_ROOM_IS_FULL, c)
			return
		}

		if len(rif.UserInfoSlice) >= rif.Num {
			util.Error("房间人数大于剧本人数，需要踢出用户")
			release.ResponseErr(model.ERR_ROOM_IS_FULL, c)
			return
		}

		// If openIdSlice Contains openId
		exist := false
		for _, v := range rif.UnionIdSlice {
			if v == uid {
				exist = true
				break
			}
		}

		useExit := false
		for _, u := range rif.UserInfoSlice {
			if u.UnionId == uid {
				useExit = true
				break
			}
		}

		//没数据则增加用户信息
		if !useExit {
			util.Info("剧本价格:%f", rif.Price)
			if rif.Price > 0.0 {
				var cost model.ScriptUserCostReq
				cost.ScriptId = rif.ScriptId
				cost.UnionIds = append(cost.UnionIds, req.UnionId)

				//查询付款
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
				if tmp == req.UnionId {
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

			rif.UserInfoSlice = append(rif.UserInfoSlice, user) //添加用户
		}

		//用户不存在
		if !exist {
			// Add openId to openIdSlice
			rif.UnionIdSlice = append(rif.UnionIdSlice, uid)

			err = gocache.SetRoomInfo(rid, rif)
			if err != nil {
				util.Error("ERROR:%v", err.Error())
				release.ResponseErr(model.ERR_ENTERED_ROOM, c)
				return
			}

			//// Update UserCache & RoomCache
			err = gocache.SetUserRoom(uid, rid) //保存用户房间号
			if err != nil {
				util.Error("ERROR:%v", err.Error())
				release.ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
				return
			}
			err = gocache.AddRoomUser(rid, uid) //保存房间用户列表
			if err != nil {
				util.Error("ERROR:%v", err.Error())
				release.ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
				return
			}

			// Send Room Enter Msg
			//go websocket.SendRoomEnterMessage(rif.UnionIdSlice, uid)
			go websocket.SendRoomEnterMessage(rif.UnionIdSlice, user)

			//记录用户进入房间
			param := model.RoomRecordBase{}
			param.ScriptId = rif.ScriptId
			param.RoomId = rid
			param.Owner = rif.Owner
			param.UnionId = uid
			go point.RecordRoomUserJoinData(param)
		}
	}

	// Return Room Info
	count := release.CountUnchoosedPeople(rif.PlayerSlice)
	resp := model.WxJoinResp{}
	resp.Params.Path = "/pages/pick/pick?is_create=0&script_id=" +
		strconv.Itoa(rif.ScriptId) + "&unchoosed_num=" + strconv.Itoa(count) +
		"&open_id=" + uid + "&union_id=" + uid + "&room_id=" + rid

	release.ResponseOk(c, resp)
}
