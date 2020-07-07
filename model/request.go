package model

import "encoding/xml"

type CommonReq struct {
	//OpenId  string `json:"open_id"`
	RoomId  string `json:"room_id"`
	UnionId string `json:"union_id"`
}

type UserLoginReq struct {
	Code        string `json:"code"`
	EncryptData string `json:"encrypt_data"`
	IV          string `json:"iv"`
}

type WxLoginReq struct {
	OpenId  string `json:"open_id"`
	UnionId string `json:"union_id"`
}

type RoomNewReq struct {
	//OpenId   string `json:"open_id"`
	UnionId  string `json:"union_id"`
	ScriptId int    `json:"script_id"`
}

type WxCreateReq struct {
	//OpenId     string `json:"open_id"`
	UnionId    string `json:"union_id"`
	ScriptName string `json:"script_name"`
	Flag       int    `json:"flag"`
}

type RoomEnterReq struct {
	CommonReq
}

type WxJoinReq struct {
	CommonReq
	UserName string `json:"user_name"`
}

type RoomInfoReq struct {
	CommonReq
}

type RoomExitReq struct {
	CommonReq
}

type WxExitReq struct {
	//OpenId string `json:"open_id"`
	UnionId string `json:"union_id"`
}

type WxKickReq struct {
	UnionId     string `json:"union_id"`
	KickUnionId string `json:"kick_union_id"`
}

type RoleChooseReq struct {
	CommonReq
	RoleId int `json:"role_id"`
}

type GameStartReq struct {
	CommonReq
}

type RoomDeleteReq struct {
	CommonReq
}

type GameInfoReq struct {
	CommonReq
}

type ClewOpenReq struct {
	CommonReq
	Key      string `json:"key"`
	Index    []int  `json:"index"`
	Password string `json:"password"`
}

type ClewPubReq struct {
	CommonReq
	Key   string `json:"key"`
	Index []int  `json:"index"`
}

type StageNextReq struct {
	CommonReq
}

type GameVoteReq struct {
	CommonReq
	RoleId int `json:"role_id"`
}

type GameReportReq struct {
	CommonReq
}

type GameReconnectReq struct {
	CommonReq
	UserName string `json:"user_name"`
}

type GameClewReq struct {
	CommonReq
}

type ScriptGetReq struct {
	ScriptId   int    `json:"script_id"`
	LevelLower int    `json:"level_lower"`
	LevelUpper int    `json:"level_upper"`
	NumLower   int    `json:"num_lower"`
	NumUpper   int    `json:"num_upper"`
	Search     string `json:"search"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
}

//用户auth.code2Session
type Code2SessionReq struct {
	Code string `json:"code"  form:"code"`
	//Location string `json:"location"`
}

//用户auth.code2Session
type GenSigReq struct {
	OpenId string `json:"open_id"`
}

//获取用户微信手机号
type GetWxUserPhoneReq struct {
	EncryptedData string `json:"encryptedData"`
	Iv            string `json:"iv"`
	SessionKey    string `json:"session_key"`
	City          string `json:"city"`
	Province      string `json:"province"`
	Country       string `json:"country"`
	Gender        int    `json:"gender"`
	NickName      string `json:"nick_name"`
	OpenId        string `json:"open_id"`
	UnionId       string `json:"union_id"`
	Phone         string `json:"-"`
}

//获取支付参数
type GetWxPayParamsReq struct {
	OpenId  string `json:"open_id"`
	UnionId string `json:"union_id"`
	SId     int    `json:"sid"`
	//Num     int     `json:"num"`
	//Price   float64 `json:"price"`
	RoomId int `json:"room_id"`
	Flag   int `json:"flag"` //1请大家2AA付款
}

type UnionId struct {
	UnionId string `json:"union_id"`
}

//剧本用户付款情况
type ScriptUserCostReq struct {
	ScriptId int      `json:"script_id"`
	OpenId   string   `json:"open_id"`
	UnionIds []string `json:"union_id"`
}

//微信统一支付参数
type WxUnifiedPayXml struct {
	XMLName        xml.Name `xml:"xml"`
	AppId          string   `xml:"appid"`      //小程序ID
	Body           string   `xml:"body"`       //商品描述
	MchId          string   `xml:"mch_id"`     //商户号
	NonceStr       string   `xml:"nonce_str"`  //随机字符串
	NotifyUrl      string   `xml:"notify_url"` //异步接收微信支付结果通知的回调地址，通知url必须为外网可访问的url，不能携带参数。
	OpenId         string   `xml:"openid"`
	OutTradeNo     string   `xml:"out_trade_no"`     //商户订单号
	Sign           string   `xml:"sign"`             //签名
	SpbillCreateIp string   `xml:"spbill_create_ip"` //终端IP
	TotalFee       int      `xml:"total_fee"`        //订单总金额，单位为分
	TradeType      string   `xml:"trade_type"`       //交易类型 小程序取值如下：JSAPI
}

//微信支付回调
type WxPayNotifyRrlReq struct {
	ReturnCode         string `xml:"return_code"`
	ReturnMsg          string `xml:"return_msg"`
	AppId              string `xml:"appid"`
	MchId              string `xml:"mch_id"`
	DeviceInfo         string `xml:"device_info"`
	NonceStr           string `xml:"nonce_str"`
	Sign               string `xml:"sign"`
	ResultCode         string `xml:"result_code"`
	ErrCode            string `xml:"err_code"`
	ErrCodeDes         string `xml:"err_code_des"`
	OpenId             string `xml:"openid"`
	IsSubscribe        string `xml:"is_subscribe"`
	TradeType          string `xml:"trade_type"`
	BankType           string `xml:"bank_type"`
	TotalFee           int    `xml:"total_fee"`
	SettlementTotalFee int    `xml:"settlement_total_fee"`
	FeeType            string `xml:"fee_type"`
	CashFee            int    `xml:"cash_fee"`
	CashFeeType        string `xml:"cash_fee_type"`
	CouponFee          int    `xml:"coupon_fee"`
	CouponCount        int    `xml:"coupon_count"`

	TransactionId string `xml:"transaction_id"`
	OutTradeNo    string `xml:"out_trade_no"`
	Attach        string `xml:"attach"`
	TimeEnd       string `xml:"time_end"`
}

type PayParams struct {
	WxUnifiedReturn
	OutTradeNo string
	PayAmt     float64
	TotalFee   int
	Content    string
	PayContent string
}

//获取剧本问题
type ScriptQuestionGetReq struct {
	ScriptId int    `json:"script_id"`
	UnionId  string `json:"union_id"`
	PId      int    `json:"pid"`
}

//获取剧本问题
type AnswerReq struct {
	ScriptId int            `json:"script_id"`
	UnionId  string         `json:"union_id"`
	Answer   []QuestionsOpt `json:"answer"`
}

//获取剧本问题解答
type QuestionEndingReq struct {
	ScriptId int    `json:"script_id"`
	UnionId  string `json:"union_id"`
	Opt      string `json:"opt"`
}

//获取剧本问题
type GameScoreReq struct {
	UnionId     string `json:"union_id"`     //用户unionid
	ScriptScore int    `json:"script_score"` //剧本评分
	GameScore   int    `json:"game_score"`   //游戏评分
	RoomId      int    `json:"room_id"`
	ScriptId    int    `json:"script_id"`
}

//获取剧本问题
type AuditReq struct {
	Version int `json:"version"` //用户unionid
}

//获取剧本问题
type PublicOrderSyncReq struct {
	UnionId    string  `json:"union_id"`    //用户unionid
	ScriptName string  `json:"script_name"` //剧本名称
	Price      float64 `json:"price"`       //剧本付款价格
	OrderNo    string  `json:"order_no"`    //订单号
}

//获取剧本问题
type OrderStateReq struct {
	UnionId    string `json:"union_id"`    //用户unionid
	ScriptName string `json:"script_name"` //剧本名称
}

//查询用户房间号
type WxUserRoomReq struct {
	UnionId string `json:"union_id"`
}

//创建房间
type WxRoomCreateReq struct {
	UnionId  string `json:"union_id"`
	ScriptId int    `json:"script_id"`
	Flag     int    `json:"flag"`
}

//用户投降
type UserSurrenderReq struct {
	UnionId string `json:"union_id"`
	RoomId  string `json:"room_id"`
}

//用户是否为会员
type UserIsMemberReq struct {
	UnionId string `json:"union_id"`
}

//生成订单
type OrderCreateReq struct {
	UnionId    string `json:"union_id"`    //用户unionid
	ActivityId int32  `json:"activity_id"` //活动ID
}
