package release

import (
	"DetectiveMasterServer/config"
	"DetectiveMasterServer/global"
	"DetectiveMasterServer/gocache"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/util"
	"DetectiveMasterServer/websocket"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"
	"net/http"
)

// Func: wx.create handler
func WxCreate(c *gin.Context) {
	fmt.Println("WxCreate ...")

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	//fmt.Println("json:", json)

	// Get Code Param
	var req model.WxCreateReq
	err := c.Bind(&req)

	// Optional Fields List
	var optionalFields []string

	// Check Param
	if !CheckParams(req, "WxCreate", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", req)
		ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	util.Info("请求参数[%+v]", req)

	//uid := req.OpenId
	uid := req.UnionId
	sname := req.ScriptName

	//校验能否创建新房间
	ok, err := CheckJoinNewRoom(uid)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseNotice(c, model.ERR_ROOM_NOT_OVER, err.Error())
		return
	}
	if !ok {
		util.Error("状态:%v", ok)
		ResponseErr(model.ERR_ROOM_NOT_OVER, c)
		return
	}

	// Get Script Id By Name
	script, no := GetScriptIdByName(sname)
	if no > 0 {
		util.Error("剧本不存在！")
		ResponseErr(no, c)
		return
	}
	//util.Info("新建剧本基本信息:%+v", script)

	// Create Room
	rid := global.CreateRoomId(8)
	util.Info("WxCreate rid:", rid)

	var UserInfoSlice []model.UserInfo
	var unionIdSlice []string
	var user model.UserInfo
	var cost model.ScriptUserCostReq

	user.UnionId = uid
	user.Name = "房主"
	cost.ScriptId = script.Id
	cost.UnionIds = append(cost.UnionIds, req.UnionId)

	//校验用户是否付款
	util.Info("剧本价格:%v", script.Price)
	if script.Price > 0 {
		us, code := CheckUserScriptCost(cost)
		util.Debug("查询用户剧本付费情况,code:%v", code)
		if code != model.ERR_OK && code != global.ERR_DB_NOTFOUND_DATA {
			util.Error("Get User ScriptId Err [%v]", code)
			ResponseErr(model.ERR_CHECK_SCRIPT_COST, c)
			return
		}
		if code == global.ERR_DB_NOTFOUND_DATA { //未付款
			user.IsPay = 1
		} else {
			user.IsPay = us[0].CostFlag
		}
	}
	unionIdSlice = append(unionIdSlice, uid)
	UserInfoSlice = append(UserInfoSlice, user)

	util.Info("WxCreate unionIdSlice:", unionIdSlice)
	playerSlice, code := GetPlayerInfoByScriptId(script.Id)
	if code != model.ERR_OK {
		//util.Logger(util.ERROR_LEVEL, "WxCreate", "Get Role Info By Script Name Err")
		util.Error("WxCreate Get Role Info By Script Name Err")
		ResponseErr(model.ERR_GET_ROLE_INFO, c)
		return
	}
	//util.Info("WxCreate playerSlice:%+v \n code:%v", playerSlice, code)

	rif := model.RoomInfo{
		ScriptId:    script.Id,
		Num:         script.Num,
		Owner:       uid,
		Price:       float64(script.Price / 100),
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

	//roomInfoEncoding, err := json.Marshal(rif)
	//if err != nil {
	//	util.Logger(util.ERROR_LEVEL, "RoomNew", "Room Info Encoding Err:"+err.Error())
	//}
	//util.Debug("roomInfoEncoding:%v", string(roomInfoEncoding))
	//
	//boltdb.CreateOrUpdate([]byte(rid), roomInfoEncoding, "RoomBucket")
	err = gocache.SetRoomInfo(rid, rif) //保存房间信息
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	// Update UserCache & RoomCache
	//global.SetUserCache(uid, rid)
	//global.AddUserToRoomCache(uid, rid)
	//global.AddRoomDeleteTask(rid)

	err = gocache.SetUserRoom(uid, rid) //保存用户房间号
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}
	err = gocache.SetRoomUser(rid, uid) //保存房间用户列表
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	// Send RoomEnter Message
	//go websocket.SendRoomEnterMessage(unionIdSlice, uid)
	go websocket.SendRoomEnterMessage(unionIdSlice, user)

	// Return roomId & playerSlice
	resp := model.WxCreateResp{}
	resp.Params.RoomId = rid
	resp.Params.RoomName = script.Name
	//resp.Params.Path = fmt.Sprintf("/pages/loadingnew/loadingnew?is_create=1&script_id=%d"+
	//	"&unchoosed_num=%d&open_id=%v&room_id=%v&union_id=%v&num=%d&price=%.2f",
	//	script.Id, len(rif.PlayerSlice), uid, rid, req.UnionId, script.Num, float64(script.Price/100))
	resp.Params.Path = fmt.Sprintf("/pages/loadingnew/loadingnew?is_create=1&script_id=%d"+
		"&open_id=%v&room_id=%v&num=%d&price=%.2f",
		script.Id, uid, rid, script.Num, float64(script.Price/100))

	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	if err = jsonEncoder.Encode(resp); err != nil {
		ResponseErr(model.ERR_DEFAULT, c)
		return
	}

	util.Info("WxCreate:%v", bf.String())

	c.Data(http.StatusOK, "application/json", bf.Bytes())
}

// Func: Get Script Id By Name
func GetScriptIdByName(scriptName string) (res model.Script, errno int) {
	req := model.ScriptGetReq{}
	req.Search = scriptName

	conf := config.GetConfig()
	fmt.Println("GetScriptIdByName conf:", conf)
	FixParams(&req, conf)

	fmt.Println("GetScriptIdByName req:", req)

	scriptList, code := GetScripts(req)
	fmt.Println("GetScriptIdByName scriptList:", scriptList)
	fmt.Println("GetScriptIdByName code:", code)

	if code != model.ERR_OK {
		util.Logger(util.ERROR_LEVEL, "ScriptGet", "Get Scripts Err")
		return res, model.ERR_GET_SCRIPTS
	}

	if len(scriptList) > 0 {
		res = scriptList[0]
		return res, 0
	} else {
		return res, model.ERR_GET_SCRIPTS
	}
}
