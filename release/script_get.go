package release

import (
	"DetectiveMasterServer/config"
	"DetectiveMasterServer/global"
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
)

func ScriptGet(c *gin.Context) {
	fmt.Println("ScriptGet ...")

	var req model.ScriptGetReq
	//fmt.Println("req:", req)

	err := c.Bind(&req)

	util.Logger(util.INFO_LEVEL, "ScriptGet", "ScriptGet ...")

	// Optional Fields List
	optionalFields := []string{"ScriptId", "LevelLower", "LevelUpper", "NumLower", "NumUpper", "Search", "Page", "Limit"}
	fmt.Println("err:", err)
	// Check Param
	if !CheckParams(req, "ScriptGet", err, optionalFields) {
		util.Error("ERR_WRONG_FORMAT:%+v", req)
		ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}
	util.Info("请求参数:%+v", req)
	if req.ScriptId == 0 && req.LevelLower == 0 && req.LevelUpper == 0 &&
		req.NumLower == 0 && req.NumUpper == 0 && len(req.Search) == 0 {
		ResponseErr(model.ERR_WRONG_FORMAT, c)
		return
	}

	conf := config.GetConfig()
	FixParams(&req, conf)

	scriptList, code := GetScripts(req)
	if code != model.ERR_OK {
		util.Logger(util.ERROR_LEVEL, "ScriptGet", "Get Scripts Err")
		ResponseErr(model.ERR_GET_SCRIPTS, c)
		return
	}
	util.Debug("scriptList:%+v", len(scriptList))

	ScriptGetResp := model.ScriptGetResp{}
	ScriptGetResp.Params = scriptList
	c.JSON(http.StatusOK, &ScriptGetResp)
}

// Func: Check Request Params
func FixParams(req *model.ScriptGetReq, conf *config.Config) {
	if req.LevelLower < 1 || req.LevelLower > 10 {
		req.LevelLower = 1
	}
	if req.LevelUpper < 1 || req.LevelUpper > 10 {
		req.LevelUpper = 10
	}
	if req.LevelLower > req.LevelUpper {
		req.LevelLower = 1
		req.LevelUpper = 10
	}
	if req.NumLower < 1 || req.NumLower > 12 {
		req.NumLower = 1
	}
	if req.NumUpper < 1 || req.NumUpper > 12 {
		req.NumUpper = 12
	}
	if req.NumLower > req.NumUpper {
		req.NumLower = 1
		req.NumUpper = 12
	}
	if req.Page < 1 {
		req.Page = 1
	}
	req.Limit = conf.Limit
}

func randHot() int {
	hot := rand.Intn(10)
	if hot == 0 {
		return 1
	} else {
		return hot
	}
}

func GetScripts(req model.ScriptGetReq) (scriptList []model.Script, code int) {

	//var scriptList []model.Script

	taskRequest := make(map[string]interface{})
	taskRequest["ScriptId"] = req.ScriptId
	taskRequest["LevelLower"] = req.LevelLower
	taskRequest["LeverUpper"] = req.LevelUpper
	taskRequest["NumLower"] = req.NumLower
	taskRequest["NumUpper"] = req.NumUpper
	taskRequest["Search"] = req.Search
	taskRequest["Page"] = req.Page
	taskRequest["Limit"] = req.Limit

	dbResult, err := global.Task.TaskJson(global.NewDBRequest("db.ScriptMiniGet", taskRequest))
	if err != nil {
		return scriptList, model.ERR_TASK_JSON
	}

	dbcode, dbparams := global.UnwrapArrayPackage(dbResult)

	switch dbcode {
	case global.ERR_DB_OK:
		util.Info("剧本个数:%v", len(dbparams))
		//scriptList = make([]model.Script, len(dbparams))
		for _, p := range dbparams {
			var script model.Script
			pm := p.(map[string]interface{})
			if pm["Id"] != nil {
				script.Id = int(pm["Id"].(float64))
			}
			if pm["Name"] != nil {
				script.Name = pm["Name"].(string)
			}
			if pm["About"] != nil {
				script.About = pm["About"].(string)
			}
			if pm["Album"] != nil {
				script.Album = pm["Album"].(string)
			}
			if pm["Author"] != nil {
				script.Author = pm["Author"].(string)
			}
			if pm["Bad"] != nil {
				script.Bad = int(pm["Bad"].(float64))
			}
			if pm["CreateTime"] != nil {
				script.CreateTime = pm["CreateTime"].(string)
			}
			if pm["Drama"] != nil {
				script.Drama = pm["Drama"].(string)
			}
			if pm["Level"] != nil {
				script.Good = int(pm["Level"].(float64))
			}
			//util.Debug("num:%v", pm["Num"])
			if pm["Num"] != nil {
				script.Num = int(pm["Num"].(float64))
				//util.Debug("人数:%v", script.Num)
			}
			if pm["Price"] != nil {
				//script.Price = float64(pm["Price"].(float64))
				script.Price = int(pm["Price"].(float64))
			}
			if pm["ShowCase"] != nil {
				script.ShowCase = int(pm["ShowCase"].(float64))
			}
			//题目标识
			if pm["TopicFlag"] != nil {
				script.TopicFlag = pm["TopicFlag"].(bool)
			}
			//探索标识
			if pm["ExploreFlag"] != nil {
				script.ExploreFlag = pm["ExploreFlag"].(bool)
			}
			//投票标识
			if pm["VoteFlag"] != nil {
				script.VoteFlag = pm["VoteFlag"].(bool)
			}

			//剧本轮次
			if pm["Round"] != nil {
				script.Round = int(pm["Round"].(float64))
				util.Debug("轮次:%v", script.Round)
			}

			script.Hot = randHot()

			//util.Debug("人数:%v", script.Num)
			//util.Debug("剧本:%+v", script)
			scriptList = append(scriptList, script)
			//util.Debug("剧本:%+v", scriptList)
		}
		code = model.ERR_OK
		//util.Debug("剧本:%+v", scriptList)
		return scriptList, code
	default:
		code = model.ERR_DEFAULT
		return scriptList, code
	}
}
