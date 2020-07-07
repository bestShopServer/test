package release

import (
	"DetectiveMasterServer/global"
	"DetectiveMasterServer/gocache"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/release/point"
	"DetectiveMasterServer/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Func: game.report handler
/*
只有投凶
凶手 :
 没被投 100分
 被投出去 20分
 没被投出去 80分
平民
 没被投 投对了 100分
 没被投 投错了 80分
 投对了 被投出去了 40分
 投错了 被投出去了 20分
 被投但没被投出去 投对了 80分
 被投但没被投出去  投错了 60分
其他情况：
 1、若该剧本只有答题没有投凶，则每道题的实际得分相加➗每道题的最高得分的总和✖️100位=最终得分
 2、若该剧本即有答题又有投凶，则（投凶对应得分➕答题实际得分）➗（100➕每道题最高得分的总和）✖️100=最终得分
*/
func GameReport(c *gin.Context) {
	util.Info("GameReport ...")
	//var json = jsoniter.ConfigCompatibleWithStandardLibrary

	var req model.GameReportReq
	var reportUrl string
	var code int
	var lastScore int
	var roleId int

	err := c.Bind(&req)

	// Optional Fields List
	var optionalFields []string

	// Check Param
	if !CheckParams(req, "GameReport", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", req)
		ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	util.Info("请求参数:%+v", req)

	//userId := req.OpenId
	userId := req.UnionId
	roomId := req.RoomId

	// Find Room In RoomBucket
	//b := boltdb.View([]byte(roomId), "RoomBucket")
	//if b == nil {
	//	util.Error("ERR_ROOM_NOT_EXIST")
	//	ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
	//	return
	//}

	// Decode Room Info
	roomInfo := model.RoomInfo{}
	//err = json.Unmarshal(b, &roomInfo)
	//if err != nil {
	//	//util.Logger(util.ERROR_LEVEL, "GameReport", "Decoding Room Info Err:"+err.Error())
	//	util.Error("GameReport Decoding Room Info ERROR[%v]", err.Error())
	//}
	//读取房间信息
	err = gocache.GetRoomInfo(roomId, &roomInfo)
	if err != nil {
		util.Error("ERROR[%v]", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	// If You are not in this room
	exist := false
	for i, v := range roomInfo.UnionIdSlice {
		if v == userId {
			exist = true
			//用户查看结案后，从房间移出
			roomInfo.UnionIdSlice = append(roomInfo.UnionIdSlice[:i], roomInfo.UnionIdSlice[i+1:]...)
			break
		}
	}
	if !exist {
		util.Error("ERR_BELONG")
		ResponseErr(model.ERR_BELONG, c)
		return
	}

	util.Info("投票标识:%v", roomInfo.VoteFlag)
	if roomInfo.VoteFlag {
		//// Find Room In VoteBucket
		//vb := boltdb.View([]byte(roomId), "VoteBucket")
		//if vb != nil {
		//	// Decode Vote Map
		var votes map[string]bool
		//	err = json.Unmarshal(vb, &votes)
		//	if err != nil {
		//		//util.Logger(util.ERROR_LEVEL, "GameVote", "Decoding Votes Err:"+err.Error())
		//		util.Error("GameVote Decoding Votes ERROR[%v]", err.Error())
		//	}
		//
		//	// Check If Or Not All People Has Voted
		//	all := true
		//	for _, v := range votes {
		//		if !v {
		//			all = false
		//			break
		//		}
		//	}
		//	if !all {
		//		util.Error("ERR_NOT_ALL_VOTED")
		//		ResponseErr(model.ERR_NOT_ALL_VOTED, c)
		//		return
		//	}
		//}
		votes, _, err = gocache.GetVoteInfo(roomId)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			ResponseErr(model.ERR_NOT_ALL_VOTED, c)
			return
		}
		all := true
		for _, v := range votes {
			if !v {
				all = false
				break
			}
		}
		if !all {
			util.Error("ERR_NOT_ALL_VOTED")
			ResponseErr(model.ERR_NOT_ALL_VOTED, c)
			return
		}
	}
	//查看用户是不是凶手
	var isMurderer bool
	var maxNum int       //最大票数
	var mu string        //真正的凶手union_id
	var voteMu string    //投凶
	var coverVoteNum int //被投次数
	var score int        //投票得分
	cus := []string{}    //投票的凶手,会是多个凶手

	for _, v := range roomInfo.PlayerSlice {
		if v.UnionId == req.UnionId {
			isMurderer = v.Role.Murderer //获取是不是凶手
		}
		if v.Role.Murderer {
			mu = v.UnionId //记录真正的凶手
		}
	}
	util.Info("本用户是不是真正的凶手:%v 真正凶手:%v", isMurderer, mu)

	//最大值判断,出现两个最大投票情况如何处理
	for _, us := range roomInfo.UserInfoSlice {
		if us.CoverVoteNum >= maxNum { //未进行最大票数判断
			maxNum = us.CoverVoteNum
		}
	}
	util.Info("最大投票数:%v", maxNum)

	//投票得到的凶手
	for _, us := range roomInfo.UserInfoSlice {
		if us.CoverVoteNum >= maxNum { //未进行最大票数判断
			cus = append(cus, us.UnionId)
		}
		if us.UnionId == req.UnionId {
			voteMu = us.VoteUser           //获取投凶记录
			coverVoteNum = us.CoverVoteNum //记录被投次数
		}
	}
	util.Info("投凶是:%v 投票确定凶手是:%v", voteMu, cus)

	if isMurderer { //本人是凶手
		if coverVoteNum > 0 {
			//是凶手,被投出去 20分
			for _, id := range cus {
				if id == req.UnionId {
					score = 20
					break
				}
			}
			if score == 0 { //是凶手,被投，但没被投出去 80分
				score = 80
			}
			//if cu == req.UnionId {
			//	score = 20
			//} else {
			//	//是凶手,没被投出去 80分
			//	score = 80
			//}
		} else {
			score = 100 //没被投得100分
		}
	} else { //本人不是凶手
		if coverVoteNum > 0 {
			//if cu == req.UnionId { //被投出去了
			//	//投对了 被投出去了 40分
			//	if voteMu == mu {
			//		score = 40
			//	} else {
			//		score = 20 //投错了 被投出去了 20分
			//	}
			//} else {
			//	//被投但没被投出去 投对了 80分
			//	if voteMu == mu {
			//		score = 80
			//	} else {
			//		score = 60 //被投但没被投出去  投错了 60分
			//	}
			//}
			for _, id := range cus {
				if id == req.UnionId {
					//投对了 被投出去了 40分
					if voteMu == mu {
						score = 40
					} else {
						score = 20 //投错了 被投出去了 20分
					}
					break
				}
			}
			if score == 0 {
				if voteMu == mu { //被投但没被投出去 投对了 80分
					score = 80
				} else {
					score = 60 //被投但没被投出去  投错了 60分
				}
			}

		} else { //被投票次数是0
			//平民 没被投 投对了 100分
			if voteMu == mu {
				score = 100
			} else { //平民 没被投 投错了 80分
				score = 80
			}
		}
	}
	util.Info("获得得分:%v", score)

	for _, tmp := range roomInfo.PlayerSlice {
		if tmp.UnionId == userId {
			reportUrl = tmp.Role.Final
			roleId = tmp.Role.Id
		}
	}
	if len(reportUrl) == 0 {
		// Get Report Url
		scripId := roomInfo.ScriptId
		reportUrl, code = GetReportUrlByScriptId(scripId, roleId)
		if code != model.ERR_OK {
			//util.Logger(util.ERROR_LEVEL, "GameReport", "Get Report Url By Script Id Err")
			util.Info("GameReport Get Report Url By Script Id Err")
			ResponseErr(model.ERR_GET_REPORT, c)
			return
		}
	}

	//查看结案后房间游戏结束
	if roomInfo.Status != 2 && len(roomInfo.UnionIdSlice) == 0 {
		util.Debug("房间用户数:%v 房间状态:%v", len(roomInfo.UnionIdSlice), roomInfo.Status)
		roomInfo.Status = 2 //房间结束
		err = gocache.SetRoomInfo(roomId, roomInfo)
		if err != nil {
			util.Error("ERROR:%v", err.Error())
			ResponseErr(model.ERR_ENTERED_ROOM, c)
			return
		}
	}
	//用户查看结案后，从房间移出
	err = gocache.SetRoomInfo(roomId, roomInfo)
	if err != nil {
		util.Error("ERROR[%v]", err.Error())
		ResponseErr(model.ERR_ROOM_NOT_EXIST, c)
		return
	}

	if roomInfo.VoteFlag && !roomInfo.TopicFlag {
		//只有投凶
		lastScore = score
	} else if roomInfo.TopicFlag && !roomInfo.VoteFlag {
		//只有答题没有投凶:每道题的实际得分相加➗每道题的最高得分的总和✖️100位=最终得分
		topicScore, topicMaxScore, ok, err := gocache.GetRoomUserQuestionScore(roomId, userId)
		if err != nil {
			util.Error("ERROR[%v]", err.Error())
			ResponseErr(model.ERR_GET_REPORT, c)
			return
		}
		if !ok {
			util.Error("得分数据不存在:%v", ok)
			ResponseErr(model.ERR_GET_REPORT, c)
			return
		}
		util.Debug("得分:%v 最高分:%v", topicScore, topicMaxScore)
		if topicMaxScore == 0 || topicScore == 0 {
			lastScore = 0
			util.Debug("得分:%v", lastScore)
		} else {
			//lastScore = topicScore / topicMaxScore * 100 //等价于
			lastScore = topicScore * 100 / topicMaxScore //等价于规避0.5小数取整数据错误情况
			util.Debug("得分:%v", lastScore)
		}
	} else if roomInfo.TopicFlag && roomInfo.VoteFlag {
		//即有答题又有投凶:（投凶对应得分➕答题实际得分）➗（100➕每道题最高得分的总和）✖️100=最终得分
		topicScore, topicMaxScore, ok, err := gocache.GetRoomUserQuestionScore(roomId, userId)
		if err != nil {
			util.Error("ERROR[%v]", err.Error())
			ResponseErr(model.ERR_GET_REPORT, c)
			return
		}
		if !ok {
			util.Error("得分数据不存在:%v", ok)
			ResponseErr(model.ERR_GET_REPORT, c)
			return
		}
		util.Debug("得分:%v 最高分:%v", topicScore, topicMaxScore)
		if score == 0 && topicScore == 0 {
			lastScore = 0
			util.Debug("得分:%v", lastScore)
		} else {
			//lastScore = (score + topicScore) / (100 + topicMaxScore) * 100	//等价于
			lastScore = (score + topicScore) * 100 / (100 + topicMaxScore)
			util.Debug("得分:%v", lastScore)
		}
	} else {
		util.Error("不存在的情况,无需求,暂不处理...")
	}

	// Delete KV From Vote|Stage|Ap Bucket If Exist
	//boltdb.Delete([]byte(roomId), "VoteBucket")
	//boltdb.Delete([]byte(roomId), "StageBucket")
	//boltdb.Delete([]byte(roomId), "ApBucket")

	// Delete Room Cache
	//global.DeleteRoomCache(roomId)
	//global.DeleteClewCache(roomId)
	global.DeleteRoomDeleteTask(roomId)
	//gocache.DelRoomUser(roomId, )

	//util.Logger(util.INFO_LEVEL, "GameReport", reportUrl)

	// Return
	resp := model.GameReportResp{}
	//resp.Params = reportUrl
	//resp.Params.Score = score
	resp.Params.Score = lastScore
	resp.Params.ReportUrl = reportUrl
	util.Info("%+v", resp)

	//更新用户角色得分
	param := model.RoomRecordBase{}
	param.ScriptId = roomInfo.ScriptId
	param.RoomId = roomId
	param.UnionId = userId
	param.Score = resp.Params.Score
	go point.RoomUserDataUpdate(param)

	c.JSON(http.StatusOK, &resp)
}

// Func: Get Report Url By ScriptId
func GetReportUrlByScriptId(scriptId, roleId int) (string, int) {
	conn := gocache.RedisConnPool.Get()
	defer conn.Close()
	//_, _, final, ok, err := gocache.ConnGetApAndStory(conn, scriptId)
	info := model.ScriptPeopleInfo{}
	ok, err := gocache.ConnGetScriptPeopleInfo(conn, scriptId, roleId, &info)
	if err != nil {
		util.Error("ERROR:%v", err.Error())
		return info.Final, model.ERR_TASK_JSON
	}
	util.Debug("ok:%v len:%v", ok, len(info.Final))

	if ok && len(info.Final) > 0 {
		//util.Debug("OK:ap:%v story:%v", ap, story)
		return info.Final, model.ERR_OK
	}

	taskRequest := make(map[string]interface{})
	taskRequest["ScriptId"] = scriptId

	dbResult, err := global.Task.TaskJson(global.NewDBRequest("db.GameReport", taskRequest))
	if err != nil {
		return "", model.ERR_TASK_JSON
	}

	dbcode, dbparams := global.UnwrapObjectPackage(dbResult)

	switch dbcode {
	case global.ERR_DB_OK:
		return dbparams["Final"].(string), model.ERR_OK
	default:
		return "", model.ERR_DEFAULT
	}
}
