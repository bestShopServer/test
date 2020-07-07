package main

import (
	"DetectiveMasterServer/config"
	"DetectiveMasterServer/global"
	"DetectiveMasterServer/gocache"
	"DetectiveMasterServer/release"
	"DetectiveMasterServer/release/game"
	"DetectiveMasterServer/release/order"
	"DetectiveMasterServer/release/pay"
	"DetectiveMasterServer/release/question"
	"DetectiveMasterServer/release/room"
	"DetectiveMasterServer/release/user"
	"DetectiveMasterServer/util"
	"DetectiveMasterServer/websocket"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"net"
	"net/http"
	"runtime"
	"time"
)

var conf = config.GetConfig()

func startWebSocketClient() {
	fmt.Println("startWebSocketClient ...")

	//websocket.WS = websocket.DialWebSocket()
	//if websocket.WS != nil {
	//websocket.WsConn = websocket.DialWebSocket()
	//if websocket.WsConn.Conn != nil {
	//	websocket.SendHandShakeMessage()
	//}
	util.Info("startWebSocketClient ok ")
	fmt.Println("startWebSocketClient ok")

	websocket.DialRedisChannel()
}

//func startHttpServer() {
//	fmt.Println("startHttpServer ...")
//	gin.DisableConsoleColor()
//
//	// Create Access.log & Set Gin Log Writer
//	//f, _ := os.Create("access.log")
//	//gin.DefaultWriter = io.MultiWriter(f)

func InitRouter() *gin.Engine {
	fmt.Println("InitRouter ...")
	gin.DisableConsoleColor()
	router := gin.Default()
	router.Use(gin.Recovery())

	// Test Group: debug
	group := router.Group("/")
	{
		group.POST("/script.get", release.ScriptGet)
		//group.POST("/user.login", release.UserLogin)
		//group.POST("/room.new", release.RoomNew)
		group.POST("/room.info", release.RoomInfo)
		//group.POST("/room.enter", release.RoomEnter)
		//group.POST("/room.exit", release.RoomExit)
		group.POST("/role.choose", release.RoleChoose)
		group.POST("/game.start", release.GameStart)
		group.POST("/game.info", release.GameInfo)
		group.POST("/game.clew", release.GameClew)
		group.POST("/clew.open", release.ClewOpen)
		group.POST("/clew.pub", release.ClewPub)
		group.POST("/stage.next", release.StageNext)
		group.POST("/game.vote", release.GameVote)
		group.POST("/game.report", release.GameReport)
		group.POST("/game.reconnect", release.GameReconnect)
		group.POST("/wx.login", release.WxLogin)

		group.POST("/wx.create", release.WxCreate)
		//group.POST("/wx.join", release.WxJoin)
		group.POST("/wx.join", room.WxRoomJoinV2) //强制加入新房间

		group.POST("/wx.exit", release.WxExit)
		group.POST("/wx.code", release.WxCode2Session)
		group.POST("/wx.phone", release.WxUserPhone)
		group.POST("/wx.genSig", release.WxGenSig)
		group.POST("/wx.pay_params", pay.WxPayParams)     //支付参数
		group.POST("/wx.notify_url", pay.PostWxNotifyRrl) //支付回调
		group.POST("/wx.kick", release.WxKick)            //踢出房间

		//查询问题
		group.POST("/script.question", question.ScriptQuestionGet)
		//提交问题答案
		group.POST("/question.score", question.CalculationScore)

		//游戏评分
		group.POST("/game.score", game.GameScore)
		group.POST("/audit", game.Audit) //判断版本号

		//删除房间
		group.POST("/room.delete", room.RoomDelete)
		//公众号退出房间
		group.POST("/wechat.exit", release.WechatExit)

		//公众号同步剧本订单
		group.POST("/order.sync", order.PublicOrderSync)
		//查询用户订单支付状态
		group.POST("/order.state", order.PostOrderState)

		//测试同步订单信息到公众号服务
		group.POST("/order.sync.public", pay.TestOrderSyncPublic)
		//测试公众号状态
		group.POST("/order.state.public", pay.TestOrderStatePublic)

		//v1.1.1
		//查询用户房间号
		group.POST("/user.room", user.WxUserRoomBase)
		//根据剧本ID创建房间
		group.POST("/room.create", room.RoomCreate)

		//投降
		group.POST("/surrender", user.UserSurrender)

		//判断是否是会员
		group.POST("/ismember", user.UserIsMember)

		//生成订单
		//group.POST("/order.create", order.PostOrderCreate)

		//group.POST("/ping", game.GameReconnect)

	}
	util.Info("路由设置成功!")
	//fmt.Println("startHttpServer run:", conf.ServerHost)
	//router.Run(conf.ServerHost)
	return router
}

//启动服务
func startHttpServer() {

	routersInit := InitRouter()
	maxHeaderBytes := 1 << 20
	//endPoint := fmt.Sprintf("127.0.0.1:%v", conf.ServerHost)
	//server := &http.Server{
	//	Addr:           endPoint,
	//	Handler:        routersInit,
	//	ReadTimeout:    30,
	//	WriteTimeout:   30,
	//	MaxHeaderBytes: maxHeaderBytes,
	//}
	//fmt.Printf("[info] start http server listening %s\n", endPoint)
	//util.Info("服务正常启动 listening[%v]\n", endPoint)
	//
	//err := server.ListenAndServe()
	//if err != nil {
	//	fmt.Printf("启动服务ERROR[%v]", err.Error())
	//}
	server := &http.Server{
		Handler:        routersInit,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: maxHeaderBytes,
	}
	listen, err := net.Listen("tcp4", conf.ServerHost)
	if err != nil {
		util.Error("Failed to listen,err:%s", err.Error())
		panic(err)
	}
	server.SetKeepAlivesEnabled(false)
	err = server.Serve(listen)

	util.Info("服务正常启动[%v]", conf.ServerHost)
}

// Func: Init Bolt Bucket
func initEnvironment() {
	fmt.Println("initEnvironment ...")
	//// Create Database If Not Exist
	//boltdb.DB = boltdb.CreateDatabase()
	//
	//// Create Bucket If Not Exist
	//boltdb.CreateBucket("RoomBucket")
	//boltdb.CreateBucket("GameBucket")
	//boltdb.CreateBucket("StageBucket")
	//boltdb.CreateBucket("UserBucket")
	//boltdb.CreateBucket("ApBucket")
	//boltdb.CreateBucket("VoteBucket")
	fmt.Println("initEnvironment ok")

}

func startZmqClient() {
	fmt.Println("startZmqClient ...")
	//var err error

	global.Task = &util.TaskObj{
		Addr:     conf.DBAddr,
		Timeout:  conf.DBTimeout,
		Host:     conf.DBHost,
		MaxQueue: conf.MaxQueue,
	}
	util.Info("db addr:", conf.DBAddr)
	//global.DBChannel = goczmq.NewDealerChanneler(fmt.Sprintf("%v", conf.DBAddr))
	//util.Sock, err = goczmq.NewDealer(fmt.Sprintf("%v", conf.DBAddr))
	//if err != nil {
	//	util.Error("ERROR:%v", err.Error())
	//}
	util.InitZeroMQClient(conf.DBAddr, conf.DBTimeout)
	//zeromq.InitMQSockPool(conf.DBAddr, 100)
	util.Info("Zero MQ NewDealer Ok!")
}

//初始化
func init() {
	fmt.Printf("init ...")
	util.InitLog()

	//打开Redis连接池
	err := gocache.Setup()
	if err != nil {
		util.Error("打开Redis连接池失败[%v]", err.Error())
	}

	//添加定时任务
	c := cron.New()
	//second min hour day month week
	err = c.AddFunc("0 0 0 * * *", util.BackLogFile)
	if err != nil {
		util.Error("定时任务执行失败[%v]", err.Error())
	}
	//备份日志
	//second min hour day month week
	err = c.AddFunc("0 */15 * * * *", util.CheckLogFileSize)
	if err != nil {
		util.Error("定时任务执行失败[%v]", err.Error())
	}

	c.Start()
}

func main() {
	// 利用cpu多核来处理http请求，这个没有用go默认就是单核处理http的，这个压测过了，请一定要相信我
	runtime.GOMAXPROCS(runtime.NumCPU())
	//initEnvironment()
	startWebSocketClient()
	startZmqClient()

	startHttpServer()
	util.Info("初始化mongodb...")
	//global.InitMongoDB()
}
