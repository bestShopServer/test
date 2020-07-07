package model

type ErrResp struct {
	Err int    `json:"err"`
	Msg string `json:"msg"`
}

type RoomNewResp struct {
	ErrResp
	Params RoomNewRespParam `json:"params"`
}

type RoomInfoResp struct {
	ErrResp
	Params RoomNewRespParam `json:"params"`
}

type RoomEnterResp struct {
	ErrResp
	Params RoomNewRespParam `json:"params"`
}

type RoomNewRespParam struct {
	RoomId   string   `json:"room_id"`
	RoomInfo RoomInfo `json:"room_info"`
}

type WxCreateResp struct {
	ErrResp
	Params WxCreateParams `json:"params"`
}

type WxCreateParams struct {
	Path     string `json:"path"`
	RoomId   string `json:"room_id"`
	RoomName string `json:"room_name"`
}

type WxJoinResp struct {
	ErrResp
	Params WxJoinParams `json:"params"`
}

type WxJoinParams struct {
	Path string `json:"path"`
}

type RoleChooseResp struct {
	ErrResp
	Params string `json:"params"`
}

type GameStartResp struct {
	ErrResp
	Params string `json:"params"`
}

type GameInfoParams struct {
	GameInfo
	Player    PlayerInfo  `json:"player"`
	GameStage interface{} `json:"game_stage"`
}

type GameInfoResp struct {
	ErrResp
	Params GameInfoParams `json:"params"`
}

type RoomExitResp struct {
	ErrResp
	Params string `json:"params"`
}

type WxExitResp struct {
	ErrResp
	Params string `json:"params"`
}

type ClewOpenResp struct {
	ErrResp
	Params string `json:"params"`
}

type ClewPubResp struct {
	ErrResp
	Params string `json:"params"`
}

type StageNextParams struct {
	GameStage interface{} `json:"game_stage"`
}

type StageNextResp struct {
	ErrResp
	Params StageNextParams `json:"params"`
}

type GameVoteResp struct {
	ErrResp
	Params string `json:"params"`
}

type GameReportResp struct {
	ErrResp
	//Params string `json:"params"`
	Params GameReport `json:"params"`
}

type ScriptGetResp struct {
	Err    int64    `json:"err"`
	Msg    string   `json:"msg"`
	Params []Script `json:"params"`
}

type ReconnectResp struct {
	ErrResp
	Params ReconnectInfo `json:"params"`
}

type GameClewResp struct {
	ErrResp
	Params GameClew `json:"params"`
}

//获取code
type Code2SessionParam struct {
	OpenId     string `json:"openid"`      //用户唯一标识
	SessionKey string `json:"session_key"` //会话密钥
	Token      string `json:"token"`
	UnionId    string `json:"unionid"`
	//Phone      string `json:"phone"`      //是否需要获取手机号
}
type Code2SessionResp struct {
	ErrResp
	Code2SessionParam `json:"params"`
}

//获取code
type ImGenSigParam struct {
	UserSig string `json:"user_sig"` //IM UserSig
}
type ImGenSigResp struct {
	ErrResp
	ImGenSigParam `json:"params"`
}

type watermark struct {
	AppId     string `json:"appid"`
	Timestamp int64  `json:"timestamp"`
}

//获取用户微信手机号
type WxUserPhoneResp struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
	watermark       `json:"watermark"`
	Token           string `json:"token"`
}

//返回前端用户表示

//获取用户微信手机号
type GetWxUserPhoneResp struct {
	ErrResp
	Res WxUserPhoneResp `json:"params"`
}

//微信支付返回信息
type WxUnifiedReturn struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	AppId      string `xml:"appid"`  //小程序ID
	MchId      string `xml:"mch_id"` //商户号
	DeviceInfo string `xml:"device_info"`
	NonceStr   string `xml:"nonce_str"`   //随机字符串
	Sign       string `xml:"sign"`        //签名
	ResultCode string `xml:"result_code"` //处理结果
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"` //错误信息描述
	TradeType  string `xml:"trade_type"`   //交易类型 小程序取值如下：JSAPI
	PrepayId   string `xml:"prepay_id"`    //预支付交易会话标识
	CodeUrl    string `xml:"code_url"`     //二维码链接
}

//获取支付参数
type MarketPayParams struct {
	AppId     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
	UserScript
}

//获取用户微信手机号
type MarketPayParamsResp struct {
	ErrResp
	Res MarketPayParams `json:"params"`
}

//响应微信回调数据
type XmlResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
}

//剧本问题
type ScriptQuestionGetResp struct {
	ErrResp
	//Params []Questions `json:"params"`
	Params QuestionAndAnswer `json:"params"`
}

//剧本问题接单
type AnswerResp struct {
	ErrResp
	Params Ending `json:"params"`
}

//审核接口
type Audit struct {
	IsShow int `json:"is_show"`
}

//是否展示审核页面
type AuditResp struct {
	ErrResp
	Params Audit `json:"params"`
}

//用户房间号
type UserRoomBase struct {
	RoomId string `json:"room_id"`
	MemberBaseInfo
}

//响应用户房间号
type UserRoomResp struct {
	ErrResp
	Params UserRoomBase `json:"params"`
}

//用户创建房间号
type RoomCreateBase struct {
	RoomId   string  `json:"room_id"`
	RoomName string  `json:"room_name"`
	Num      int     `json:"num"`
	Price    float64 `json:"price"`
}

//响应用户创建房间号
type RoomCreateResp struct {
	ErrResp
	Params RoomCreateBase `json:"params"`
}

//响应用户生成订单号
type OrderCreateParams struct {
	OrderNo string `json:"order_no"`
	InvTime string `json:"inv_time"` //失效时间
}

//响应用户是否会员
type UserIsMemberResp struct {
	ErrResp
	Params MemberBaseInfo `json:"params"`
}

//响应用户生成订单号
type OrderCreateResp struct {
	ErrResp
	Params OrderCreateParams `json:"params"`
}
