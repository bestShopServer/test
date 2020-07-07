package debug

import (
	"DetectiveMasterServer/model"
	"DetectiveMasterServer/release"
	"DetectiveMasterServer/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"log"
	"net/http"
)

type ScriptGetResp struct {
	Err    int64    `json:"err"`
	Msg    string   `json:"msg"`
	Params []Script `json:"params"`
}

type Script struct {
	About      string `json:"about"`
	Album      string `json:"album"`
	Author     string `json:"author"`
	Bad        int64  `json:"bad"`
	CreateTime string `json:"createTime"`
	Drama      string `json:"drama"`
	Good       int64  `json:"good"`
	Hot        int64  `json:"hot"`
	Id         int64  `json:"id"`
	Level      int64  `json:"level"`
	Name       string `json:"name"`
	Num        int64  `json:"num"`
	Played     bool   `json:"played"`
}

// TODO: Add Page & Add Script Limit In Config.json
type ScriptGetReq struct {
	ScriptId   int64   `json:"script_id"`
	LevelLower float64 `json:"level_lower"`
	LevelUpper float64 `json:"level_upper"`
	NumLower   float64 `json:"num_lower"`
	NumUpper   float64 `json:"num_upper"`
	Search     string  `json:"search"`
}

func ScriptGetEndpoint(c *gin.Context) {
	scriptJsonBytes := loadJsonToBytes("./script.json")

	// Check Script Json File
	if !gjson.ValidBytes(scriptJsonBytes) {
		util.Logger(util.ERROR_LEVEL, "script", "Invalid script json file.")
		release.ResponseErr(model.ERR_SERVER_ABNORMAL, c)
		return
	}

	var reqBody ScriptGetReq
	err := c.Bind(&reqBody)
	if err != nil {
		log.Println("Err: ParseRequest ", err)
		c.JSON(http.StatusOK, ScriptGetResp{
			Err: 4001,
			Msg: "错误参数格式",
		})
	} else {

		if !CheckReqParams(reqBody) {
			c.JSON(http.StatusOK, ScriptGetResp{
				Err: 4001,
				Msg: "错误参数格式",
			})
		} else {
			c.JSON(http.StatusOK, searchScript(reqBody, scriptJsonBytes))
		}
	}
}

// TODO: Update Check 1 - 10
// Func: Check Request Params
func CheckReqParams(req ScriptGetReq) bool {
	if req.LevelLower < 0 || req.LevelUpper < 0 ||
		req.NumLower < 0 || req.NumUpper < 0 {
		return false
	}
	return true
}

func searchScript(req ScriptGetReq, scriptJsonBytes []byte) ScriptGetResp {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var script Script
	var scriptList []Script
	ret := ScriptGetResp{}

	// Makeup Filter Path
	scriptIdPath := fmt.Sprintf("#[id==%d]#", req.ScriptId)
	levelLowerPath := fmt.Sprintf("#[level>=%f]#", req.LevelLower)
	levelUpperPath := fmt.Sprintf("#[level<=%f]#", req.LevelUpper)
	numLowerPath := fmt.Sprintf("#[num>=%f]#", req.NumLower)
	numUpperPath := fmt.Sprintf("#[num<=%f]#", req.NumUpper)
	nameSearchPath := fmt.Sprintf("#[name%%\"*%s*\"]#", req.Search)

	// Get Script Params
	result := gjson.Get(gjson.ParseBytes(scriptJsonBytes).String(), "params")

	// Filter Script Id
	if req.ScriptId != 0 {
		result = gjson.Get(result.String(), scriptIdPath)
	}

	// Filter Level
	if !util.Float64Equal(req.LevelLower, 0.0) {
		result = gjson.Get(result.String(), levelLowerPath)
	}
	if !util.Float64Equal(req.LevelUpper, 0.0) {
		result = gjson.Get(result.String(), levelUpperPath)
	}

	// Filter Num
	if !util.Float64Equal(req.NumLower, 0.0) {
		result = gjson.Get(result.String(), numLowerPath)
	}
	if !util.Float64Equal(req.NumUpper, 0.0) {
		result = gjson.Get(result.String(), numUpperPath)
	}

	// Filter Name
	result = gjson.Get(result.String(), nameSearchPath)
	for _, name := range result.Array() {
		if json.Unmarshal([]byte(name.String()), &script) == nil {
			scriptList = append(scriptList, script)
		}
	}

	ret.Err = 0
	ret.Msg = ""
	ret.Params = scriptList
	return ret
}
