package model

import "time"

type UserInfo struct {
	UnionId      string `json:"union_id"`
	Member       int    `json:"member"` //会员登记0非会员1会员
	Name         string `json:"name"`
	IsPay        int    `json:"is_pay"`         //付款标识 1 未付款 2 已付款
	CoverVoteNum int    `json:"cover_vote_num"` //被投票次数，被投最多的就是凶手了
	VoteUser     string `json:"vote_user"`      //投票凶手
	//IsSeeReport	 int 	`json:"is_see_report"`	//是否查看结案 0未看 1查看
	Surrender bool `json:"surrender"` //投降标识true false
}

type ScreenQuestion struct {
	Ssqid    int            `json:"ssqid"`
	Question string         `json:"question"`
	Flag     int            `json:"flag"` //是否有问题1是2否
	Answers  []ScreenAnswer `json:"answers"`
}
type ScreenAnswer struct {
	Ssaid  int    `json:"ssaid"`
	Answer string `json:"answer"`
	Screen string `json:"screen"`
}

type ScreenInfo struct {
	ScreenNum int `json:"screen_num"`
	//Content   []string `json:"screen_contents"`
	Content []ScreenQuestion `json:"screen_contents"`
}

type PlayerInfo struct {
	//OpenId string   `json:"open_id"`
	UnionId    string       `json:"union_id"`
	Role       RoleInfo     `json:"role"`
	VoteRoleId int          `json:"vote_role_id"` //投角色为凶手
	VoteResult bool         `json:"vote_result"`  //投票结果对不对
	Screens    []ScreenInfo `json:"screens"`
}

type RoleInfo struct {
	Id       int    `json:"id"`
	Album    string `json:"album"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Sex      int    `json:"sex"`
	Tall     int    `json:"tall"`
	About    string `json:"about"`
	Vote     bool   `json:"vote"`     //可否被投票
	Pofround string `json:"pofround"` //选择人物花费的AP点
	Choice   bool   `json:"choice"`   //任务可否被选
	Murderer bool   `json:"murderer"` //是否是凶手
	//KeywordFlag bool   `json:"keyword_flag"` //关键字标识
	//TopicFlag   bool   `json:"topic_flag"`   //问题标识
	Final string `json:"final"`
}

//房间信息
type RoomInfo struct {
	ScriptId int     `json:"script_id"`
	Round    int     `json:"round"` //剧本轮
	Num      int     `json:"num"`
	Price    float64 `json:"price"`
	Flag     int     `json:"flag"` //1请大家2AA付款
	Owner    string  `json:"owner"`
	Status   int     `json:"status"` //房间状态 2游戏结束
	//OpenIdSlice []string     `json:"open_id_slice"`
	UserInfoSlice []UserInfo   `json:"user_info_slice"`
	UnionIdSlice  []string     `json:"union_id_slice"`
	PlayerSlice   []PlayerInfo `json:"player_slice"`
	ExploreFlag   bool         `json:"explore_flag"` //探索标识
	TopicFlag     bool         `json:"topic_flag"`   //问题标识
	VoteFlag      bool         `json:"vote_flag"`    //投票标识
}

////关键字
//type KeyWord struct {
//	Keyword string `json:"keyword"`
//	Content string `json:"content"`
//}
//
//type GameClewContent struct {
//	T int    `json:"t"`
//	D string `json:"d"`
//}
//
//type GameSubClew struct {
//	UnionId  string            `json:"union_id"`
//	Id       int               `json:"id"`
//	Ap       int               `json:"ap"`
//	Type     int               `json:"type"`
//	Content  []GameClewContent `json:"content"`
//	Password string            `json:"password"`
//	Status   int               `json:"status"`
//	Sub      []GameSubClew     `json:"sub"`
//}
//type GameClewValue struct {
//	Album    string        `json:"album"`
//	KeyList  []GameSubClew `json:"key_list"`
//	TotalNum int           `json:"total_num"`
//	FetchNum int           `json:"fetch_num"`
//}

type CacheGameClew struct {
	Id    int             `json:"id"`
	Round int             `json:"round"`
	Sub   []CacheGameClew `json:"sub"`
}

type CacheGameClews struct {
	Name string          `json:"name"`
	Clew []CacheGameClew `json:"clew"`
}
type CacheClew struct {
	Key map[string]CacheGameClew `json:"key"`
}

type GameClews struct {
	Round int         `json:"round"`
	Key   interface{} `json:"key"`
}

//关键字回复轮次数组
type KeyWord struct {
	Skid    int    `json:"skid"`
	Round   int    `json:"round"`
	Keyword string `json:"keyword"`
	Content string `json:"content"`
	Album   string `json:"album"`
	SpId    int    `json:"-"`
	EndId   int    `json:"end_id"`
	Resume  string `json:"resume"`
}

//关键字关系
type KeyWordRelation struct {
	Sid    int    `json:"sid"`
	Pid    int    `json:"pid"`
	Round  int    `json:"round"`
	Skids  string `json:"skids"`
	EndId  int    `json:"end_id"`
	Resume string `json:"resume"`
	Remark string `json:"remark"`
}

type GameInfo struct {
	Story           string            `json:"story"`
	Task            MainTask          `json:"task"`
	About           string            `json:"about"`
	Ap              []int             `json:"ap"`
	Clew            []Clew            `json:"clew"`
	Explore         []KeyWord         `json:"explore"`
	ExploreRelation []KeyWordRelation `json:"explore_relation"`
}

type MainTask struct {
	Main string     `json:"main"`
	Side []SideTask `json:"sub"`
}

type SideTask struct {
	Id      int      `json:"id"`
	Caption string   `json:"caption"`
	Options []Option `json:"options"`
}

type Option struct {
	Id     int    `json:"id"`
	Option string `json:"option"`
}

type Clew struct {
	Round int            `json:"round"`
	Key   map[string]Key `json:"key"`
}

type Key struct {
	Album    string          `json:"album"`
	KeyList  []MainKeyDetail `json:"key_list"`
	TotalNum int             `json:"total_num"`
	FetchNum int             `json:"fetch_num"`
}

type KeyDetail struct {
	//OpenId   string       `json:"open_id"`
	UnionId string `json:"union_id"`
	Id      int    `json:"id"`
	Ap      int    `json:"ap"`
	Type    int    `json:"type"`
	//Content  []KeyContent `json:"content"`
	Content  string `json:"content"`
	Title    string `json:"title"`
	Question string `json:"question"`
	Password string `json:"password"`
	Status   int    `json:"status"`
	Opids    string `json:"opids"`
	UnOpids  string `json:"un_opids"`
}

type KeyContent struct {
	T int    `json:"t"`
	D string `json:"d"`
}

type MainKeyDetail struct {
	KeyDetail
	Sub []SubKeyDetail `json:"sub"`
}

type SubKeyDetail struct {
	KeyDetail
	Sub []SubKeyDetail `json:"sub"`
}

type StageInfo struct {
	Round   int              `json:"round"`
	OnClick map[int][]string `json:"on_click"`
}

type ReconnectInfo struct {
	RoomInfo   RoomInfo    `json:"room_info"`
	GameStage  interface{} `json:"game_stage"`
	GameInfo   GameInfo    `json:"game_info"`
	GameClew   GameClew    `json:"game_clew"`
	NotVoteNum interface{} `json:"not_vote_num"`
	LeftAp     int         `json:"left_ap"`
	StageClick bool        `json:"stage_click"`
}

type Script struct {
	About       string `json:"about"`
	Album       string `json:"album"`
	Author      string `json:"author"`
	Bad         int    `json:"bad"`
	CreateTime  string `json:"createTime"`
	Drama       string `json:"drama"`
	Good        int    `json:"good"`
	Hot         int    `json:"hot"`
	Id          int    `json:"id"`
	Level       int    `json:"level"`
	Name        string `json:"name"`
	Num         int    `json:"num"`
	Price       int    `json:"price"`
	ShowCase    int    `json:"show_case"`
	ExploreFlag bool   `json:"explore_flag"` //探索标识
	TopicFlag   bool   `json:"topic_flag"`   //问题标识
	VoteFlag    bool   `json:"vote_flag"`    //投票标识
	Round       int    `json:"round"`        //剧本轮次
}

type Role struct {
	Album    string `json:"album"`
	Name     string `json:"name"`
	OpenNum  int    `json:"open_num"`
	TotalNum int    `json:"total_num"`
}

type GameClew struct {
	Roles []Role         `json:"roles"`
	Key   map[string]Key `json:"key"`
}

type GameReport struct {
	Score     int    `json:"score"`
	ReportUrl string `json:"report_url"`
}

//用户剧本状态有效期，回复小程序支付签名
type UserScript struct {
	ScriptId   int       `json:"script_id"`
	ScriptName string    `json:"script_name"`
	UserId     int       `json:"user_id"`
	Status     int       `json:"status"`
	EffTime    time.Time `json:"eff_time"`
	InvTime    time.Time `json:"inv_time"`
	CostFlag   int       `json:"cost_flag"` //付款标识 1 未付款 2 已付款
	UnpaidNum  int       `json:"unpaid_num"`
	UnionId    string    `json:"union_id"`
}

type Answer struct {
	AnsId      int    `json:"ans_id"`
	AnsOpt     string `json:"ans_opt"`
	AnsTitle   string `json:"ans_title"`
	NextQueOpt string `json:"next_que_opt"`
	EndOpt     string `json:"end_opt"`
	Score      int    `json:"score"`
	Flag       int    `json:"flag"`
	EndId      int    `json:"end_id"`
	Resume     string `json:"resume"`
}

type AnswerEndingRelation struct {
	Sid    int    `json:"sid"`
	Round  int    `json:"round"`
	AnsIds string `json:"ans_ids"`
	EndId  int    `json:"end_id"`
	Resume string `json:"resume"`
	Remark string `json:"remark"`
}

//用户剧本状态有效期，回复小程序支付签名
type Questions struct {
	QueId    int      `json:"que_id"`
	QueOpt   string   `json:"que_opt"`
	QueTitle string   `json:"que_title"`
	QueType  int      `json:"que_type"`
	MaxScore int      `json:"max_score"` //问题对应答案最高得分
	Ans      []Answer `json:"ans"`
}

//用户剧本状态有效期，回复小程序支付签名
type QuestionAndAnswer struct {
	Question []Questions            `json:"question"`
	Relation []AnswerEndingRelation `json:"relation"`
}

//type AnswerOpt struct {
//	AnsId int `json:"ans_id"`
//}

//用户剧本状态有效期，回复小程序支付签名
type QuestionsOpt struct {
	QueId  int    `json:"que_id"`
	QueOpt string `json:"que_opt"`
	Answer []int  `json:"ans_ids"`
}

//用户剧本问题和得分
type Ending struct {
	Score  int    `json:"score"`
	Resume string `json:"resume"`
}

//用户基本数据
type UserBase struct {
	UnionId  string `json:"union_id"`
	Phone    string `json:"phone"`
	City     string `json:"city"`
	Province string `json:"province"`
	Country  string `json:"country"`
	Gender   int    `json:"gender"`
	NickName string `json:"nick_name"`
	OpenId   string `json:"open_id"`
}

//记录房间数据
type RecordRoomData struct {
	RoomId    string
	UnionId   string
	Owner     string
	StartTime time.Time
	EndTime   time.Time
}

//剧本人物信息
type ScriptPeopleInfo struct {
	ScriptId    int    `json:"sid"`
	RoleId      int    `json:"spid"`
	Name        string `json:"name"`
	Final       string `json:"final"`
	Story       string `json:"story"`
	About       string `json:"about"`
	Infomation  string `json:"infomation"`
	Choice      int    `json:"choice"`
	KeywordFlag int    `json:"keyword_flag"`
	Ap          string `json:"ap"`
	Screen      string `json:"screen"`
}
