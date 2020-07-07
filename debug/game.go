package debug

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
)

type GameInfoReq struct {
	RoomId string `json:"room_id"`
	OpenId string `json:"open_id"`
}

type GameInfoResp struct {
	Err    int64    `json:"err"`
	Msg    string   `json:"msg"`
	Params GameInfo `json:"params"`
}

type GameInfo struct {
	Story string   `json:"story"`
	Task  MainTask `json:"task"`
	About string   `json:"about"`
	Clew  []Clew   `json:"clew"`
}

type MainTask struct {
	Main string     `json:"main"`
	Side []SideTask `json:"sub"`
}

type SideTask struct {
	Id      int64    `json:"id"`
	Caption string   `json:"caption"`
	Options []Option `json:"options"`
}

type Option struct {
	Id     int64  `json:"id"`
	Option string `json:"option"`
}

type Clew struct {
	Round int64          `json:"round"`
	Key   map[string]Key `json:"key"`
}

type Key struct {
	Album   string      `json:"album"`
	KeyList []KeyDetail `json:"key_list"`
}

type KeyDetail struct {
	Id       int64          `json:"id"`
	Ap       int64          `json:"ap"`
	Type     int64          `json:"type"`
	Content  string         `json:"content"`
	Password string         `json:"password"`
	Status   int64          `json:"status"`
	Sub      []SubKeyDetail `json:"sub"`
}

type SubKeyDetail struct {
	Id       int64  `json:"id"`
	Ap       int64  `json:"ap"`
	Type     int64  `json:"type"`
	Content  string `json:"content"`
	Password string `json:"password"`
	Status   int64  `json:"status"`
}

func GameInfoEndpoint(c *gin.Context) {
	gameInfoResp := GameInfoResp{}
	loadJson("./game.json", &gameInfoResp)

	var reqBody GameInfoReq
	err := c.Bind(&reqBody)

	if err != nil {
		log.Println("Err: ParseRequest ", err)
		c.JSON(http.StatusOK, GameInfoResp{
			Err: 4001,
			Msg: "错误参数格式",
		})
	} else {
		c.JSON(http.StatusOK, gameInfoResp)
	}
}
